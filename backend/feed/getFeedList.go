package feed

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"

	"MrRSS/backend"
)

func GetFeedList(db *sql.DB) []backend.FeedsInfo {
	result := []backend.FeedsInfo{}

	// Print the feeds in the Feeds table
	rows, err := db.Query("SELECT [Link], [Category] FROM [Feeds]")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var link string
		var category string
		err = rows.Scan(&link, &category)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, backend.FeedsInfo{Link: link, Category: category})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}
