package database

import (
	"database/sql"
	"time"

	"MrRSS/internal/models"
)

// AddFeed adds a new feed or updates an existing one.
// Returns the feed ID and any error encountered.
func (db *DB) AddFeed(feed *models.Feed) (int64, error) {
	db.WaitForReady()

	// Check if feed already exists
	var existingID int64
	err := db.QueryRow("SELECT id FROM feeds WHERE url = ?", feed.URL).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Feed doesn't exist, insert new
		query := `INSERT INTO feeds (title, url, link, description, category, image_url, script_path, hide_from_timeline, proxy_url, proxy_enabled, refresh_interval, last_updated) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		result, err := db.Exec(query, feed.Title, feed.URL, feed.Link, feed.Description, feed.Category, feed.ImageURL, feed.ScriptPath, feed.HideFromTimeline, feed.ProxyURL, feed.ProxyEnabled, feed.RefreshInterval, time.Now())
		if err != nil {
			return 0, err
		}
		newID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		return newID, nil
	} else if err != nil {
		return 0, err
	}

	// Feed exists, update it
	query := `UPDATE feeds SET title = ?, link = ?, description = ?, category = ?, image_url = ?, script_path = ?, hide_from_timeline = ?, proxy_url = ?, proxy_enabled = ?, refresh_interval = ?, last_updated = ? WHERE id = ?`
	_, err = db.Exec(query, feed.Title, feed.Link, feed.Description, feed.Category, feed.ImageURL, feed.ScriptPath, feed.HideFromTimeline, feed.ProxyURL, feed.ProxyEnabled, feed.RefreshInterval, time.Now(), existingID)
	return existingID, err
}

// DeleteFeed deletes a feed and all its articles.
func (db *DB) DeleteFeed(id int64) error {
	db.WaitForReady()
	// First delete associated articles
	_, err := db.Exec("DELETE FROM articles WHERE feed_id = ?", id)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM feeds WHERE id = ?", id)
	return err
}

// GetFeeds returns all feeds.
func (db *DB) GetFeeds() ([]models.Feed, error) {
	db.WaitForReady()
	rows, err := db.Query("SELECT id, title, url, link, description, category, image_url, last_updated, last_error, COALESCE(discovery_completed, 0), COALESCE(script_path, ''), COALESCE(hide_from_timeline, 0), COALESCE(proxy_url, ''), COALESCE(proxy_enabled, 0), COALESCE(refresh_interval, 0) FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []models.Feed
	for rows.Next() {
		var f models.Feed
		var link, category, imageURL, lastError, scriptPath, proxyURL sql.NullString
		if err := rows.Scan(&f.ID, &f.Title, &f.URL, &link, &f.Description, &category, &imageURL, &f.LastUpdated, &lastError, &f.DiscoveryCompleted, &scriptPath, &f.HideFromTimeline, &proxyURL, &f.ProxyEnabled, &f.RefreshInterval); err != nil {
			return nil, err
		}
		f.Link = link.String
		f.Category = category.String
		f.ImageURL = imageURL.String
		f.LastError = lastError.String
		f.ScriptPath = scriptPath.String
		f.ProxyURL = proxyURL.String
		feeds = append(feeds, f)
	}
	return feeds, nil
}

// GetFeedByID retrieves a specific feed by its ID.
func (db *DB) GetFeedByID(id int64) (*models.Feed, error) {
	db.WaitForReady()
	row := db.QueryRow("SELECT id, title, url, link, description, category, image_url, last_updated, last_error, COALESCE(discovery_completed, 0), COALESCE(script_path, ''), COALESCE(hide_from_timeline, 0), COALESCE(proxy_url, ''), COALESCE(proxy_enabled, 0), COALESCE(refresh_interval, 0) FROM feeds WHERE id = ?", id)

	var f models.Feed
	var link, category, imageURL, lastError, scriptPath, proxyURL sql.NullString
	if err := row.Scan(&f.ID, &f.Title, &f.URL, &link, &f.Description, &category, &imageURL, &f.LastUpdated, &lastError, &f.DiscoveryCompleted, &scriptPath, &f.HideFromTimeline, &proxyURL, &f.ProxyEnabled, &f.RefreshInterval); err != nil {
		return nil, err
	}
	f.Link = link.String
	f.Category = category.String
	f.ImageURL = imageURL.String
	f.LastError = lastError.String
	f.ScriptPath = scriptPath.String
	f.ProxyURL = proxyURL.String

	return &f, nil
}

// GetAllFeedURLs returns a set of all subscribed RSS feed URLs for deduplication.
func (db *DB) GetAllFeedURLs() (map[string]bool, error) {
	db.WaitForReady()
	rows, err := db.Query("SELECT url FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make(map[string]bool)
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls[url] = true
	}
	return urls, rows.Err()
}

// UpdateFeed updates feed title, URL, category, script_path, hide_from_timeline, proxy settings, and refresh_interval.
func (db *DB) UpdateFeed(id int64, title, url, category, scriptPath string, hideFromTimeline bool, proxyURL string, proxyEnabled bool, refreshInterval int) error {
	db.WaitForReady()
	_, err := db.Exec("UPDATE feeds SET title = ?, url = ?, category = ?, script_path = ?, hide_from_timeline = ?, proxy_url = ?, proxy_enabled = ?, refresh_interval = ? WHERE id = ?", title, url, category, scriptPath, hideFromTimeline, proxyURL, proxyEnabled, refreshInterval, id)
	return err
}

// UpdateFeedCategory updates a feed's category.
func (db *DB) UpdateFeedCategory(id int64, category string) error {
	db.WaitForReady()
	_, err := db.Exec("UPDATE feeds SET category = ? WHERE id = ?", category, id)
	return err
}

// UpdateFeedImage updates a feed's image URL.
func (db *DB) UpdateFeedImage(id int64, imageURL string) error {
	db.WaitForReady()
	_, err := db.Exec("UPDATE feeds SET image_url = ? WHERE id = ?", imageURL, id)
	return err
}

// UpdateFeedLink updates a feed's homepage link.
func (db *DB) UpdateFeedLink(id int64, link string) error {
	db.WaitForReady()
	_, err := db.Exec("UPDATE feeds SET link = ? WHERE id = ?", link, id)
	return err
}

// UpdateFeedError updates a feed's error message.
func (db *DB) UpdateFeedError(id int64, errorMsg string) error {
	db.WaitForReady()
	_, err := db.Exec("UPDATE feeds SET last_error = ? WHERE id = ?", errorMsg, id)
	return err
}

// MarkFeedDiscovered marks a feed as having completed discovery.
func (db *DB) MarkFeedDiscovered(id int64) error {
	db.WaitForReady()
	_, err := db.Exec("UPDATE feeds SET discovery_completed = 1 WHERE id = ?", id)
	return err
}
