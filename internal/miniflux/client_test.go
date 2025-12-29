package miniflux

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"MrRSS/internal/models"
)

// TestNewClient tests the client creation
func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		serverURL string
		apiKey    string
		wantURL   string
	}{
		{
			name:      "URL without /v1 suffix",
			serverURL: "https://miniflux.example.com",
			apiKey:    "test-api-key",
			wantURL:   "https://miniflux.example.com/v1",
		},
		{
			name:      "URL with /v1 suffix",
			serverURL: "https://miniflux.example.com/v1",
			apiKey:    "test-api-key",
			wantURL:   "https://miniflux.example.com/v1",
		},
		{
			name:      "URL with trailing slash",
			serverURL: "https://miniflux.example.com/",
			apiKey:    "test-api-key",
			wantURL:   "https://miniflux.example.com/v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.serverURL, tt.apiKey)
			if client.baseURL != tt.wantURL {
				t.Errorf("NewClient() baseURL = %v, want %v", client.baseURL, tt.wantURL)
			}
			if client.apiKey != tt.apiKey {
				t.Errorf("NewClient() apiKey = %v, want %v", client.apiKey, tt.apiKey)
			}
		})
	}
}

// TestTestConnection tests the connection testing functionality
func TestTestConnection(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "Successful connection",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "Unauthorized",
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "Internal server error",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check authentication header
				if r.Header.Get("X-Auth-Token") == "" {
					t.Error("Expected X-Auth-Token header")
				}
				// Check endpoint
				if r.URL.Path != "/v1/me" {
					t.Errorf("Expected /v1/me, got %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := NewClient(server.URL, "test-api-key")
			ctx := context.Background()
			err := client.TestConnection(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("TestConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetFeeds tests the feeds retrieval functionality
func TestGetFeeds(t *testing.T) {
	mockFeeds := []Feed{
		{
			ID:      1,
			Title:   "Test Feed 1",
			FeedURL: "https://example.com/feed1.xml",
			SiteURL: "https://example.com",
			Category: struct {
				ID    int64  `json:"id"`
				Title string `json:"title"`
			}{ID: 1, Title: "Tech"},
		},
		{
			ID:      2,
			Title:   "Test Feed 2",
			FeedURL: "https://example.com/feed2.xml",
			SiteURL: "https://example.com",
			Category: struct {
				ID    int64  `json:"id"`
				Title string `json:"title"`
			}{ID: 2, Title: "News"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authentication
		if r.Header.Get("X-Auth-Token") == "" {
			t.Error("Expected X-Auth-Token header")
		}
		// Check endpoint
		if r.URL.Path != "/v1/feeds" {
			t.Errorf("Expected /v1/feeds, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockFeeds)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")
	ctx := context.Background()
	feeds, err := client.GetFeeds(ctx)

	if err != nil {
		t.Fatalf("GetFeeds() error = %v", err)
	}

	if len(feeds) != len(mockFeeds) {
		t.Errorf("GetFeeds() got %d feeds, want %d", len(feeds), len(mockFeeds))
	}

	if feeds[0].Title != mockFeeds[0].Title {
		t.Errorf("GetFeeds() feed title = %v, want %v", feeds[0].Title, mockFeeds[0].Title)
	}
}

// TestGetEntries tests the entries retrieval functionality
func TestGetEntries(t *testing.T) {
	mockEntries := []Entry{
		{
			ID:          1,
			Title:       "Test Entry 1",
			URL:         "https://example.com/entry1",
			Content:     "<p>Test content 1</p>",
			PublishedAt: time.Now(),
			Status:      "unread",
			Starred:     false,
		},
		{
			ID:          2,
			Title:       "Test Entry 2",
			URL:         "https://example.com/entry2",
			Content:     "<p>Test content 2</p>",
			PublishedAt: time.Now(),
			Status:      "unread",
			Starred:     true,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authentication
		if r.Header.Get("X-Auth-Token") == "" {
			t.Error("Expected X-Auth-Token header")
		}
		// Check endpoint and query params
		if r.URL.Path != "/v1/entries" {
			t.Errorf("Expected /v1/entries, got %s", r.URL.Path)
		}
		status := r.URL.Query().Get("status")
		if status != "unread" {
			t.Errorf("Expected status=unread, got %s", status)
		}
		limit := r.URL.Query().Get("limit")
		if limit != "100" {
			t.Errorf("Expected limit=100, got %s", limit)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total":   len(mockEntries),
			"entries": mockEntries,
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")
	ctx := context.Background()
	entries, err := client.GetEntries(ctx, "unread", 100)

	if err != nil {
		t.Fatalf("GetEntries() error = %v", err)
	}

	if len(entries) != len(mockEntries) {
		t.Errorf("GetEntries() got %d entries, want %d", len(entries), len(mockEntries))
	}

	if entries[0].Title != mockEntries[0].Title {
		t.Errorf("GetEntries() entry title = %v, want %v", entries[0].Title, mockEntries[0].Title)
	}

	if !entries[1].Starred {
		t.Error("GetEntries() expected entry to be starred")
	}
}

// TestUpdateEntries tests the entries update functionality
func TestUpdateEntries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authentication
		if r.Header.Get("X-Auth-Token") == "" {
			t.Error("Expected X-Auth-Token header")
		}
		// Check method
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}
		// Check endpoint
		if r.URL.Path != "/v1/entries" {
			t.Errorf("Expected /v1/entries, got %s", r.URL.Path)
		}
		// Check content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Expected Content-Type: application/json")
		}

		// Parse body
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")
	ctx := context.Background()
	err := client.UpdateEntries(ctx, []int64{1, 2, 3}, "read")

	if err != nil {
		t.Errorf("UpdateEntries() error = %v", err)
	}
}

// MockDatabase implements the Database interface for testing
type MockDatabase struct {
	feeds    []models.Feed
	articles []models.Article
}

func (m *MockDatabase) GetFeeds() ([]models.Feed, error) {
	return m.feeds, nil
}

func (m *MockDatabase) AddFeed(feed *models.Feed) (int64, error) {
	feed.ID = int64(len(m.feeds) + 1)
	m.feeds = append(m.feeds, *feed)
	return feed.ID, nil
}

func (m *MockDatabase) SaveArticles(ctx context.Context, articles []*models.Article) error {
	for _, article := range articles {
		m.articles = append(m.articles, *article)
	}
	return nil
}

func (m *MockDatabase) GetArticles(filter string, feedID int64, category string, showHidden bool, limit, offset int) ([]models.Article, error) {
	result := []models.Article{}
	for _, article := range m.articles {
		if article.FeedID == feedID {
			result = append(result, article)
		}
	}
	return result, nil
}

// TestSync tests the full synchronization workflow
func TestSync(t *testing.T) {
	mockFeeds := []Feed{
		{
			ID:      1,
			Title:   "Test Feed",
			FeedURL: "https://example.com/feed.xml",
			SiteURL: "https://example.com",
			Category: struct {
				ID    int64  `json:"id"`
				Title string `json:"title"`
			}{ID: 1, Title: "Tech"},
		},
	}

	mockEntries := []Entry{
		{
			ID:          1,
			Title:       "Test Entry",
			URL:         "https://example.com/entry1",
			Content:     "<p>Test content</p>",
			PublishedAt: time.Now(),
			Status:      "unread",
			Starred:     false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/feeds":
			json.NewEncoder(w).Encode(mockFeeds)
		case "/v1/entries":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"total":   len(mockEntries),
				"entries": mockEntries,
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	mockDB := &MockDatabase{
		feeds:    []models.Feed{},
		articles: []models.Article{},
	}

	syncService := NewSyncService(server.URL, "test-api-key", mockDB)
	ctx := context.Background()

	err := syncService.Sync(ctx)
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	// Check if feeds were added (1 from Miniflux + 1 special Miniflux feed)
	if len(mockDB.feeds) != 2 {
		t.Errorf("Expected 2 feeds after sync, got %d", len(mockDB.feeds))
	}

	// Check if articles were added
	if len(mockDB.articles) != 1 {
		t.Errorf("Expected 1 article after sync, got %d", len(mockDB.articles))
	}

	// Verify the special Miniflux feed was created
	foundSpecialFeed := false
	for _, feed := range mockDB.feeds {
		if feed.URL == "miniflux://synced" {
			foundSpecialFeed = true
			if feed.Title != "Miniflux Synced Articles" {
				t.Errorf("Special feed title = %v, want 'Miniflux Synced Articles'", feed.Title)
			}
		}
	}
	if !foundSpecialFeed {
		t.Error("Special Miniflux feed was not created")
	}
}

// TestSyncDeduplication tests that duplicate articles are not added
func TestSyncDeduplication(t *testing.T) {
	mockFeeds := []Feed{
		{
			ID:      1,
			Title:   "Test Feed",
			FeedURL: "https://example.com/feed.xml",
			SiteURL: "https://example.com",
			Category: struct {
				ID    int64  `json:"id"`
				Title string `json:"title"`
			}{ID: 1, Title: "Tech"},
		},
	}

	mockEntries := []Entry{
		{
			ID:          1,
			Title:       "Test Entry",
			URL:         "https://example.com/entry1",
			Content:     "<p>Test content</p>",
			PublishedAt: time.Now(),
			Status:      "unread",
			Starred:     false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/feeds":
			json.NewEncoder(w).Encode(mockFeeds)
		case "/v1/entries":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"total":   len(mockEntries),
				"entries": mockEntries,
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	mockDB := &MockDatabase{
		feeds:    []models.Feed{},
		articles: []models.Article{},
	}

	syncService := NewSyncService(server.URL, "test-api-key", mockDB)
	ctx := context.Background()

	// First sync
	err := syncService.Sync(ctx)
	if err != nil {
		t.Fatalf("First Sync() error = %v", err)
	}

	articlesAfterFirstSync := len(mockDB.articles)

	// Second sync - should not add duplicate articles
	err = syncService.Sync(ctx)
	if err != nil {
		t.Fatalf("Second Sync() error = %v", err)
	}

	if len(mockDB.articles) != articlesAfterFirstSync {
		t.Errorf("Expected %d articles after second sync (no duplicates), got %d",
			articlesAfterFirstSync, len(mockDB.articles))
	}
}
