package main

import (
	"context"
	"database/sql"

	"MrRSS/backend"
	"MrRSS/backend/feed"
	"MrRSS/backend/history"
)

// App struct
type App struct {
	ctx context.Context
	db  *sql.DB
}

// NewApp creates a new App application struct
func NewApp(db *sql.DB) *App {
	return &App{
		db: db,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) InitDatabase() {
	feeds := []backend.FeedsInfo{
		{Link: "https://www.kawabangga.com/feed", Category: "RSS/Atom"},
		{Link: "https://jvns.ca/atom.xml", Category: "RSS/Atom"},
		{Link: "https://www.ruanyifeng.com/blog/atom.xml", Category: "RSS/Atom"},
		{Link: "https://www.appinn.com/feed/", Category: "RSS/Atom"},
	}
	feed.SetFeedList(a.db, feeds)
}

func (a *App) GetFeedList() []backend.FeedsInfo {
	return feed.GetFeedList(a.db)
}

func (a *App) SetFeedList(feeds []backend.FeedsInfo) {
	feed.SetFeedList(a.db, feeds)
}

func (a *App) DeleteFeedList(feeds backend.FeedsInfo) {
	feed.DeleteFeedList(a.db, feeds)
}

func (a *App) GetFeedContent() []backend.FeedContentsInfo {
	return feed.GetFeedContent(a.db)
}

func (a *App) GetHistory() []backend.FeedContentsInfo {
	return history.GetHistory(a.db)
}

func (a *App) SetHistory(historys []backend.FeedContentsInfo) {
	history.SetHistory(a.db, historys)
}

func (a *App) SetHistoryReaded(feed backend.FeedContentsInfo) {
	history.SetHistoryReaded(a.db, feed)
}

func (a *App) ClearHistory() {
	history.ClearHistory(a.db)
}
