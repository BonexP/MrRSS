package freshrss

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"MrRSS/internal/models"
)

// Client represents a FreshRSS API client
type Client struct {
	baseURL    string
	username   string
	password   string
	authToken  string
	httpClient *http.Client
}

// NewClient creates a new FreshRSS API client
func NewClient(serverURL, username, password string) *Client {
	// Ensure URL ends with /api/greader.php
	if !strings.HasSuffix(serverURL, "/api/greader.php") {
		serverURL = strings.TrimSuffix(serverURL, "/") + "/api/greader.php"
	}

	return &Client{
		baseURL:  serverURL,
		username: username,
		password: password,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
			},
		},
	}
}

// Login authenticates with the FreshRSS server and retrieves an auth token
func (c *Client) Login(ctx context.Context) error {
	data := url.Values{}
	data.Set("Email", c.username)
	data.Set("Passwd", c.password)

	req, err := http.NewRequestWithContext(ctx, "POST",
		c.baseURL+"/accounts/ClientLogin",
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read login response: %w", err)
	}

	// Parse response: SID=token\nAuth=token
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Auth=") {
			c.authToken = strings.TrimPrefix(line, "Auth=")
			return nil
		}
	}

	return fmt.Errorf("auth token not found in response")
}

// GetToken retrieves a write token for modifying operations
func (c *Client) GetToken(ctx context.Context) (string, error) {
	if c.authToken == "" {
		return "", fmt.Errorf("not authenticated")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/reader/api/0/token", nil)
	if err != nil {
		return "", fmt.Errorf("create token request: %w", err)
	}

	req.Header.Set("Authorization", "GoogleLogin auth="+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status %d", resp.StatusCode)
	}

	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read token response: %w", err)
	}

	return string(token), nil
}

// Subscription represents a feed subscription
type Subscription struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	URL        string     `json:"url"`
	Categories []Category `json:"categories"`
}

// Category represents a feed category
type Category struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// GetCategories retrieves all categories/tags from FreshRSS
func (c *Client) GetCategories(ctx context.Context) ([]Category, error) {
	if c.authToken == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	req, err := http.NewRequestWithContext(ctx, "GET",
		c.baseURL+"/reader/api/0/tag/list?output=json", nil)
	if err != nil {
		return nil, fmt.Errorf("create categories request: %w", err)
	}

	req.Header.Set("Authorization", "GoogleLogin auth="+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("categories request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("categories request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Tags []struct {
			ID   string `json:"id"`
			Type string `json:"type"` // "folder" or "tag"
		} `json:"tags"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode categories response: %w", err)
	}

	// Convert tags to categories
	categories := make([]Category, 0, len(result.Tags))
	for _, tag := range result.Tags {
		// Extract label from ID (FreshRSS uses "user/-/label/LabelName" format)
		if strings.HasPrefix(tag.ID, "user/-/label/") {
			label := strings.TrimPrefix(tag.ID, "user/-/label/")
			categories = append(categories, Category{
				ID:    tag.ID,
				Label: label,
			})
		}
	}

	return categories, nil
}

// GetSubscriptions retrieves all feed subscriptions
func (c *Client) GetSubscriptions(ctx context.Context) ([]Subscription, error) {
	if c.authToken == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	req, err := http.NewRequestWithContext(ctx, "GET",
		c.baseURL+"/reader/api/0/subscription/list?output=json", nil)
	if err != nil {
		return nil, fmt.Errorf("create subscriptions request: %w", err)
	}

	req.Header.Set("Authorization", "GoogleLogin auth="+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("subscriptions request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("subscriptions request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Subscriptions []Subscription `json:"subscriptions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode subscriptions response: %w", err)
	}

	return result.Subscriptions, nil
}

// Article represents a FreshRSS article
type Article struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	URL        string    `json:"canonical,omitempty"`
	Content    string    `json:"summary,omitempty"`
	Published  time.Time `json:"published"`
	Updated    time.Time `json:"updated"`
	Author     string    `json:"author,omitempty"`
	Categories []string  `json:"categories,omitempty"`
}

// GetUnreadArticles retrieves unread articles
func (c *Client) GetUnreadArticles(ctx context.Context, maxItems int) ([]Article, error) {
	if c.authToken == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	streamURL := fmt.Sprintf("%s/reader/api/0/stream/contents/user/-/state/com.google/reading-list?output=json&n=%d",
		c.baseURL, maxItems)

	req, err := http.NewRequestWithContext(ctx, "GET", streamURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create articles request: %w", err)
	}

	req.Header.Set("Authorization", "GoogleLogin auth="+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("articles request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("articles request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Items []struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			Canonical []struct {
				Href string `json:"href"`
			} `json:"canonical"`
			Summary struct {
				Content string `json:"content"`
			} `json:"summary"`
			Published  int64    `json:"published"`
			Updated    int64    `json:"updated"`
			Author     string   `json:"author"`
			Categories []string `json:"categories"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode articles response: %w", err)
	}

	articles := make([]Article, len(result.Items))
	for i, item := range result.Items {
		var articleURL string
		if len(item.Canonical) > 0 {
			articleURL = item.Canonical[0].Href
		}

		articles[i] = Article{
			ID:         item.ID,
			Title:      item.Title,
			URL:        articleURL,
			Content:    item.Summary.Content,
			Published:  time.Unix(item.Published, 0),
			Updated:    time.Unix(item.Updated, 0),
			Author:     item.Author,
			Categories: item.Categories,
		}
	}

	return articles, nil
}

// MarkAsRead marks articles as read
func (c *Client) MarkAsRead(ctx context.Context, articleIDs []string) error {
	if c.authToken == "" {
		return fmt.Errorf("not authenticated")
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("get token: %w", err)
	}

	data := url.Values{}
	data.Set("T", token)

	for _, id := range articleIDs {
		data.Set("i", id)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		c.baseURL+"/reader/api/0/edit-tag",
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("create mark read request: %w", err)
	}

	req.Header.Set("Authorization", "GoogleLogin auth="+c.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add read tag
	data.Set("a", "user/-/state/com.google/read")
	req.Body = io.NopCloser(strings.NewReader(data.Encode()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("mark read request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("mark read failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SubscribeToFeed subscribes to a new feed
func (c *Client) SubscribeToFeed(ctx context.Context, feedURL, title string) error {
	if c.authToken == "" {
		return fmt.Errorf("not authenticated")
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("get token: %w", err)
	}

	data := url.Values{}
	data.Set("T", token)
	data.Set("s", "feed/"+feedURL)
	if title != "" {
		data.Set("t", title)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		c.baseURL+"/reader/api/0/subscription/edit",
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("create subscribe request: %w", err)
	}

	req.Header.Set("Authorization", "GoogleLogin auth="+c.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add subscription
	data.Set("ac", "subscribe")
	req.Body = io.NopCloser(strings.NewReader(data.Encode()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("subscribe request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("subscribe failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SyncService handles synchronization between MrRSS and FreshRSS
type SyncService struct {
	client *Client
	db     Database
}

// Database interface for FreshRSS sync operations
type Database interface {
	GetFeeds() ([]models.Feed, error)
	AddFeed(feed *models.Feed) (int64, error)
	SaveArticles(ctx context.Context, articles []*models.Article) error
	GetArticles(filter string, feedID int64, category string, showHidden bool, limit, offset int) ([]models.Article, error)
}

// NewSyncService creates a new sync service
func NewSyncService(serverURL, username, password string, db Database) *SyncService {
	return &SyncService{
		client: NewClient(serverURL, username, password),
		db:     db,
	}
}

// Sync performs a bidirectional sync
func (s *SyncService) Sync(ctx context.Context) error {
	// Login to FreshRSS
	if err := s.client.Login(ctx); err != nil {
		return fmt.Errorf("login to FreshRSS: %w", err)
	}

	// Get categories from FreshRSS to build category hierarchy
	categories, err := s.client.GetCategories(ctx)
	if err != nil {
		log.Printf("Failed to get categories, continuing without category sync: %v", err)
		categories = []Category{} // Continue without categories
	}

	// Create category map for quick lookup
	categoryMap := make(map[string]Category)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat
	}

	// Get subscriptions from FreshRSS
	subscriptions, err := s.client.GetSubscriptions(ctx)
	if err != nil {
		return fmt.Errorf("get subscriptions: %w", err)
	}

	// Sync feeds: Add missing feeds to local database
	localFeeds, err := s.db.GetFeeds()
	if err != nil {
		return fmt.Errorf("get local feeds: %w", err)
	}

	// Create a map of local feed URLs for quick lookup
	localFeedMap := make(map[string]int64)
	for _, feed := range localFeeds {
		localFeedMap[feed.URL] = feed.ID
	}

	// Add missing feeds
	for _, sub := range subscriptions {
		if _, exists := localFeedMap[sub.URL]; !exists {
			// Build category path from FreshRSS categories (support nested folders)
			category := s.buildCategoryPath(sub.Categories, categoryMap)

			feed := &models.Feed{
				Title:       sub.Title,
				URL:         sub.URL,
				Category:    category,
				LastUpdated: time.Now(),
			}

			_, err := s.db.AddFeed(feed)
			if err != nil {
				log.Printf("Failed to add feed %s: %v", sub.URL, err)
				continue
			}
			log.Printf("Added feed: %s (category: %s)", sub.Title, category)
		}
	}

	// Get unread articles from FreshRSS
	freshArticles, err := s.client.GetUnreadArticles(ctx, 100) // Get up to 100 unread articles
	if err != nil {
		return fmt.Errorf("get unread articles: %w", err)
	}

	// Create or get FreshRSS feed for synced articles
	freshRSSFeedID, err := s.getOrCreateFreshRSSFeed()
	if err != nil {
		return fmt.Errorf("create FreshRSS feed: %w", err)
	}

	// Get existing FreshRSS articles to avoid duplicates
	existingArticles, err := s.db.GetArticles("unread", freshRSSFeedID, "", false, 1000, 0) // Get unread articles only
	if err != nil {
		return fmt.Errorf("get existing articles: %w", err)
	}

	// Create a map of existing article URLs for quick lookup
	existingArticleMap := make(map[string]bool)
	for _, article := range existingArticles {
		existingArticleMap[article.URL] = true
	}

	// Convert FreshRSS articles to MrRSS articles (only new ones)
	mrssArticles := make([]*models.Article, 0, len(freshArticles))
	for _, freshArt := range freshArticles {
		// Skip if article already exists
		if existingArticleMap[freshArt.URL] {
			continue
		}

		article := &models.Article{
			FeedID:      freshRSSFeedID,
			Title:       freshArt.Title,
			URL:         freshArt.URL,
			Summary:     freshArt.Content, // Store FreshRSS content as summary
			PublishedAt: freshArt.Published,
			IsRead:      false, // FreshRSS unread articles
			IsFavorite:  false,
			IsHidden:    false,
		}
		mrssArticles = append(mrssArticles, article)
	}

	// Save new articles to database
	if len(mrssArticles) > 0 {
		if err := s.db.SaveArticles(ctx, mrssArticles); err != nil {
			return fmt.Errorf("save articles: %w", err)
		}
		log.Printf("Synced %d new articles from FreshRSS", len(mrssArticles))
	}

	log.Printf("FreshRSS sync completed successfully")
	return nil
}
func (s *SyncService) getOrCreateFreshRSSFeed() (int64, error) {
	// Check if FreshRSS feed already exists
	feeds, err := s.db.GetFeeds()
	if err != nil {
		return 0, err
	}

	for _, feed := range feeds {
		if feed.URL == "freshrss://synced" {
			return feed.ID, nil
		}
	}

	// Create new FreshRSS feed
	freshRSSFeed := &models.Feed{
		Title:       "FreshRSS Synced Articles",
		URL:         "freshrss://synced",
		Description: "Articles synced from FreshRSS server",
		Category:    "FreshRSS",
		LastUpdated: time.Now(),
	}

	return s.db.AddFeed(freshRSSFeed)
}

// buildCategoryPath builds a category path from FreshRSS categories
// Supports nested folder structure by parsing category labels that contain "/"
func (s *SyncService) buildCategoryPath(categories []Category, categoryMap map[string]Category) string {
	if len(categories) == 0 {
		return ""
	}

	// Use the first category from the subscription
	categoryID := categories[0].ID

	// Look up the category in our map to get the full label
	if cat, exists := categoryMap[categoryID]; exists {
		label := cat.Label
		// FreshRSS supports nested categories with "/" separator
		// The label itself may contain "/" for hierarchy (e.g., "Tech/News")
		// MrRSS already uses "/" as category separator, so we can use it directly
		return label
	}

	// Fallback: try to extract label from category ID
	if strings.HasPrefix(categoryID, "user/-/label/") {
		label := strings.TrimPrefix(categoryID, "user/-/label/")
		// Check if this label exists in our category map (in case of different ID formats)
		for _, cat := range categoryMap {
			if cat.Label == label {
				return label
			}
		}
		return label
	}

	// Last resort: use the ID as-is if it looks like a label
	if !strings.HasPrefix(categoryID, "user/-/") {
		return categoryID
	}

	return ""
}
