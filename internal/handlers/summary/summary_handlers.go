package summary

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"MrRSS/internal/handlers/core"
	"MrRSS/internal/summary"
	"MrRSS/internal/utils"
)

// HandleSummarizeArticle generates a summary for an article's content.
func HandleSummarizeArticle(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ArticleID int64  `json:"article_id"`
		Length    string `json:"length"` // "short", "medium", "long"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate length parameter
	summaryLength := summary.Medium
	switch req.Length {
	case "short":
		summaryLength = summary.Short
	case "long":
		summaryLength = summary.Long
	case "medium", "":
		summaryLength = summary.Medium
	default:
		http.Error(w, "Invalid length parameter. Use 'short', 'medium', or 'long'", http.StatusBadRequest)
		return
	}

	// Get the article content
	content, err := getArticleContent(h, req.ArticleID)
	if err != nil {
		log.Printf("Error getting article content for summary: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if content == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"summary":      "",
			"is_too_short": true,
			"error":        "No content available for this article",
		})
		return
	}

	// Get summary provider from settings (with default)
	provider, err := h.DB.GetSetting("summary_provider")
	if err != nil || provider == "" {
		provider = "local" // Default to local algorithm
	}

	var result summary.SummaryResult

	if provider == "ai" {
		// Use AI summarization
		apiKey, err := h.DB.GetSetting("summary_ai_api_key")
		if err != nil || apiKey == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "missing_ai_api_key",
			})
			return
		}

		// Get endpoint and model with fallback to defaults
		endpoint, _ := h.DB.GetSetting("summary_ai_endpoint")
		model, _ := h.DB.GetSetting("summary_ai_model")
		systemPrompt, _ := h.DB.GetSetting("summary_ai_system_prompt")

		aiSummarizer := summary.NewAISummarizer(apiKey, endpoint, model)
		if systemPrompt != "" {
			aiSummarizer.SetSystemPrompt(systemPrompt)
		}
		aiResult, err := aiSummarizer.Summarize(content, summaryLength)
		if err != nil {
			log.Printf("Error generating AI summary: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = aiResult
	} else {
		// Use local algorithm
		summarizer := summary.NewSummarizer()
		result = summarizer.Summarize(content, summaryLength)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"summary":        result.Summary,
		"sentence_count": result.SentenceCount,
		"is_too_short":   result.IsTooShort,
	})
}

// getArticleContent fetches the content of an article by ID
func getArticleContent(h *core.Handler, articleID int64) (string, error) {
	// Get show_hidden_articles setting to include all articles
	allArticles, err := h.DB.GetArticles("", 0, "", true, 50000, 0)
	if err != nil {
		return "", err
	}

	var article *struct {
		FeedID int64
		URL    string
	}

	for _, a := range allArticles {
		if a.ID == articleID {
			article = &struct {
				FeedID int64
				URL    string
			}{
				FeedID: a.FeedID,
				URL:    a.URL,
			}
			break
		}
	}

	if article == nil {
		return "", nil
	}

	// Get the feed URL
	feeds, err := h.DB.GetFeeds()
	if err != nil {
		return "", err
	}

	var feedURL string
	for _, f := range feeds {
		if f.ID == article.FeedID {
			feedURL = f.URL
			break
		}
	}

	if feedURL == "" {
		return "", nil
	}

	// Parse the feed to get fresh content
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	parsedFeed, err := h.Fetcher.ParseFeed(ctx, feedURL)
	if err != nil {
		return "", err
	}

	// Find the article in the feed by URL (use normalized comparison for robustness)
	for _, item := range parsedFeed.Items {
		if utils.URLsMatch(item.Link, article.URL) {
			if item.Content != "" {
				return utils.CleanHTML(item.Content), nil
			}
			return utils.CleanHTML(item.Description), nil
		}
	}

	return "", nil
}
