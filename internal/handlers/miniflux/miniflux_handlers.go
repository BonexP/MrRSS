package miniflux

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"MrRSS/internal/handlers/core"
	"MrRSS/internal/handlers/response"
	"MrRSS/internal/miniflux"
)

// HandleSync starts a full Miniflux bidirectional sync
func HandleSync(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}

	enabled, err := h.DB.GetSetting("miniflux_enabled")
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	if enabled != "true" {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}

	serverURL, apiKey, err := h.DB.GetMinifluxConfig()
	if err != nil || serverURL == "" || apiKey == "" {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}

	syncService := miniflux.NewBidirectionalSyncService(serverURL, apiKey, h.DB)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		result, err := syncService.Sync(ctx)
		if err != nil {
			log.Printf("[Miniflux Handler] Sync error: %v", err)
		}
		log.Printf("[Miniflux Handler] Sync completed: pull=%d push=%d errors=%d",
			result.PullChangesCount, result.PushChangesCount, len(result.Errors))

		_ = h.DB.SetSetting("miniflux_last_sync_time", time.Now().Format(time.RFC3339))
	}()

	response.JSON(w, map[string]interface{}{
		"status":  "sync_started",
		"message": "Miniflux synchronization started",
	})
}

// HandleSyncFeed starts a sync for a single Miniflux feed
func HandleSyncFeed(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}

	feedIDStr := r.URL.Query().Get("feed_id")
	if feedIDStr == "" {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}

	feedID, err := strconv.ParseInt(feedIDStr, 10, 64)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	enabled, err := h.DB.GetSetting("miniflux_enabled")
	if err != nil || enabled != "true" {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}

	serverURL, apiKey, err := h.DB.GetMinifluxConfig()
	if err != nil || serverURL == "" || apiKey == "" {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}

	syncService := miniflux.NewBidirectionalSyncService(serverURL, apiKey, h.DB)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		count, err := syncService.SyncFeed(ctx, feedID)
		if err != nil {
			log.Printf("[Miniflux Handler] SyncFeed error for feed %d: %v", feedID, err)
		} else {
			log.Printf("[Miniflux Handler] SyncFeed completed for feed %d: %d entries", feedID, count)
		}
	}()

	response.JSON(w, map[string]interface{}{
		"status":  "sync_started",
		"message": "Feed synchronization started",
	})
}

// HandleTestConnection tests the Miniflux connection
func HandleTestConnection(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}

	serverURL, apiKey, err := h.DB.GetMinifluxConfig()
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
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
		log.Printf("[Miniflux Handler] Connection test failed: %v", err)
		response.JSON(w, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response.JSON(w, map[string]interface{}{
		"success": true,
		"message": "Connection to Miniflux server successful",
	})
}

// HandleSyncStatus returns the current Miniflux sync status
func HandleSyncStatus(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, nil, http.StatusMethodNotAllowed)
		return
	}

	pendingCount, _ := h.DB.GetPendingMinifluxSyncCount()
	failedItems, _ := h.DB.GetFailedMinifluxSyncItems(5)
	lastSyncTime, _ := h.DB.GetSetting("miniflux_last_sync_time")

	type failedItem struct {
		ArticleID  int64  `json:"article_id"`
		ArticleURL string `json:"article_url"`
		Action     string `json:"action"`
		Error      string `json:"error"`
	}

	var failed []failedItem
	for _, item := range failedItems {
		errMsg := ""
		if item.SyncError != nil {
			errMsg = *item.SyncError
		}
		failed = append(failed, failedItem{
			ArticleID:  item.ArticleID,
			ArticleURL: item.ArticleURL,
			Action:     string(item.Action),
			Error:      errMsg,
		})
	}

	response.JSON(w, map[string]interface{}{
		"pending_changes": pendingCount,
		"failed_items":    failed,
		"last_sync_time":  lastSyncTime,
	})
}
