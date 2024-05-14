package history

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"

	"MrRSS/backend"
)

func SetHistory(db *sql.DB, history []backend.FeedContentsInfo) {
	var err error
	// Clear the History table before inserting new data
	_, err = db.Exec("DELETE FROM [History]")
	if err != nil {
		log.Fatal(err)
	}

	// Insert history into the History table
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO [History]([FeedTitle], [FeedImage], [Title], [Link], [TimeSince], [Time], [Image], [Content], [Readed]) values(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, item := range history {
		_, err = stmt.Exec(item.FeedTitle, item.FeedImage, item.Title, item.Link, item.TimeSince, item.Time, item.Image, item.Content, item.Readed)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func SetHistoryReaded(db *sql.DB, history backend.FeedContentsInfo) {
	// Update the Readed field in the History table
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("UPDATE [History] SET [Readed] = ? WHERE [Link] = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(history.Readed, history.Link)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func ClearHistory(db *sql.DB) {
	// Clear the History table
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("DELETE FROM [History]")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}
