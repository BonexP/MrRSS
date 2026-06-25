package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// InitMinifluxSyncTable creates the miniflux_sync_queue table if it doesn't exist
func InitMinifluxSyncTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS miniflux_sync_queue (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		article_id INTEGER NOT NULL,
		article_url TEXT NOT NULL,
		sync_action TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		synced_at INTEGER,
		sync_error TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_miniflux_sync_article ON miniflux_sync_queue(article_id);
	CREATE INDEX IF NOT EXISTS idx_miniflux_sync_synced ON miniflux_sync_queue(synced_at);
	CREATE INDEX IF NOT EXISTS idx_miniflux_sync_url ON miniflux_sync_queue(article_url);
	`

	_, err := db.Exec(query)
	return err
}

// EnqueueMinifluxSyncChange adds a state change to the Miniflux sync queue
func (db *DB) EnqueueMinifluxSyncChange(articleID int64, articleURL string, action SyncAction) error {
	db.WaitForReady()

	query := `
	INSERT INTO miniflux_sync_queue (article_id, article_url, sync_action, created_at)
	VALUES (?, ?, ?, ?)
	`

	result, err := db.Exec(query, articleID, articleURL, string(action), time.Now().Unix())
	if err != nil {
		return fmt.Errorf("enqueue miniflux sync change: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("[Miniflux EnqueueSyncChange] Warning: could not get insert ID: %v", err)
	} else {
		log.Printf("[Miniflux EnqueueSyncChange] Successfully enqueued sync item ID=%d: articleID=%d url=%s action=%s",
			id, articleID, articleURL, action)
	}

	return nil
}

// GetPendingMinifluxSyncChanges retrieves all pending Miniflux sync changes
func (db *DB) GetPendingMinifluxSyncChanges(limit int) ([]SyncQueueItem, error) {
	db.WaitForReady()

	query := `
	SELECT id, article_id, article_url, sync_action, created_at, synced_at, sync_error
	FROM miniflux_sync_queue
	WHERE synced_at IS NULL
	ORDER BY created_at ASC
	LIMIT ?
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("get pending miniflux sync changes: %w", err)
	}
	defer rows.Close()

	var items []SyncQueueItem
	for rows.Next() {
		var item SyncQueueItem
		var syncedAt sql.NullInt64
		var syncError sql.NullString
		var action string
		var createdAt int64

		err := rows.Scan(
			&item.ID,
			&item.ArticleID,
			&item.ArticleURL,
			&action,
			&createdAt,
			&syncedAt,
			&syncError,
		)
		if err != nil {
			return nil, fmt.Errorf("scan miniflux sync queue item: %w", err)
		}

		item.Action = SyncAction(action)
		item.CreatedAt = time.Unix(createdAt, 0)

		if syncedAt.Valid {
			t := time.Unix(syncedAt.Int64, 0)
			item.SyncedAt = &t
		}

		if syncError.Valid {
			item.SyncError = &syncError.String
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate miniflux sync queue items: %w", err)
	}

	log.Printf("[Miniflux GetPendingSyncChanges] Retrieved %d pending items (limit=%d)", len(items), limit)

	return items, nil
}

// MarkMinifluxSynced marks Miniflux sync queue items as successfully synced
func (db *DB) MarkMinifluxSynced(itemIDs []int64) error {
	db.WaitForReady()

	if len(itemIDs) == 0 {
		return nil
	}

	query := `UPDATE miniflux_sync_queue SET synced_at = ? WHERE id = ?`
	now := time.Now().Unix()

	for _, id := range itemIDs {
		_, err := db.Exec(query, now, id)
		if err != nil {
			return fmt.Errorf("mark miniflux synced item %d: %w", id, err)
		}
	}

	return nil
}

// MarkMinifluxSyncFailed marks a Miniflux sync queue item as failed
func (db *DB) MarkMinifluxSyncFailed(itemID int64, errMsg string) error {
	db.WaitForReady()

	query := `UPDATE miniflux_sync_queue SET sync_error = ? WHERE id = ?`

	_, err := db.Exec(query, errMsg, itemID)
	if err != nil {
		return fmt.Errorf("mark miniflux sync failed: %w", err)
	}

	return nil
}

// ClearPendingMinifluxSyncForArticle removes all pending Miniflux sync changes for an article
func (db *DB) ClearPendingMinifluxSyncForArticle(articleID int64) error {
	db.WaitForReady()

	query := `DELETE FROM miniflux_sync_queue WHERE article_id = ? AND synced_at IS NULL`

	_, err := db.Exec(query, articleID)
	if err != nil {
		return fmt.Errorf("clear pending miniflux sync for article: %w", err)
	}

	return nil
}

// GetPendingMinifluxSyncCount returns the count of pending Miniflux sync changes
func (db *DB) GetPendingMinifluxSyncCount() (int, error) {
	db.WaitForReady()

	var count int
	query := `SELECT COUNT(*) FROM miniflux_sync_queue WHERE synced_at IS NULL`

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get pending miniflux sync count: %w", err)
	}

	return count, nil
}

// DeleteOldMinifluxSyncedItems removes old successfully synced Miniflux items
func (db *DB) DeleteOldMinifluxSyncedItems(olderThan time.Duration) error {
	db.WaitForReady()

	cutoff := time.Now().Add(-olderThan).Unix()
	query := `DELETE FROM miniflux_sync_queue WHERE synced_at IS NOT NULL AND synced_at < ?`

	_, err := db.Exec(query, cutoff)
	if err != nil {
		return fmt.Errorf("delete old miniflux synced items: %w", err)
	}

	return nil
}

// GetFailedMinifluxSyncItems returns Miniflux sync items that failed
func (db *DB) GetFailedMinifluxSyncItems(limit int) ([]SyncQueueItem, error) {
	db.WaitForReady()

	query := `
	SELECT id, article_id, article_url, sync_action, created_at, synced_at, sync_error
	FROM miniflux_sync_queue
	WHERE sync_error IS NOT NULL
	ORDER BY created_at DESC
	LIMIT ?
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("get failed miniflux sync items: %w", err)
	}
	defer rows.Close()

	var items []SyncQueueItem
	for rows.Next() {
		var item SyncQueueItem
		var syncedAt sql.NullInt64
		var syncError sql.NullString
		var action string
		var createdAt int64

		err := rows.Scan(
			&item.ID,
			&item.ArticleID,
			&item.ArticleURL,
			&action,
			&createdAt,
			&syncedAt,
			&syncError,
		)
		if err != nil {
			return nil, fmt.Errorf("scan miniflux sync queue item: %w", err)
		}

		item.Action = SyncAction(action)
		item.CreatedAt = time.Unix(createdAt, 0)

		if syncedAt.Valid {
			t := time.Unix(syncedAt.Int64, 0)
			item.SyncedAt = &t
		}

		if syncError.Valid {
			item.SyncError = &syncError.String
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate miniflux sync queue items: %w", err)
	}

	return items, nil
}
