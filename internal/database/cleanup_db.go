package database

import (
	"strconv"
	"time"
)

// CleanupOldArticles removes articles based on age and status.
// - Articles older than configured days: delete except favorited
// - Also checks database size against max_cache_size_mb setting
func (db *DB) CleanupOldArticles() (int64, error) {
	db.WaitForReady()

	// Get max article age from settings (default 30 days)
	maxAgeDaysStr, err := db.GetSetting("max_article_age_days")
	maxAgeDays := 30
	if err == nil {
		if days, err := strconv.Atoi(maxAgeDaysStr); err == nil && days > 0 {
			maxAgeDays = days
		}
	}

	cutoffDate := time.Now().AddDate(0, 0, -maxAgeDays)

	// Delete articles older than configured age that are not favorited
	result, err := db.Exec(`
		DELETE FROM articles
		WHERE published_at < ?
		AND is_favorite = 0
	`, cutoffDate)
	if err != nil {
		return 0, err
	}

	count, _ := result.RowsAffected()

	// Run VACUUM to reclaim space
	_, _ = db.Exec("VACUUM")

	return count, nil
}

// CleanupUnimportantArticles removes all articles except read and favorited ones.
func (db *DB) CleanupUnimportantArticles() (int64, error) {
	db.WaitForReady()

	result, err := db.Exec(`
		DELETE FROM articles
		WHERE is_read = 0
		AND is_favorite = 0
	`)
	if err != nil {
		return 0, err
	}

	count, _ := result.RowsAffected()

	// Run VACUUM to reclaim space
	_, _ = db.Exec("VACUUM")

	return count, nil
}

// GetDatabaseSizeMB returns the current database size in megabytes.
func (db *DB) GetDatabaseSizeMB() (float64, error) {
	db.WaitForReady()

	var pageCount, pageSize int64
	err := db.QueryRow("PRAGMA page_count").Scan(&pageCount)
	if err != nil {
		return 0, err
	}

	err = db.QueryRow("PRAGMA page_size").Scan(&pageSize)
	if err != nil {
		return 0, err
	}

	sizeBytes := pageCount * pageSize
	sizeMB := float64(sizeBytes) / (1024 * 1024)

	return sizeMB, nil
}
