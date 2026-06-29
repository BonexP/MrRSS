package miniflux

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	userID     int64
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
	FeedID   int64  `json:"feed_id"`
	IconID   int64  `json:"icon_id"`
	Data     string `json:"data"`
	MimeType string `json:"mime_type"`
}

// Category represents a Miniflux category
type Category struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	FeedCount   int    `json:"feed_count,omitempty"`
	TotalUnread int    `json:"total_unread,omitempty"`
}

// Entry represents a Miniflux entry (article)
type Entry struct {
	ID          int64         `json:"id"`
	Title       string        `json:"title"`
	URL         string        `json:"url"`
	FeedID      int64         `json:"feed_id"`
	Content     string        `json:"content,omitempty"`
	PublishedAt time.Time     `json:"published_at"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   *time.Time    `json:"updated_at,omitempty"`
	Status      string        `json:"status"`
	Starred     bool          `json:"starred"`
	FeedTitle   string        `json:"feed_title,omitempty"`
	Author      string        `json:"author,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
	ReadingTime int           `json:"reading_time"`
	Enclosures  []interface{} `json:"enclosures"`
	Feed        *Feed         `json:"feed,omitempty"`
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

	c.userID = user.ID
	Logger.Printf("[Miniflux] Connected as user: %s (ID: %d)", user.Username, user.ID)
	return nil
}

// GetUserID returns the cached user ID (must call TestConnection first)
func (c *Client) GetUserID(ctx context.Context) (int64, error) {
	if c.userID > 0 {
		return c.userID, nil
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/me", nil)
	if err != nil {
		return 0, err
	}

	var user User
	if err := c.do(req, &user); err != nil {
		return 0, fmt.Errorf("get user failed: %w", err)
	}

	c.userID = user.ID
	return user.ID, nil
}

// GetFeeds retrieves all feeds from Miniflux
func (c *Client) GetFeeds(ctx context.Context) ([]Feed, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/feeds", nil)
	if err != nil {
		return nil, err
	}

	var feeds []Feed
	if err := c.do(req, &feeds); err != nil {
		Logger.Printf("[Miniflux Client] GetFeeds failed: %v", err)
		return nil, fmt.Errorf("get feeds failed: %w", err)
	}

	Logger.Printf("[Miniflux Client] GetFeeds: %d feeds", len(feeds))
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

// GetFeedIcon retrieves a feed icon by feed ID
func (c *Client) GetFeedIcon(ctx context.Context, feedID int64) (*FeedIcon, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/feeds/%d/icon", feedID), nil)
	if err != nil {
		return nil, err
	}

	var icon FeedIcon
	if err := c.do(req, &icon); err != nil {
		return nil, fmt.Errorf("get feed icon %d failed: %w", feedID, err)
	}

	return &icon, nil
}

// GetFeedIconByIconID retrieves a feed icon by icon ID
func (c *Client) GetFeedIconByIconID(ctx context.Context, iconID int64) (*FeedIcon, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/icons/%d", iconID), nil)
	if err != nil {
		return nil, err
	}

	var icon FeedIcon
	if err := c.do(req, &icon); err != nil {
		return nil, fmt.Errorf("get icon %d failed: %w", iconID, err)
	}

	return &icon, nil
}

// GetCategories retrieves all categories from Miniflux
func (c *Client) GetCategories(ctx context.Context) ([]Category, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/categories", nil)
	if err != nil {
		return nil, err
	}

	var categories []Category
	if err := c.do(req, &categories); err != nil {
		Logger.Printf("[Miniflux Client] GetCategories failed: %v", err)
		return nil, fmt.Errorf("get categories failed: %w", err)
	}

	return categories, nil
}

// GetCategoriesWithCounts retrieves all categories with unread/feed counts
func (c *Client) GetCategoriesWithCounts(ctx context.Context) ([]Category, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/categories?counts=true", nil)
	if err != nil {
		return nil, err
	}

	var categories []Category
	if err := c.do(req, &categories); err != nil {
		Logger.Printf("[Miniflux Client] GetCategoriesWithCounts failed: %v", err)
		return nil, fmt.Errorf("get categories with counts failed: %w", err)
	}

	Logger.Printf("[Miniflux Client] GetCategoriesWithCounts: %d categories", len(categories))
	return categories, nil
}

// GetEntry retrieves a single entry by ID
func (c *Client) GetEntry(ctx context.Context, entryID int64) (*Entry, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/entries/%d", entryID), nil)
	if err != nil {
		return nil, err
	}

	var entry Entry
	if err := c.do(req, &entry); err != nil {
		return nil, fmt.Errorf("get entry %d failed: %w", entryID, err)
	}

	return &entry, nil
}

// GetFeedEntry retrieves a single entry within a feed context
func (c *Client) GetFeedEntry(ctx context.Context, feedID, entryID int64) (*Entry, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/feeds/%d/entries/%d", feedID, entryID), nil)
	if err != nil {
		return nil, err
	}

	var entry Entry
	if err := c.do(req, &entry); err != nil {
		return nil, fmt.Errorf("get feed entry %d/%d failed: %w", feedID, entryID, err)
	}

	return &entry, nil
}

func buildEntryFilterParams(filter EntryFilter) string {
	params := url.Values{}
	if filter.FeedID > 0 {
		params.Set("feed_id", strconv.FormatInt(filter.FeedID, 10))
	}
	if filter.CategoryID > 0 {
		params.Set("category_id", strconv.FormatInt(filter.CategoryID, 10))
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
	if filter.Order != "" {
		params.Set("order", filter.Order)
	}
	if filter.Search != "" {
		params.Set("search", filter.Search)
	}
	if filter.AfterEntryID > 0 {
		params.Set("after_entry_id", strconv.FormatInt(filter.AfterEntryID, 10))
	}
	return params.Encode()
}

// GetEntries retrieves entries with optional filters
func (c *Client) GetEntries(ctx context.Context, filter EntryFilter) (*EntriesResponse, error) {
	path := "/entries"
	if encoded := buildEntryFilterParams(filter); encoded != "" {
		path += "?" + encoded
	}

	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result EntriesResponse
	if err := c.do(req, &result); err != nil {
		Logger.Printf("[Miniflux Client] GetEntries failed: %v", err)
		return nil, fmt.Errorf("get entries failed: %w", err)
	}

	Logger.Printf("[Miniflux Client] GetEntries: %d entries (total=%d)", len(result.Entries), result.Total)
	return &result, nil
}

// GetFeedEntries retrieves entries for a specific feed
func (c *Client) GetFeedEntries(ctx context.Context, feedID int64, filter EntryFilter) (*EntriesResponse, error) {
	filter.FeedID = feedID
	path := fmt.Sprintf("/feeds/%d/entries", feedID)
	if encoded := buildEntryFilterParams(filter); encoded != "" {
		path += "?" + encoded
	}

	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result EntriesResponse
	if err := c.do(req, &result); err != nil {
		Logger.Printf("[Miniflux Client] GetFeedEntries (feed=%d) failed: %v", feedID, err)
		return nil, fmt.Errorf("get feed entries failed: %w", err)
	}

	Logger.Printf("[Miniflux Client] GetFeedEntries (feed=%d): %d entries (total=%d)", feedID, len(result.Entries), result.Total)
	return &result, nil
}

// GetCategoryEntries retrieves entries for a specific category
func (c *Client) GetCategoryEntries(ctx context.Context, categoryID int64, filter EntryFilter) (*EntriesResponse, error) {
	filter.CategoryID = categoryID
	path := fmt.Sprintf("/categories/%d/entries", categoryID)
	if encoded := buildEntryFilterParams(filter); encoded != "" {
		path += "?" + encoded
	}

	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result EntriesResponse
	if err := c.do(req, &result); err != nil {
		Logger.Printf("[Miniflux Client] GetCategoryEntries (category=%d) failed: %v", categoryID, err)
		return nil, fmt.Errorf("get category entries failed: %w", err)
	}

	Logger.Printf("[Miniflux Client] GetCategoryEntries (category=%d): %d entries (total=%d)", categoryID, len(result.Entries), result.Total)
	return &result, nil
}

// CountersResponse is the response from the counters API
type CountersResponse struct {
	Reads   map[string]int `json:"reads"`
	Unreads map[string]int `json:"unreads"`
}

// EntryFilter represents filter parameters for GetEntries
type EntryFilter struct {
	FeedID       int64
	CategoryID   int64
	Status       string // "read", "unread", "removed"
	Starred      bool
	Limit        int
	Offset       int
	Direction    string // "asc" or "desc"
	Order        string // "id", "status", "published_at", "category_title", "category_id"
	Search       string
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
		Logger.Printf("[Miniflux Client] UpdateEntries failed (%d entries, status=%s): %v", len(entryIDs), status, err)
		return fmt.Errorf("update entries failed: %w", err)
	}

	Logger.Printf("[Miniflux Client] UpdateEntries: %d entries -> status=%s", len(entryIDs), status)
	return nil
}

// ToggleBookmark toggles the bookmark status of an entry
func (c *Client) ToggleBookmark(ctx context.Context, entryID int64) error {
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/entries/%d/bookmark", entryID), nil)
	if err != nil {
		return err
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("toggle bookmark failed: %w", err)
	}

	return nil
}

// MarkFeedAsRead marks all entries in a feed as read
func (c *Client) MarkFeedAsRead(ctx context.Context, feedID int64) error {
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/feeds/%d/mark-all-as-read", feedID), nil)
	if err != nil {
		return err
	}

	if err := c.do(req, nil); err != nil {
		Logger.Printf("[Miniflux Client] MarkFeedAsRead (feed=%d) failed: %v", feedID, err)
		return fmt.Errorf("mark feed as read failed: %w", err)
	}

	Logger.Printf("[Miniflux Client] MarkFeedAsRead (feed=%d)", feedID)
	return nil
}

// MarkCategoryAsRead marks all entries in a category as read
func (c *Client) MarkCategoryAsRead(ctx context.Context, categoryID int64) error {
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/categories/%d/mark-all-as-read", categoryID), nil)
	if err != nil {
		return err
	}

	if err := c.do(req, nil); err != nil {
		Logger.Printf("[Miniflux Client] MarkCategoryAsRead (category=%d) failed: %v", categoryID, err)
		return fmt.Errorf("mark category as read failed: %w", err)
	}

	Logger.Printf("[Miniflux Client] MarkCategoryAsRead (category=%d)", categoryID)
	return nil
}

// MarkAllAsRead marks all entries for the current user as read
func (c *Client) MarkAllAsRead(ctx context.Context) error {
	userID, err := c.GetUserID(ctx)
	if err != nil {
		return fmt.Errorf("get user ID: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/users/%d/mark-all-as-read", userID), nil)
	if err != nil {
		return err
	}

	if err := c.do(req, nil); err != nil {
		Logger.Printf("[Miniflux Client] MarkAllAsRead failed: %v", err)
		return fmt.Errorf("mark all as read failed: %w", err)
	}

	Logger.Printf("[Miniflux Client] MarkAllAsRead: success")
	return nil
}

// GetCounters retrieves unread/read counts for each feed
func (c *Client) GetCounters(ctx context.Context) (*CountersResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/feeds/counters", nil)
	if err != nil {
		return nil, err
	}

	var counters CountersResponse
	if err := c.do(req, &counters); err != nil {
		Logger.Printf("[Miniflux Client] GetCounters failed: %v", err)
		return nil, fmt.Errorf("get counters failed: %w", err)
	}

	Logger.Printf("[Miniflux Client] GetCounters: unreads=%d, reads=%d", len(counters.Unreads), len(counters.Reads))
	return &counters, nil
}
