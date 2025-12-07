package article

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"MrRSS/internal/handlers/core"
	"MrRSS/internal/models"
	"MrRSS/internal/utils"
)

// HandleGetArticleContent fetches the article content from RSS feed dynamically.
func HandleGetArticleContent(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	articleIDStr := r.URL.Query().Get("id")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// Get all articles to find the one we need
	allArticles, err := h.DB.GetArticles("", 0, "", false, 1000, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var article *models.Article
	for i := range allArticles {
		if allArticles[i].ID == articleID {
			article = &allArticles[i]
			break
		}
	}

	if article == nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	// Get the feed to fetch fresh content
	feeds, err := h.DB.GetFeeds()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var feedURL string
	for i := range feeds {
		if feeds[i].ID == article.FeedID {
			feedURL = feeds[i].URL
			break
		}
	}

	if feedURL == "" {
		http.Error(w, "Feed not found", http.StatusNotFound)
		return
	}

	// Parse the feed to get fresh content
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	parsedFeed, err := h.Fetcher.ParseFeed(ctx, feedURL)
	if err != nil {
		log.Printf("Error parsing feed for article content: %v", err)
		http.Error(w, "Failed to fetch article content", http.StatusInternalServerError)
		return
	}

	// Find the article in the feed by URL (use normalized comparison for robustness)
	var content string
	for _, item := range parsedFeed.Items {
		if utils.URLsMatch(item.Link, article.URL) {
			content = item.Content
			if content == "" {
				content = item.Description
			}
			// Clean HTML to fix malformed tags that can cause rendering issues
			content = utils.CleanHTML(content)
			break
		}
	}

	json.NewEncoder(w).Encode(map[string]string{
		"content": content,
	})
}
