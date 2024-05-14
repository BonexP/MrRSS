package main

import (
	"MrRSS/backend"
	"database/sql"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Get the database file path
	dbFilePath := backend.GetDbFilePath("rss.db")
	db, dbErr := sql.Open("sqlite", dbFilePath)
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	// Ensure the database connection is closed when the function ends
	defer db.Close()

	backend.InitDatabase(db)

	// Create an instance of the app structure
	app := NewApp(db)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "MrRSS",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
