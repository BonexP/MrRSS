package summary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"MrRSS/internal/config"
)

// AISummarizer implements summarization using OpenAI-compatible APIs (GPT, Claude, etc.).
type AISummarizer struct {
	APIKey       string
	Endpoint     string
	Model        string
	SystemPrompt string
	client       *http.Client
}

// NewAISummarizer creates a new AI summarizer with the given credentials.
// endpoint should be the API base URL (e.g., "https://api.openai.com/v1" for OpenAI)
// model should be the model name (e.g., "gpt-4o-mini", "claude-3-haiku-20240307")
// Summary settings are independent from translation settings.
func NewAISummarizer(apiKey, endpoint, model string) *AISummarizer {
	defaults := config.Get()
	// Use only summary-specific endpoint and model; no fallback to translation settings
	if endpoint == "" {
		endpoint = defaults.SummaryAIEndpoint
	}
	if model == "" {
		model = defaults.SummaryAIModel
	}
	return &AISummarizer{
		APIKey:       apiKey,
		Endpoint:     strings.TrimSuffix(endpoint, "/"),
		Model:        model,
		SystemPrompt: "", // Will be set from settings when used
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

// SetSystemPrompt sets a custom system prompt for the summarizer.
func (s *AISummarizer) SetSystemPrompt(prompt string) {
	s.SystemPrompt = prompt
}

// Summarize generates a summary of the given text using an OpenAI-compatible API.
func (s *AISummarizer) Summarize(text string, length SummaryLength) (SummaryResult, error) {
	// Clean the text first
	cleanedText := cleanText(text)

	// Check if text is too short
	if len(cleanedText) < MinContentLength {
		return SummaryResult{
			Summary:    cleanedText,
			IsTooShort: true,
		}, nil
	}

	// Truncate text if too long to save tokens
	// Use rune slicing to avoid breaking multi-byte UTF-8 characters (e.g., Chinese, emoji)
	runes := []rune(cleanedText)
	if len(runes) > MaxInputCharsForAI {
		cleanedText = string(runes[:MaxInputCharsForAI])
	}

	targetWords := getTargetWordCount(length)

	// Use custom system prompt if provided, otherwise use default
	systemPrompt := s.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = "You are a summarizer. Generate a concise summary of the given text. Output ONLY the summary, nothing else."
	}
	userPrompt := fmt.Sprintf("Summarize the following text in approximately %d words:\n\n%s", targetWords, cleanedText)

	requestBody := map[string]interface{}{
		"model": s.Model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0.3, // Low temperature for consistent summaries
		"max_tokens":  512, // Limit output tokens for summaries
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return SummaryResult{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Validate endpoint URL to prevent SSRF attacks
	apiURL := s.Endpoint + "/chat/completions"
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return SummaryResult{}, fmt.Errorf("invalid API endpoint URL: %w", err)
	}
	if parsedURL.Scheme != "https" {
		return SummaryResult{}, fmt.Errorf("API endpoint must use HTTPS for security")
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return SummaryResult{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return SummaryResult{}, fmt.Errorf("ai api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		// Try to decode error response; if it fails, return generic error with status code
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil && errorResp.Error.Message != "" {
			return SummaryResult{}, fmt.Errorf("ai api error: %s", errorResp.Error.Message)
		}
		return SummaryResult{}, fmt.Errorf("ai api returned status: %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return SummaryResult{}, fmt.Errorf("failed to decode ai response: %w", err)
	}

	if len(result.Choices) > 0 && result.Choices[0].Message.Content != "" {
		// Clean up the response
		summary := strings.TrimSpace(result.Choices[0].Message.Content)

		// Count sentences in the summary
		sentences := splitSentences(summary)

		return SummaryResult{
			Summary:       summary,
			SentenceCount: len(sentences),
			IsTooShort:    false,
		}, nil
	}

	return SummaryResult{}, fmt.Errorf("no summary found in ai response")
}
