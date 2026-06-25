package miniflux

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"MrRSS/internal/database"
	"MrRSS/internal/models"
)

const (
	entriesPerPage = 100
	maxSyncWorkers = 3
)

// SyncResult holds the result of a Miniflux sync operation
type SyncResult struct {
	PullSuccess      bool
	PullChangesCount int
	PushSuccess      bool
	PushChangesCount int
	Errors           []string
	Duration         time.Duration
	LastSyncTime     time.Time
}

// BidirectionalSyncService handles bidirectional sync with Miniflux
type BidirectionalSyncService struct {
	client *Client
	db     *database.DB
}

// NewBidirectionalSyncService creates a new Miniflux sync service
func NewBidirectionalSyncService(serverURL, apiKey string, db *database.DB) *BidirectionalSyncService {
	return &BidirectionalSyncService{
		client: NewClient(serverURL, apiKey),
		db:     db,
	}
}

// Sync performs a full bidirectional sync with Miniflux
func (s *BidirectionalSyncService) Sync(ctx context.Context) (*SyncResult, error) {
	start := time.Now()
	result := &SyncResult{LastSyncTime: time.Now()}

	log.Printf("[Miniflux Sync] Starting bidirectional sync")

	if err := s.pullFromServer(ctx, result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("pull error: %v", err))
		log.Printf("[Miniflux Sync] Pull failed: %v", err)
	} else {
		result.PullSuccess = true
		log.Printf("[Miniflux Sync] Pull completed: %d changes", result.PullChangesCount)
	}

	if err := s.pushToServer(ctx, result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("push error: %v", err))
		log.Printf("[Miniflux Sync] Push failed: %v", err)
	} else {
		result.PushSuccess = true
		log.Printf("[Miniflux Sync] Push completed: %d changes", result.PushChangesCount)
	}

	result.Duration = time.Since(start)
	log.Printf("[Miniflux Sync] Sync completed in %v: pull=%v(%d) push=%v(%d) errors=%d",
		result.Duration, result.PullSuccess, result.PullChangesCount,
		result.PushSuccess, result.PushChangesCount, len(result.Errors))

	return result, nil
}

// SyncFeed pulls changes for a single Miniflux feed
func (s *BidirectionalSyncService) SyncFeed(ctx context.Context, minifluxFeedID int64) (int, error) {
	log.Printf("[Miniflux Sync] Syncing feed %d", minifluxFeedID)

	feed, err := s.client.GetFeed(ctx, minifluxFeedID)
	if err != nil {
		return 0, fmt.Errorf("get miniflux feed %d: %w", minifluxFeedID, err)
	}

	localFeedID, err := s.ensureLocalFeed(ctx, feed)
	if err != nil {
		return 0, fmt.Errorf("ensure local feed: %w", err)
	}

	count, err := s.pullEntries(ctx, feed.ID, localFeedID, nil)
	if err != nil {
		return 0, fmt.Errorf("pull entries: %w", err)
	}

	return count, nil
}

// SyncArticleStatus pushes a single article's status change to Miniflux
func (s *BidirectionalSyncService) SyncArticleStatus(ctx context.Context, articleID int64, articleURL string, action database.SyncAction) error {
	article, err := s.db.GetArticleByURL(articleURL)
	if err != nil {
		return fmt.Errorf("get article by URL: %w", err)
	}
	if article == nil {
		return fmt.Errorf("article not found: %s", articleURL)
	}

	entryID := article.MinifluxEntryID
	if entryID == 0 {
		filter := EntryFilter{Limit: 1}
		resp, findErr := s.client.GetEntries(ctx, filter)
		if findErr != nil {
			return fmt.Errorf("find miniflux entry: %w", findErr)
		}
		if len(resp.Entries) == 0 {
			return fmt.Errorf("no miniflux entry found for URL: %s", articleURL)
		}
		entryID = resp.Entries[0].ID
		_ = s.db.UpdateMinifluxEntryID(articleID, entryID)
	}

	var status string
	var starred *bool
	switch action {
	case database.SyncActionMarkRead:
		status = "read"
	case database.SyncActionMarkUnread:
		status = "unread"
	case database.SyncActionStar:
		starred = boolPtr(true)
	case database.SyncActionUnstar:
		starred = boolPtr(false)
	}

	if err := s.client.UpdateEntries(ctx, []int64{entryID}, status, starred); err != nil {
		if enqueueErr := s.db.EnqueueMinifluxSyncChange(articleID, articleURL, action); enqueueErr != nil {
			log.Printf("[Miniflux Sync] Error enqueuing failed sync: %v", enqueueErr)
		}
		return fmt.Errorf("update entry %d: %w", entryID, err)
	}

	_ = s.db.ClearPendingMinifluxSyncForArticle(articleID)
	return nil
}

// pullFromServer pulls feeds and entries from Miniflux
func (s *BidirectionalSyncService) pullFromServer(ctx context.Context, result *SyncResult) error {
	log.Printf("[Miniflux Sync] Pulling feeds from Miniflux")

	feeds, err := s.client.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("get miniflux feeds: %w", err)
	}

	log.Printf("[Miniflux Sync] Found %d feeds on Miniflux server", len(feeds))

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxSyncWorkers)
	var mu sync.Mutex
	var pullErrors []string

	for _, feed := range feeds {
		wg.Add(1)
		sem <- struct{}{}
		go func(f Feed) {
			defer wg.Done()
			defer func() { <-sem }()

			localFeedID, err := s.ensureLocalFeed(ctx, &f)
			if err != nil {
				errMsg := fmt.Sprintf("feed %d (%s): %v", f.ID, f.Title, err)
				log.Printf("[Miniflux Sync] Error ensuring feed: %s", errMsg)
				mu.Lock()
				pullErrors = append(pullErrors, errMsg)
				mu.Unlock()
				return
			}

			count, err := s.pullEntries(ctx, f.ID, localFeedID, nil)
			if err != nil {
				errMsg := fmt.Sprintf("pull entries for feed %d: %v", f.ID, err)
				log.Printf("[Miniflux Sync] %s", errMsg)
				mu.Lock()
				pullErrors = append(pullErrors, errMsg)
				mu.Unlock()
				return
			}

			mu.Lock()
			result.PullChangesCount += count
			mu.Unlock()
		}(feed)
	}

	wg.Wait()
	s.removeDeletedFeeds(ctx, feeds)

	if len(pullErrors) > 0 {
		return fmt.Errorf("pull errors: %v", pullErrors)
	}

	return nil
}

// pushToServer pushes local status changes to Miniflux
func (s *BidirectionalSyncService) pushToServer(ctx context.Context, result *SyncResult) error {
	log.Printf("[Miniflux Sync] Pushing changes to Miniflux")

	failedItems, err := s.db.GetFailedMinifluxSyncItems(50)
	if err != nil {
		log.Printf("[Miniflux Sync] Error getting failed items: %v", err)
	} else {
		for _, item := range failedItems {
			if err := s.SyncArticleStatus(ctx, item.ArticleID, item.ArticleURL, item.Action); err != nil {
				log.Printf("[Miniflux Sync] Retry failed for article %d: %v", item.ArticleID, err)
			}
		}
	}

	pendingItems, err := s.db.GetPendingMinifluxSyncChanges(100)
	if err != nil {
		return fmt.Errorf("get pending changes: %w", err)
	}

	if len(pendingItems) == 0 {
		log.Printf("[Miniflux Sync] No pending changes to push")
		return nil
	}

	log.Printf("[Miniflux Sync] Processing %d pending changes", len(pendingItems))

	batchErrors := 0
	var syncedIDs []int64

	for _, item := range pendingItems {
		article, err := s.db.GetArticleByURL(item.ArticleURL)
		if err != nil || article == nil {
			log.Printf("[Miniflux Sync] Article not found for URL %s: %v", item.ArticleURL, err)
			_ = s.db.MarkMinifluxSyncFailed(item.ID, fmt.Sprintf("article not found: %v", err))
			batchErrors++
			continue
		}

		entryID := article.MinifluxEntryID
		if entryID == 0 {
			log.Printf("[Miniflux Sync] No Miniflux entry ID for article %d, skipping", article.ID)
			_ = s.db.MarkMinifluxSyncFailed(item.ID, "no entry ID")
			batchErrors++
			continue
		}

		var status string
		var starred *bool
		switch item.Action {
		case database.SyncActionMarkRead:
			status = "read"
		case database.SyncActionMarkUnread:
			status = "unread"
		case database.SyncActionStar:
			starred = boolPtr(true)
		case database.SyncActionUnstar:
			starred = boolPtr(false)
		}

		if err := s.client.UpdateEntries(ctx, []int64{entryID}, status, starred); err != nil {
			log.Printf("[Miniflux Sync] Failed to update entry %d: %v", entryID, err)
			_ = s.db.MarkMinifluxSyncFailed(item.ID, err.Error())
			batchErrors++
			continue
		}

		syncedIDs = append(syncedIDs, item.ID)
		result.PushChangesCount++
	}

	if len(syncedIDs) > 0 {
		_ = s.db.MarkMinifluxSynced(syncedIDs)
	}

	if batchErrors > 0 {
		return fmt.Errorf("%d push errors", batchErrors)
	}

	return nil
}

// ensureLocalFeed finds or creates a local feed for a Miniflux feed
func (s *BidirectionalSyncService) ensureLocalFeed(ctx context.Context, feed *Feed) (int64, error) {
	existing, err := s.db.GetFeedByMinifluxID(feed.ID)
	if err == nil && existing != nil {
		return existing.ID, nil
	}

	category := ""
	if feed.Category != nil {
		category = feed.Category.Title
	}

	siteURL := feed.SiteURL
	if siteURL == "" {
		siteURL = feed.FeedURL
	}

	localFeed := &models.Feed{
		Title:            feed.Title,
		URL:              feed.FeedURL,
		Link:             siteURL,
		Category:         category,
		IsMinifluxSource: true,
		MinifluxFeedID:   feed.ID,
	}

	newID, err := s.db.AddFeed(localFeed)
	if err != nil {
		return 0, fmt.Errorf("add feed: %w", err)
	}

	log.Printf("[Miniflux Sync] Created local feed '%s' (ID: %d) for Miniflux feed %d", feed.Title, newID, feed.ID)
	return newID, nil
}

// pullEntries pulls entries from a Miniflux feed and saves them locally
func (s *BidirectionalSyncService) pullEntries(ctx context.Context, minifluxFeedID, localFeedID int64, afterEntryID *int64) (int, error) {
	filter := EntryFilter{
		FeedID:    minifluxFeedID,
		Limit:     entriesPerPage,
		Direction: "asc",
	}

	if afterEntryID != nil {
		filter.AfterEntryID = *afterEntryID
	}

	var allEntries []Entry
	offset := 0

	for {
		filter.Offset = offset
		resp, err := s.client.GetEntries(ctx, filter)
		if err != nil {
			return 0, fmt.Errorf("get entries at offset %d: %w", offset, err)
		}

		if len(resp.Entries) == 0 {
			break
		}

		allEntries = append(allEntries, resp.Entries...)
		offset += len(resp.Entries)

		if offset >= resp.Total {
			break
		}
	}

	if len(allEntries) == 0 {
		return 0, nil
	}

	articles := make([]*models.Article, 0, len(allEntries))
	for _, entry := range allEntries {
		article := convertEntryToArticle(entry, localFeedID)
		articles = append(articles, article)
	}

	if err := s.db.SaveArticles(ctx, articles); err != nil {
		return 0, fmt.Errorf("save articles: %w", err)
	}

	syncStatusBatch := make([]Entry, 0, len(allEntries))
	for _, entry := range allEntries {
		uniqueID := generateUniqueID(entry.Title, localFeedID, entry.PublishedAt)
		var articleID int64
		err := s.db.QueryRow("SELECT id FROM articles WHERE unique_id = ? AND feed_id = ?", uniqueID, localFeedID).Scan(&articleID)
		if err == nil {
			_ = s.db.UpdateMinifluxEntryID(articleID, entry.ID)
		}
		syncStatusBatch = append(syncStatusBatch, entry)
	}

	s.syncArticleStatus(ctx, syncStatusBatch, localFeedID)

	log.Printf("[Miniflux Sync] Pulled %d entries for feed %d", len(allEntries), minifluxFeedID)
	return len(allEntries), nil
}

// syncArticleStatus syncs read/starred status from Miniflux entries to local articles
func (s *BidirectionalSyncService) syncArticleStatus(ctx context.Context, entries []Entry, localFeedID int64) {
	for _, entry := range entries {
		uniqueID := generateUniqueID(entry.Title, localFeedID, entry.PublishedAt)
		isRead := entry.Status == "read"
		isStarred := entry.Starred
		_, _ = s.db.Exec("UPDATE articles SET is_read = ?, is_favorite = ? WHERE unique_id = ? AND feed_id = ?",
			isRead, isStarred, uniqueID, localFeedID)
	}
}

// removeDeletedFeeds removes local feeds that no longer exist on Miniflux server
func (s *BidirectionalSyncService) removeDeletedFeeds(ctx context.Context, remoteFeeds []Feed) {
	remoteFeedIDs := make(map[int64]bool)
	for _, f := range remoteFeeds {
		remoteFeedIDs[f.ID] = true
	}

	rows, err := s.db.Query("SELECT id, miniflux_feed_id FROM feeds WHERE COALESCE(is_miniflux_source, 0) = 1")
	if err != nil {
		log.Printf("[Miniflux Sync] Error querying local feeds: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var localID, minifluxID int64
		if err := rows.Scan(&localID, &minifluxID); err != nil {
			continue
		}
		if !remoteFeedIDs[minifluxID] {
			log.Printf("[Miniflux Sync] Removing feed %d (Miniflux ID %d) - no longer on server", localID, minifluxID)
			_ = s.db.DeleteFeed(localID)
		}
	}
}

// convertEntryToArticle converts a Miniflux entry to a local article model
func convertEntryToArticle(entry Entry, feedID int64) *models.Article {
	content := ""
	if entry.Content != nil {
		content = entry.Content.Content
	}

	author := ""
	if entry.Author != nil {
		author = entry.Author.Name
	}

	return &models.Article{
		FeedID:          feedID,
		Title:           entry.Title,
		URL:             entry.URL,
		PublishedAt:     entry.PublishedAt,
		IsRead:          entry.Status == "read",
		IsFavorite:      entry.Starred,
		Author:          author,
		Summary:         content,
		MinifluxEntryID: entry.ID,
	}
}

func generateUniqueID(title string, feedID int64, publishedAt time.Time) string {
	return fmt.Sprintf("%d-%s-%d", feedID, title, publishedAt.Unix())
}

func boolPtr(b bool) *bool {
	return &b
}
