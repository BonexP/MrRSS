package feed

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"MrRSS/internal/handlers/core"
)

// HandleFeeds returns all feeds.
func HandleFeeds(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	feeds, err := h.DB.GetFeeds()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(feeds)
}

// HandleAddFeed adds a new feed subscription and immediately fetches its articles.
func HandleAddFeed(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL             string `json:"url"`
		Category        string `json:"category"`
		Title           string `json:"title"`
		ScriptPath      string `json:"script_path"`
		HideFromTimeline bool  `json:"hide_from_timeline"`
		ProxyURL        string `json:"proxy_url"`
		ProxyEnabled    bool   `json:"proxy_enabled"`
		RefreshInterval int    `json:"refresh_interval"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var feedID int64
	var err error
	if req.ScriptPath != "" {
		// Add feed using custom script
		feedID, err = h.Fetcher.AddScriptSubscription(req.ScriptPath, req.Category, req.Title)
	} else {
		// Add feed using URL
		feedID, err = h.Fetcher.AddSubscription(req.URL, req.Category, req.Title)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update all feed settings
	feed, err := h.DB.GetFeedByID(feedID)
	if err != nil {
		// Log the error but don't fail the request - feed was created successfully
		// The settings can be set later via edit
		http.Error(w, "feed created but failed to update settings: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.DB.UpdateFeed(feed.ID, feed.Title, feed.URL, feed.Category, feed.ScriptPath, req.HideFromTimeline, req.ProxyURL, req.ProxyEnabled, req.RefreshInterval); err != nil {
		http.Error(w, "feed created but failed to update settings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Immediately fetch articles for the newly added feed in background
	go func() {
		feed, err := h.DB.GetFeedByID(feedID)
		if err != nil {
			return
		}
		h.Fetcher.FetchSingleFeed(context.Background(), *feed)
	}()

	w.WriteHeader(http.StatusOK)
}

// HandleDeleteFeed deletes a feed subscription.
func HandleDeleteFeed(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.DB.DeleteFeed(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleUpdateFeed updates a feed's properties.
func HandleUpdateFeed(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID              int64  `json:"id"`
		Title           string `json:"title"`
		URL             string `json:"url"`
		Category        string `json:"category"`
		ScriptPath      string `json:"script_path"`
		HideFromTimeline bool  `json:"hide_from_timeline"`
		ProxyURL        string `json:"proxy_url"`
		ProxyEnabled    bool   `json:"proxy_enabled"`
		RefreshInterval int    `json:"refresh_interval"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.DB.UpdateFeed(req.ID, req.Title, req.URL, req.Category, req.ScriptPath, req.HideFromTimeline, req.ProxyURL, req.ProxyEnabled, req.RefreshInterval); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleRefreshFeed refreshes a single feed by ID with progress tracking.
func HandleRefreshFeed(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid feed ID", http.StatusBadRequest)
		return
	}

	feed, err := h.DB.GetFeedByID(id)
	if err != nil {
		http.Error(w, "Feed not found", http.StatusNotFound)
		return
	}

	// Refresh the feed in background with progress tracking
	go h.Fetcher.FetchSingleFeed(context.Background(), *feed)

	w.WriteHeader(http.StatusOK)
}
