package miniflux

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// Logger is the package-level logger for Miniflux debugging.
// It writes to both stderr and a log file in the system temp directory.
var Logger *log.Logger

func init() {
	logDir := os.TempDir()
	logPath := filepath.Join(logDir, "mrrss_miniflux_debug.log")

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		Logger = log.New(os.Stderr, "[Miniflux] ", log.LstdFlags)
		return
	}

	Logger = log.New(io.MultiWriter(os.Stderr, f), "[Miniflux] ", log.LstdFlags)
	Logger.Printf("Miniflux debug log initialized: %s", logPath)
}
