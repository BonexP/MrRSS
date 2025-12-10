package models

import "time"

type Feed struct {
	ID                 int64     `json:"id"`
	Title              string    `json:"title"`
	URL                string    `json:"url"`
	Link               string    `json:"link"` // Website homepage link
	Description        string    `json:"description"`
	Category           string    `json:"category"`
	ImageURL           string    `json:"image_url"` // New field
	LastUpdated        time.Time `json:"last_updated"`
	LastError          string    `json:"last_error,omitempty"`  // Track last fetch error
	DiscoveryCompleted bool      `json:"discovery_completed"`   // Track if discovery has been run
	ScriptPath         string    `json:"script_path,omitempty"` // Path to custom script for fetching feed
	HideFromTimeline   bool      `json:"hide_from_timeline"`    // Hide articles from timeline views
	ProxyURL           string    `json:"proxy_url,omitempty"`   // Custom proxy URL for this feed (overrides global)
	ProxyEnabled       bool      `json:"proxy_enabled"`         // Whether to use proxy for this feed
	RefreshInterval    int       `json:"refresh_interval"`      // Custom refresh interval in minutes (0 = use global)
}

type Article struct {
	ID              int64     `json:"id"`
	FeedID          int64     `json:"feed_id"`
	Title           string    `json:"title"`
	URL             string    `json:"url"`
	ImageURL        string    `json:"image_url"`
	AudioURL        string    `json:"audio_url"`
	VideoURL        string    `json:"video_url"` // YouTube video URL for embedded player
	PublishedAt     time.Time `json:"published_at"`
	IsRead          bool      `json:"is_read"`
	IsFavorite      bool      `json:"is_favorite"`
	IsHidden        bool      `json:"is_hidden"`
	IsReadLater     bool      `json:"is_read_later"`
	FeedTitle       string    `json:"feed_title,omitempty"` // Joined field
	TranslatedTitle string    `json:"translated_title"`
}
