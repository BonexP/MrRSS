package miniflux

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client is a Miniflux REST API v1 client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Miniflux API client
func NewClient(serverURL, apiKey string) *Client {
	baseURL := strings.TrimRight(serverURL, "/")
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL += "/v1"
	}

	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("X-Auth-Token", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) do(req *http.Request, result interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Miniflux API error (status %d): %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("parse response failed: %w", err)
		}
	}

	return nil
}

// User represents a Miniflux user
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// Feed represents a Miniflux feed
type Feed struct {
	ID              int64      `json:"id"`
	Title           string     `json:"title"`
	FeedURL         string     `json:"feed_url"`
	SiteURL         string     `json:"site_url"`
	Category        *Category  `json:"category,omitempty"`
	Icon            *FeedIcon  `json:"icon,omitempty"`
	LastRefreshedAt *time.Time `json:"last_refreshed_at,omitempty"`
}

// FeedIcon represents a Miniflux feed icon
type FeedIcon struct {
	Data     string `json:"data"`
	MimeType string `json:"mime_type"`
}

// Category represents a Miniflux category
type Category struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// Entry represents a Miniflux entry (article)
type Entry struct {
	ID          int64         `json:"id"`
	Title       string        `json:"title"`
	URL         string        `json:"url"`
	FeedID      int64         `json:"feed_id"`
	Content     *EntryContent `json:"content,omitempty"`
	PublishedAt time.Time     `json:"published_at"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   *time.Time    `json:"updated_at,omitempty"`
	Status      string        `json:"status"`
	Starred     bool          `json:"starred"`
	FeedTitle   string        `json:"feed_title,omitempty"`
	Author      *EntryAuthor  `json:"author,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
}

// EntryContent represents the content of a Miniflux entry
type EntryContent struct {
	Content  string `json:"content"`
	MimeType string `json:"mime_type"`
}

// EntryAuthor represents the author of a Miniflux entry
type EntryAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

// EntriesResponse is the response from the entries API
type EntriesResponse struct {
	Total   int     `json:"total"`
	Entries []Entry `json:"entries"`
}

// EntryUpdateRequest is the request body for updating entries
type EntryUpdateRequest struct {
	EntryIDs []int64 `json:"entry_ids"`
	Status   string  `json:"status,omitempty"`
	Starred  *bool   `json:"starred,omitempty"`
}

// TestConnection verifies the Miniflux connection by fetching the current user
func (c *Client) TestConnection(ctx context.Context) error {
	req, err := c.newRequest(ctx, http.MethodGet, "/me", nil)
	if err != nil {
		return err
	}

	var user User
	if err := c.do(req, &user); err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	log.Printf("[Miniflux] Connected as user: %s (ID: %d)", user.Username, user.ID)
	return nil
}

// GetFeeds retrieves all feeds from Miniflux
func (c *Client) GetFeeds(ctx context.Context) ([]Feed, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/feeds", nil)
	if err != nil {
		return nil, err
	}

	var feeds []Feed
	if err := c.do(req, &feeds); err != nil {
		return nil, fmt.Errorf("get feeds failed: %w", err)
	}

	return feeds, nil
}

// GetFeed retrieves a single feed by ID
func (c *Client) GetFeed(ctx context.Context, feedID int64) (*Feed, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/feeds/%d", feedID), nil)
	if err != nil {
		return nil, err
	}

	var feed Feed
	if err := c.do(req, &feed); err != nil {
		return nil, fmt.Errorf("get feed %d failed: %w", feedID, err)
	}

	return &feed, nil
}

// GetCategories retrieves all categories from Miniflux
func (c *Client) GetCategories(ctx context.Context) ([]Category, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/categories", nil)
	if err != nil {
		return nil, err
	}

	var categories []Category
	if err := c.do(req, &categories); err != nil {
		return nil, fmt.Errorf("get categories failed: %w", err)
	}

	return categories, nil
}

// GetEntries retrieves entries with optional filters
func (c *Client) GetEntries(ctx context.Context, filter EntryFilter) (*EntriesResponse, error) {
	params := url.Values{}
	if filter.FeedID > 0 {
		params.Set("feed_id", strconv.FormatInt(filter.FeedID, 10))
	}
	if filter.Status != "" {
		params.Set("status", filter.Status)
	}
	if filter.Starred {
		params.Set("starred", "true")
	}
	if filter.Limit > 0 {
		params.Set("limit", strconv.Itoa(filter.Limit))
	}
	if filter.Offset > 0 {
		params.Set("offset", strconv.Itoa(filter.Offset))
	}
	if filter.Direction != "" {
		params.Set("direction", filter.Direction)
	}
	if filter.AfterEntryID > 0 {
		params.Set("after_entry_id", strconv.FormatInt(filter.AfterEntryID, 10))
	}

	path := "/entries"
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result EntriesResponse
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get entries failed: %w", err)
	}

	return &result, nil
}

// EntryFilter represents filter parameters for GetEntries
type EntryFilter struct {
	FeedID       int64
	Status       string // "read", "unread", "removed"
	Starred      bool
	Limit        int
	Offset       int
	Direction    string // "asc" or "desc"
	AfterEntryID int64
}

// UpdateEntries updates entries (mark as read/unread, star/unstar)
func (c *Client) UpdateEntries(ctx context.Context, entryIDs []int64, status string, starred *bool) error {
	if len(entryIDs) == 0 {
		return nil
	}

	reqBody := EntryUpdateRequest{
		EntryIDs: entryIDs,
		Status:   status,
		Starred:  starred,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPut, "/entries", bytes.NewReader(body))
	if err != nil {
		return err
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("update entries failed: %w", err)
	}

	return nil
}

// GetCounters retrieves unread count for each feed
func (c *Client) GetCounters(ctx context.Context) (map[int64]int, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/feeds/counters", nil)
	if err != nil {
		return nil, err
	}

	var counters map[string]int
	if err := c.do(req, &counters); err != nil {
		return nil, fmt.Errorf("get counters failed: %w", err)
	}

	// Convert string keys to int64
	result := make(map[int64]int, len(counters))
	for k, v := range counters {
		id, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			continue
		}
		result[id] = v
	}

	return result, nil
}
