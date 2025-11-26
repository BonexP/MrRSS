package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"MrRSS/internal/database"
	"MrRSS/internal/feed"
	"MrRSS/internal/models"
	"MrRSS/internal/translation"
)

func setupTestHandler(t *testing.T) *Handler {
	t.Helper()

	// Create temporary database
	dbFile := "test_handlers.db"
	t.Cleanup(func() { os.Remove(dbFile) })

	db, err := database.NewDB(dbFile)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	if err := db.Init(); err != nil {
		t.Fatalf("Failed to init database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	// Setup fetcher and translator
	translator := translation.NewGoogleFreeTranslator()
	fetcher := feed.NewFetcher(db, translator)

	handler := NewHandler(db, fetcher, translator)
	return handler
}

func TestHandler_HandleArticles(t *testing.T) {
	handler := setupTestHandler(t)

	// Create test data
	feed := &models.Feed{
		Title: "Test Feed",
		URL:   "https://example.com/feed.xml",
	}
	err := handler.DB.AddFeed(feed)
	if err != nil {
		t.Fatalf("Failed to add test feed: %v", err)
	}

	article := &models.Article{
		FeedID:      1, // Assuming first feed gets ID 1
		Title:       "Test Article",
		URL:         "https://example.com/article1",
		Content:     "Test content",
		PublishedAt: time.Now(),
	}
	err = handler.DB.SaveArticle(article)
	if err != nil {
		t.Fatalf("Failed to add test article: %v", err)
	}

	// Create request
	req := httptest.NewRequest("GET", "/api/articles", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.HandleArticles(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []models.Article
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) == 0 {
		t.Error("Expected at least one article")
	}
}

func TestHandler_HandleAddFeed(t *testing.T) {
	t.Skip("Skipping HandleAddFeed test due to network dependency")
}

func TestHandler_HandleFeeds(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/api/feeds", nil)
	w := httptest.NewRecorder()

	handler.HandleFeeds(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []models.Feed
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should return empty array initially
	if len(response) != 0 {
		t.Errorf("Expected empty array, got %d feeds", len(response))
	}
}
