package update

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"MrRSS/internal/handlers/core"
	"MrRSS/internal/utils"
	"MrRSS/internal/version"
)

// HandleCheckUpdates checks for the latest version on GitHub.
func HandleCheckUpdates(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentVersion := version.Version
	const githubAPI = "https://api.github.com/repos/WCY-dt/MrRSS/releases/latest"

	// Create HTTP client with global proxy support
	var proxyURL string
	proxyEnabled, _ := h.DB.GetSetting("proxy_enabled")
	if proxyEnabled == "true" {
		// Build proxy URL from global settings (use encrypted methods for credentials)
		proxyType, _ := h.DB.GetSetting("proxy_type")
		proxyHost, _ := h.DB.GetSetting("proxy_host")
		proxyPort, _ := h.DB.GetSetting("proxy_port")
		proxyUsername, _ := h.DB.GetEncryptedSetting("proxy_username")
		proxyPassword, _ := h.DB.GetEncryptedSetting("proxy_password")
		proxyURL = utils.BuildProxyURL(proxyType, proxyHost, proxyPort, proxyUsername, proxyPassword)
	}

	client, err := utils.CreateHTTPClient(proxyURL, 30*time.Second)
	if err != nil {
		log.Printf("Error creating HTTP client: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current_version": currentVersion,
			"error":           "Failed to create HTTP client",
		})
		return
	}

	resp, err := client.Get(githubAPI)
	if err != nil {
		log.Printf("Error checking for updates: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current_version": currentVersion,
			"error":           "Failed to check for updates",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("GitHub API returned status: %d", resp.StatusCode)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current_version": currentVersion,
			"error":           "Failed to fetch latest release",
		})
		return
	}

	var release struct {
		TagName     string `json:"tag_name"`
		Name        string `json:"name"`
		HTMLURL     string `json:"html_url"`
		Body        string `json:"body"`
		PublishedAt string `json:"published_at"`
		Assets      []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Size               int64  `json:"size"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		log.Printf("Error decoding release info: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current_version": currentVersion,
			"error":           "Failed to parse release information",
		})
		return
	}

	// Remove 'v' prefix if present for comparison
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	hasUpdate := compareVersions(latestVersion, currentVersion) > 0

	// Find the appropriate download URL based on platform
	var downloadURL string
	var assetName string
	var assetSize int64
	platform := runtime.GOOS
	arch := runtime.GOARCH

	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)

		// Match platform-specific installer/package with architecture
		// Asset naming convention: MrRSS-{version}-{platform}-{arch}-installer.{ext}
		platformArch := platform + "-" + arch

		if platform == "windows" {
			// For Windows, prefer installer.exe, fallback to .zip
			if strings.Contains(name, platformArch) && strings.HasSuffix(name, "-installer.exe") {
				downloadURL = asset.BrowserDownloadURL
				assetName = asset.Name
				assetSize = asset.Size
				break
			}
		} else if platform == "linux" {
			// For Linux, prefer .AppImage, fallback to .tar.gz
			if strings.Contains(name, platformArch) && strings.HasSuffix(name, ".appimage") {
				downloadURL = asset.BrowserDownloadURL
				assetName = asset.Name
				assetSize = asset.Size
				break
			}
		} else if platform == "darwin" {
			// For macOS, use universal build (supports both arm64 and amd64)
			if strings.Contains(name, "darwin-universal") && strings.HasSuffix(name, ".dmg") {
				downloadURL = asset.BrowserDownloadURL
				assetName = asset.Name
				assetSize = asset.Size
				break
			}
		}
	}

	response := map[string]interface{}{
		"current_version": currentVersion,
		"latest_version":  latestVersion,
		"has_update":      hasUpdate,
		"platform":        platform,
		"arch":            arch,
	}

	if downloadURL != "" {
		response["download_url"] = downloadURL
		response["asset_name"] = assetName
		response["asset_size"] = assetSize
	}

	json.NewEncoder(w).Encode(response)
}
