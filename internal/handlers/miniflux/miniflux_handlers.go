package miniflux

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"MrRSS/internal/handlers/core"
	"MrRSS/internal/handlers/response"
	"MrRSS/internal/miniflux"
)

func getMinifluxClient(h *core.Handler) (*miniflux.Client, error) {
	serverURL, _ := h.DB.GetSetting("miniflux_server_url")
	apiKey, err := h.DB.GetEncryptedSetting("miniflux_api_key")
	if err != nil {
		return nil, err
	}
	return miniflux.NewClient(serverURL, apiKey), nil
}

func checkMinifluxEnabled(h *core.Handler) bool {
	enabled, _ := h.DB.GetSetting("miniflux_enabled")
	return enabled == "true"
}

// HandleFeeds proxies Miniflux feeds list
func HandleFeeds(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}
	if !checkMinifluxEnabled(h) {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}
	client, err := getMinifluxClient(h)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	feeds, err := client.GetFeeds(ctx)
	if err != nil {
		miniflux.Logger.Printf("[Miniflux Handler] HandleFeeds failed: %v", err)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	miniflux.Logger.Printf("[Miniflux Handler] HandleFeeds: %d feeds", len(feeds))
	response.JSON(w, feeds)
}

// HandleCategories proxies Miniflux categories list
func HandleCategories(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}
	if !checkMinifluxEnabled(h) {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}
	client, err := getMinifluxClient(h)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	categories, err := client.GetCategoriesWithCounts(ctx)
	if err != nil {
		miniflux.Logger.Printf("[Miniflux Handler] HandleCategories failed: %v", err)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	miniflux.Logger.Printf("[Miniflux Handler] HandleCategories: %d categories", len(categories))
	response.JSON(w, categories)
}

// HandleCounters proxies Miniflux unread/read counters
func HandleCounters(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}
	if !checkMinifluxEnabled(h) {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}
	client, err := getMinifluxClient(h)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	counters, err := client.GetCounters(ctx)
	if err != nil {
		miniflux.Logger.Printf("[Miniflux Handler] HandleCounters failed: %v", err)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	miniflux.Logger.Printf("[Miniflux Handler] HandleCounters: unreads=%d, reads=%d", len(counters.Unreads), len(counters.Reads))
	response.JSON(w, counters)
}

// HandleArticles proxies Miniflux articles list
func HandleArticles(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}
	if !checkMinifluxEnabled(h) {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}
	client, err := getMinifluxClient(h)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	filter := miniflux.EntryFilter{
		Direction: "desc",
		Order:     "published_at",
		Limit:     50,
	}

	if v := r.URL.Query().Get("feed_id"); v != "" {
		filter.FeedID, _ = strconv.ParseInt(v, 10, 64)
	}
	if v := r.URL.Query().Get("category_id"); v != "" {
		filter.CategoryID, _ = strconv.ParseInt(v, 10, 64)
	}
	if v := r.URL.Query().Get("status"); v != "" {
		filter.Status = v
	}
	if v := r.URL.Query().Get("starred"); v == "true" {
		filter.Starred = true
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			filter.Limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			filter.Offset = n
		}
	}
	if v := r.URL.Query().Get("search"); v != "" {
		filter.Search = v
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	var result *miniflux.EntriesResponse
	var endpoint string
	if filter.FeedID > 0 {
		endpoint = fmt.Sprintf("feed_id=%d", filter.FeedID)
		result, err = client.GetFeedEntries(ctx, filter.FeedID, filter)
	} else if filter.CategoryID > 0 {
		endpoint = fmt.Sprintf("category_id=%d", filter.CategoryID)
		result, err = client.GetCategoryEntries(ctx, filter.CategoryID, filter)
	} else {
		endpoint = "all"
		result, err = client.GetEntries(ctx, filter)
	}
	if err != nil {
		miniflux.Logger.Printf("[Miniflux Handler] HandleArticles failed (%s): %v", endpoint, err)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	miniflux.Logger.Printf("[Miniflux Handler] HandleArticles (%s): %d entries (total=%d)", endpoint, len(result.Entries), result.Total)
	response.JSON(w, result)
}

// HandleArticle proxies a single Miniflux article
func HandleArticle(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}
	if !checkMinifluxEnabled(h) {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}
	client, err := getMinifluxClient(h)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	entry, err := client.GetEntry(ctx, id)
	if err != nil {
		miniflux.Logger.Printf("[Miniflux Handler] HandleArticle failed (id=%d): %v", id, err)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	miniflux.Logger.Printf("[Miniflux Handler] HandleArticle (id=%d): %q", id, entry.Title)
	response.JSON(w, entry)
}

// HandleUpdateStatus updates article status (read/unread/star)
func HandleUpdateStatus(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}
	if !checkMinifluxEnabled(h) {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}
	client, err := getMinifluxClient(h)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	var req struct {
		EntryIDs []int64 `json:"entry_ids"`
		Status   string  `json:"status"`
		Starred  *bool   `json:"starred"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	if err := client.UpdateEntries(ctx, req.EntryIDs, req.Status, req.Starred); err != nil {
		miniflux.Logger.Printf("[Miniflux Handler] HandleUpdateStatus failed: %v", err)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	miniflux.Logger.Printf("[Miniflux Handler] HandleUpdateStatus: %d entries -> status=%s, starred=%v", len(req.EntryIDs), req.Status, req.Starred)
	response.JSON(w, map[string]bool{"success": true})
}

// HandleMarkFeedRead marks all entries in a feed as read
func HandleMarkFeedRead(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}
	if !checkMinifluxEnabled(h) {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}
	client, err := getMinifluxClient(h)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	if err := client.MarkFeedAsRead(ctx, id); err != nil {
		miniflux.Logger.Printf("[Miniflux Handler] HandleMarkFeedRead failed (id=%d): %v", id, err)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	miniflux.Logger.Printf("[Miniflux Handler] HandleMarkFeedRead (id=%d)", id)
	response.JSON(w, map[string]bool{"success": true})
}

// HandleMarkCategoryRead marks all entries in a category as read
func HandleMarkCategoryRead(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}
	if !checkMinifluxEnabled(h) {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}
	client, err := getMinifluxClient(h)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	if err := client.MarkCategoryAsRead(ctx, id); err != nil {
		miniflux.Logger.Printf("[Miniflux Handler] HandleMarkCategoryRead failed (id=%d): %v", id, err)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	miniflux.Logger.Printf("[Miniflux Handler] HandleMarkCategoryRead (id=%d)", id)
	response.JSON(w, map[string]bool{"success": true})
}

// HandleMarkAllRead marks all entries as read
func HandleMarkAllRead(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}
	if !checkMinifluxEnabled(h) {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}
	client, err := getMinifluxClient(h)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	if err := client.MarkAllAsRead(ctx); err != nil {
		miniflux.Logger.Printf("[Miniflux Handler] HandleMarkAllRead failed: %v", err)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	miniflux.Logger.Printf("[Miniflux Handler] HandleMarkAllRead: success")
	response.JSON(w, map[string]bool{"success": true})
}

// HandleTestConnection tests the Miniflux connection
func HandleTestConnection(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}

	serverURL, _ := h.DB.GetSetting("miniflux_server_url")
	apiKey, _ := h.DB.GetEncryptedSetting("miniflux_api_key")
	if serverURL == "" {
		response.JSON(w, map[string]interface{}{
			"success": false,
			"error":   "Miniflux server URL is not configured",
		})
		return
	}
	if apiKey == "" {
		response.JSON(w, map[string]interface{}{
			"success": false,
			"error":   "Miniflux API key is not configured",
		})
		return
	}

	client := miniflux.NewClient(serverURL, apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.TestConnection(ctx); err != nil {
		miniflux.Logger.Printf("[Miniflux Handler] Connection test failed: %v", err)
		response.JSON(w, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	miniflux.Logger.Printf("[Miniflux Handler] Connection test successful (server=%s)", serverURL)

	response.JSON(w, map[string]interface{}{
		"success": true,
		"message": "Connection to Miniflux server successful",
	})
}

// HandleSync performs a sync with the Miniflux server:
// pulls feeds/categories/counters and pushes pending status changes
func HandleSync(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}
	if !checkMinifluxEnabled(h) {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}

	serverURL, _ := h.DB.GetSetting("miniflux_server_url")
	apiKey, _ := h.DB.GetEncryptedSetting("miniflux_api_key")
	if serverURL == "" || apiKey == "" {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}

	client := miniflux.NewClient(serverURL, apiKey)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		miniflux.Logger.Printf("[Miniflux Sync] Starting sync...")

		if err := client.TestConnection(ctx); err != nil {
			miniflux.Logger.Printf("[Miniflux Sync] Connection failed: %v", err)
			return
		}

		feeds, err := client.GetFeeds(ctx)
		if err != nil {
			miniflux.Logger.Printf("[Miniflux Sync] Failed to get feeds: %v", err)
			return
		}
		miniflux.Logger.Printf("[Miniflux Sync] Fetched %d feeds", len(feeds))

		_, err = client.GetCategoriesWithCounts(ctx)
		if err != nil {
			miniflux.Logger.Printf("[Miniflux Sync] Failed to get categories: %v", err)
		}

		_, err = client.GetCounters(ctx)
		if err != nil {
			miniflux.Logger.Printf("[Miniflux Sync] Failed to get counters: %v", err)
		}

		lastSyncTime := time.Now().Format(time.RFC3339)
		_ = h.DB.SetSetting("miniflux_last_sync_time", lastSyncTime)

		miniflux.Logger.Printf("[Miniflux Sync] Sync completed: %d feeds, %s", len(feeds), lastSyncTime)
	}()

	response.JSON(w, map[string]interface{}{
		"status":  "sync_started",
		"message": "Miniflux synchronization started",
	})
}

// HandleSyncStatus returns the current Miniflux sync status
func HandleSyncStatus(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}

	lastSyncStr, _ := h.DB.GetSetting("miniflux_last_sync_time")
	var lastSyncTime *time.Time
	if lastSyncStr != "" {
		if ts, err := time.Parse(time.RFC3339, lastSyncStr); err == nil {
			lastSyncTime = &ts
		}
	}

	response.JSON(w, map[string]interface{}{
		"last_sync_time": lastSyncTime,
	})
}
