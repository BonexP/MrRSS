package feed

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"

	"MrRSS/backend"
)

func SetFeedList(db *sql.DB, feeds []backend.FeedsInfo) {
	// Insert feeds into the Feeds table
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO [Feeds]([Link], [Category]) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, feed := range feeds {
		_, err = stmt.Exec(feed.Link, feed.Category)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func DeleteFeedList(db *sql.DB, feed backend.FeedsInfo) {
	// Delete feeds from the Feeds table
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("DELETE FROM [Feeds] WHERE [Link] = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(feed.Link)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}
