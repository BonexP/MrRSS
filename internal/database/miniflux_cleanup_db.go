package database

import (
	"log"
)

// CleanupMinifluxData removes all Miniflux-related feeds, articles, and sync queue items
func (db *DB) CleanupMinifluxData() error {
	db.WaitForReady()

	log.Printf("[Miniflux Cleanup] Starting cleanup of Miniflux data...")

	feedIDs, err := db.getMinifluxFeedIDs()
	if err != nil {
		log.Printf("[Miniflux Cleanup] Error getting Miniflux feed IDs: %v", err)
		return err
	}

	if len(feedIDs) == 0 {
		log.Printf("[Miniflux Cleanup] No Miniflux feeds found, nothing to clean")
		return nil
	}

	log.Printf("[Miniflux Cleanup] Found %d Miniflux feeds to delete", len(feedIDs))

	_, err = db.Exec(`DELETE FROM article_contents WHERE article_id IN (
		SELECT id FROM articles WHERE feed_id IN (
			SELECT id FROM feeds WHERE is_miniflux_source = 1
		)
	)`)
	if err != nil {
		log.Printf("[Miniflux Cleanup] Error deleting article contents: %v", err)
	} else {
		log.Printf("[Miniflux Cleanup] Deleted article contents for Miniflux feeds")
	}

	result, err := db.Exec(`DELETE FROM articles WHERE feed_id IN (
		SELECT id FROM feeds WHERE is_miniflux_source = 1
	)`)
	if err != nil {
		log.Printf("[Miniflux Cleanup] Error deleting articles: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("[Miniflux Cleanup] Deleted %d articles from Miniflux feeds", rowsAffected)

	result, err = db.Exec("DELETE FROM feeds WHERE is_miniflux_source = 1")
	if err != nil {
		log.Printf("[Miniflux Cleanup] Error deleting feeds: %v", err)
		return err
	}

	rowsAffected, _ = result.RowsAffected()
	log.Printf("[Miniflux Cleanup] Deleted %d Miniflux feeds", rowsAffected)

	_, err = db.Exec("DELETE FROM miniflux_sync_queue")
	if err != nil {
		log.Printf("[Miniflux Cleanup] Error clearing sync queue: %v", err)
	} else {
		log.Printf("[Miniflux Cleanup] Cleared Miniflux sync queue")
	}

	log.Printf("[Miniflux Cleanup] Completed successfully")

	return nil
}

func (db *DB) getMinifluxFeedIDs() ([]int64, error) {
	rows, err := db.Query("SELECT id FROM feeds WHERE is_miniflux_source = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		feedIDs = append(feedIDs, id)
	}

	return feedIDs, rows.Err()
}
