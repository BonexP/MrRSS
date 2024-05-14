package backend

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
)

type FeedsInfo struct {
	Link     string
	Category string
}

type FeedContentsInfo struct {
	FeedTitle string
	FeedImage string
	Title     string
	Link      string
	TimeSince string
	Time      string
	Image     string
	Content   string
	Readed    bool
}

func InitDatabase(db *sql.DB) {
	var err error
	var sqlStmt string

	// Check if the feeds table exists
	sqlStmt = `
		CREATE TABLE IF NOT EXISTS [Feeds] ([Link] VARCHAR NOT NULL PRIMARY KEY, [Category] VARCHAR);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the history table exists
	sqlStmt = `
		CREATE TABLE IF NOT EXISTS [History] ([FeedTitle] VARCHAR, [FeedImage] VARCHAR, [Title] VARCHAR, [Link] VARCHAR NOT NULL PRIMARY KEY, [TimeSince] VARCHAR, [Time] VARCHAR, [Image] VARCHAR, [Content] TEXT, [Readed] BOOLEAN);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDbFilePath(dbName string) string {
	var dbFilePath string
	if os.Getenv("DEV_MODE") == "true" {
		dbFilePath = fmt.Sprintf("data/%s", dbName)
	} else {
		configDir, _ := os.UserConfigDir()
		dbFilePath = filepath.Join(configDir, "MrRSS", "database", dbName)
	}
	return dbFilePath
}
