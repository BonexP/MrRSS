package history

import (
	"database/sql"
	"log"
	"sort"

	_ "github.com/glebarez/go-sqlite"

	"MrRSS/backend"
)

func GetHistory(db *sql.DB) []backend.FeedContentsInfo {
	result := []backend.FeedContentsInfo{}

	// Get the items in the History table
	rows, err := db.Query("SELECT [FeedTitle], [FeedImage], [Title], [Link], [TimeSince], [Time], [Image], [Content], [Readed] FROM [History]")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var feedTitle string
		var feedImage string
		var title string
		var link string
		var timeSince string
		var time string
		var image string
		var content string
		var readed bool
		err = rows.Scan(&feedTitle, &feedImage, &title, &link, &timeSince, &time, &image, &content, &readed)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, backend.FeedContentsInfo{FeedTitle: feedTitle, FeedImage: feedImage, Title: title, Link: link, TimeSince: timeSince, Time: time, Image: image, Content: content, Readed: readed})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// Sort the result by time
	sort.Slice(result, func(i, j int) bool {
		return result[i].Time > result[j].Time
	})

	return result
}

func CheckInHistory(db *sql.DB, feed backend.FeedContentsInfo) bool {
	// Check if the item is in the History table
	rows, err := db.Query("SELECT [Link] FROM [History] WHERE [Link] = ?", feed.Link)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return rows.Next()
}

func GetHistoryReaded(db *sql.DB, feed backend.FeedContentsInfo) bool {
	// Get the Readed field in the History table
	rows, err := db.Query("SELECT [Readed] FROM [History] WHERE [Link] = ?", feed.Link)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var readed bool
	if rows.Next() {
		err = rows.Scan(&readed)
		if err != nil {
			log.Fatal(err)
		}
	}

	return readed
}
