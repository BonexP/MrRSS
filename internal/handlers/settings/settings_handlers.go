package settings

import (
	"encoding/json"
	"log"
	"net/http"

	"MrRSS/internal/handlers/core"
	"MrRSS/internal/utils"
)

// HandleSettings handles GET and POST requests for application settings.
func HandleSettings(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		interval, _ := h.DB.GetSetting("update_interval")
		refreshMode, _ := h.DB.GetSetting("refresh_mode")
		translationEnabled, _ := h.DB.GetSetting("translation_enabled")
		targetLang, _ := h.DB.GetSetting("target_language")
		provider, _ := h.DB.GetSetting("translation_provider")
		apiKey, _ := h.DB.GetEncryptedSetting("deepl_api_key")
		baiduAppID, _ := h.DB.GetSetting("baidu_app_id")
		baiduSecretKey, _ := h.DB.GetEncryptedSetting("baidu_secret_key")
		aiAPIKey, _ := h.DB.GetEncryptedSetting("ai_api_key")
		aiEndpoint, _ := h.DB.GetSetting("ai_endpoint")
		aiModel, _ := h.DB.GetSetting("ai_model")
		aiSystemPrompt, _ := h.DB.GetSetting("ai_system_prompt")
		autoCleanup, _ := h.DB.GetSetting("auto_cleanup_enabled")
		maxCacheSize, _ := h.DB.GetSetting("max_cache_size_mb")
		maxArticleAge, _ := h.DB.GetSetting("max_article_age_days")
		language, _ := h.DB.GetSetting("language")
		theme, _ := h.DB.GetSetting("theme")
		lastUpdate, _ := h.DB.GetSetting("last_article_update")
		showHidden, _ := h.DB.GetSetting("show_hidden_articles")
		startupOnBoot, _ := h.DB.GetSetting("startup_on_boot")
		closeToTray, _ := h.DB.GetSetting("close_to_tray")
		shortcuts, _ := h.DB.GetSetting("shortcuts")
		rules, _ := h.DB.GetSetting("rules")
		defaultViewMode, _ := h.DB.GetSetting("default_view_mode")
		mediaCacheEnabled, _ := h.DB.GetSetting("media_cache_enabled")
		mediaCacheMaxSizeMB, _ := h.DB.GetSetting("media_cache_max_size_mb")
		mediaCacheMaxAgeDays, _ := h.DB.GetSetting("media_cache_max_age_days")
		summaryEnabled, _ := h.DB.GetSetting("summary_enabled")
		summaryLength, _ := h.DB.GetSetting("summary_length")
		summaryProvider, _ := h.DB.GetSetting("summary_provider")
		summaryAIAPIKey, _ := h.DB.GetEncryptedSetting("summary_ai_api_key")
		summaryAIEndpoint, _ := h.DB.GetSetting("summary_ai_endpoint")
		summaryAIModel, _ := h.DB.GetSetting("summary_ai_model")
		summaryAISystemPrompt, _ := h.DB.GetSetting("summary_ai_system_prompt")
		proxyEnabled, _ := h.DB.GetSetting("proxy_enabled")
		proxyType, _ := h.DB.GetSetting("proxy_type")
		proxyHost, _ := h.DB.GetSetting("proxy_host")
		proxyPort, _ := h.DB.GetSetting("proxy_port")
		proxyUsername, _ := h.DB.GetEncryptedSetting("proxy_username")
		proxyPassword, _ := h.DB.GetEncryptedSetting("proxy_password")
		googleTranslateEndpoint, _ := h.DB.GetSetting("google_translate_endpoint")
		showArticlePreviewImages, _ := h.DB.GetSetting("show_article_preview_images")
		json.NewEncoder(w).Encode(map[string]string{
			"update_interval":             interval,
			"refresh_mode":                refreshMode,
			"translation_enabled":         translationEnabled,
			"target_language":             targetLang,
			"translation_provider":        provider,
			"deepl_api_key":               apiKey,
			"baidu_app_id":                baiduAppID,
			"baidu_secret_key":            baiduSecretKey,
			"ai_api_key":                  aiAPIKey,
			"ai_endpoint":                 aiEndpoint,
			"ai_model":                    aiModel,
			"ai_system_prompt":            aiSystemPrompt,
			"auto_cleanup_enabled":        autoCleanup,
			"max_cache_size_mb":           maxCacheSize,
			"max_article_age_days":        maxArticleAge,
			"language":                    language,
			"theme":                       theme,
			"last_article_update":         lastUpdate,
			"show_hidden_articles":        showHidden,
			"startup_on_boot":             startupOnBoot,
			"close_to_tray":               closeToTray,
			"shortcuts":                   shortcuts,
			"rules":                       rules,
			"default_view_mode":           defaultViewMode,
			"media_cache_enabled":         mediaCacheEnabled,
			"media_cache_max_size_mb":     mediaCacheMaxSizeMB,
			"media_cache_max_age_days":    mediaCacheMaxAgeDays,
			"summary_enabled":             summaryEnabled,
			"summary_length":              summaryLength,
			"summary_provider":            summaryProvider,
			"summary_ai_api_key":          summaryAIAPIKey,
			"summary_ai_endpoint":         summaryAIEndpoint,
			"summary_ai_model":            summaryAIModel,
			"summary_ai_system_prompt":    summaryAISystemPrompt,
			"proxy_enabled":               proxyEnabled,
			"proxy_type":                  proxyType,
			"proxy_host":                  proxyHost,
			"proxy_port":                  proxyPort,
			"proxy_username":              proxyUsername,
			"proxy_password":              proxyPassword,
			"google_translate_endpoint":   googleTranslateEndpoint,
			"show_article_preview_images": showArticlePreviewImages,
		})
	case http.MethodPost:
		var req struct {
			UpdateInterval           string `json:"update_interval"`
			RefreshMode              string `json:"refresh_mode"`
			TranslationEnabled       string `json:"translation_enabled"`
			TargetLanguage           string `json:"target_language"`
			TranslationProvider      string `json:"translation_provider"`
			DeepLAPIKey              string `json:"deepl_api_key"`
			BaiduAppID               string `json:"baidu_app_id"`
			BaiduSecretKey           string `json:"baidu_secret_key"`
			AIAPIKey                 string `json:"ai_api_key"`
			AIEndpoint               string `json:"ai_endpoint"`
			AIModel                  string `json:"ai_model"`
			AISystemPrompt           string `json:"ai_system_prompt"`
			AutoCleanupEnabled       string `json:"auto_cleanup_enabled"`
			MaxCacheSizeMB           string `json:"max_cache_size_mb"`
			MaxArticleAgeDays        string `json:"max_article_age_days"`
			Language                 string `json:"language"`
			Theme                    string `json:"theme"`
			ShowHiddenArticles       string `json:"show_hidden_articles"`
			StartupOnBoot            string `json:"startup_on_boot"`
			CloseToTray              string `json:"close_to_tray"`
			Shortcuts                string `json:"shortcuts"`
			Rules                    string `json:"rules"`
			DefaultViewMode          string `json:"default_view_mode"`
			MediaCacheEnabled        string `json:"media_cache_enabled"`
			MediaCacheMaxSizeMB      string `json:"media_cache_max_size_mb"`
			MediaCacheMaxAgeDays     string `json:"media_cache_max_age_days"`
			SummaryEnabled           string `json:"summary_enabled"`
			SummaryLength            string `json:"summary_length"`
			SummaryProvider          string `json:"summary_provider"`
			SummaryAIAPIKey          string `json:"summary_ai_api_key"`
			SummaryAIEndpoint        string `json:"summary_ai_endpoint"`
			SummaryAIModel           string `json:"summary_ai_model"`
			SummaryAISystemPrompt    string `json:"summary_ai_system_prompt"`
			ProxyEnabled             string `json:"proxy_enabled"`
			ProxyType                string `json:"proxy_type"`
			ProxyHost                string `json:"proxy_host"`
			ProxyPort                string `json:"proxy_port"`
			ProxyUsername            string `json:"proxy_username"`
			ProxyPassword            string `json:"proxy_password"`
			GoogleTranslateEndpoint  string `json:"google_translate_endpoint"`
			ShowArticlePreviewImages string `json:"show_article_preview_images"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.UpdateInterval != "" {
			h.DB.SetSetting("update_interval", req.UpdateInterval)
		}
		if req.RefreshMode != "" {
			h.DB.SetSetting("refresh_mode", req.RefreshMode)
		}
		if req.TranslationEnabled != "" {
			h.DB.SetSetting("translation_enabled", req.TranslationEnabled)
		}
		if req.TargetLanguage != "" {
			h.DB.SetSetting("target_language", req.TargetLanguage)
		}
		if req.TranslationProvider != "" {
			h.DB.SetSetting("translation_provider", req.TranslationProvider)
		}
		// Always update API keys as they might be cleared (use encrypted storage for sensitive credentials)
		if err := h.DB.SetEncryptedSetting("deepl_api_key", req.DeepLAPIKey); err != nil {
			log.Printf("Failed to save DeepL API key: %v", err)
			http.Error(w, "Failed to save DeepL API key", http.StatusInternalServerError)
			return
		}
		h.DB.SetSetting("baidu_app_id", req.BaiduAppID)
		if err := h.DB.SetEncryptedSetting("baidu_secret_key", req.BaiduSecretKey); err != nil {
			log.Printf("Failed to save Baidu secret key: %v", err)
			http.Error(w, "Failed to save Baidu secret key", http.StatusInternalServerError)
			return
		}
		if err := h.DB.SetEncryptedSetting("ai_api_key", req.AIAPIKey); err != nil {
			log.Printf("Failed to save AI API key: %v", err)
			http.Error(w, "Failed to save AI API key", http.StatusInternalServerError)
			return
		}
		h.DB.SetSetting("ai_endpoint", req.AIEndpoint)
		h.DB.SetSetting("ai_model", req.AIModel)
		h.DB.SetSetting("ai_system_prompt", req.AISystemPrompt)

		if req.AutoCleanupEnabled != "" {
			h.DB.SetSetting("auto_cleanup_enabled", req.AutoCleanupEnabled)
		}

		if req.MaxCacheSizeMB != "" {
			h.DB.SetSetting("max_cache_size_mb", req.MaxCacheSizeMB)
		}

		if req.MaxArticleAgeDays != "" {
			h.DB.SetSetting("max_article_age_days", req.MaxArticleAgeDays)
		}

		if req.Language != "" {
			h.DB.SetSetting("language", req.Language)
		}

		if req.Theme != "" {
			h.DB.SetSetting("theme", req.Theme)
		}

		if req.ShowHiddenArticles != "" {
			h.DB.SetSetting("show_hidden_articles", req.ShowHiddenArticles)
		}

		if req.CloseToTray != "" {
			h.DB.SetSetting("close_to_tray", req.CloseToTray)
		}

		// Always update shortcuts as it might be cleared or modified
		h.DB.SetSetting("shortcuts", req.Shortcuts)

		// Always update rules as it might be cleared or modified
		h.DB.SetSetting("rules", req.Rules)

		if req.DefaultViewMode != "" {
			h.DB.SetSetting("default_view_mode", req.DefaultViewMode)
		}

		if req.MediaCacheEnabled != "" {
			h.DB.SetSetting("media_cache_enabled", req.MediaCacheEnabled)
		}

		if req.MediaCacheMaxSizeMB != "" {
			h.DB.SetSetting("media_cache_max_size_mb", req.MediaCacheMaxSizeMB)
		}

		if req.MediaCacheMaxAgeDays != "" {
			h.DB.SetSetting("media_cache_max_age_days", req.MediaCacheMaxAgeDays)
		}

		if req.SummaryEnabled != "" {
			h.DB.SetSetting("summary_enabled", req.SummaryEnabled)
		}

		if req.SummaryLength != "" {
			h.DB.SetSetting("summary_length", req.SummaryLength)
		}

		if req.SummaryProvider != "" {
			h.DB.SetSetting("summary_provider", req.SummaryProvider)
		}

		// Always update AI summary keys as they might be cleared (use encrypted storage for API key)
		if err := h.DB.SetEncryptedSetting("summary_ai_api_key", req.SummaryAIAPIKey); err != nil {
			log.Printf("Failed to save summary AI API key: %v", err)
			http.Error(w, "Failed to save summary AI API key", http.StatusInternalServerError)
			return
		}
		h.DB.SetSetting("summary_ai_endpoint", req.SummaryAIEndpoint)
		h.DB.SetSetting("summary_ai_model", req.SummaryAIModel)
		h.DB.SetSetting("summary_ai_system_prompt", req.SummaryAISystemPrompt)

		if req.ProxyEnabled != "" {
			h.DB.SetSetting("proxy_enabled", req.ProxyEnabled)
		}

		// Always update proxy settings as they might be cleared (use encrypted storage for credentials)
		h.DB.SetSetting("proxy_type", req.ProxyType)
		h.DB.SetSetting("proxy_host", req.ProxyHost)
		h.DB.SetSetting("proxy_port", req.ProxyPort)
		if err := h.DB.SetEncryptedSetting("proxy_username", req.ProxyUsername); err != nil {
			log.Printf("Failed to save proxy username: %v", err)
			http.Error(w, "Failed to save proxy credentials", http.StatusInternalServerError)
			return
		}
		if err := h.DB.SetEncryptedSetting("proxy_password", req.ProxyPassword); err != nil {
			log.Printf("Failed to save proxy password: %v", err)
			http.Error(w, "Failed to save proxy credentials", http.StatusInternalServerError)
			return
		}

		// Always update google_translate_endpoint as it might be reset to default
		h.DB.SetSetting("google_translate_endpoint", req.GoogleTranslateEndpoint)

		if req.ShowArticlePreviewImages != "" {
			h.DB.SetSetting("show_article_preview_images", req.ShowArticlePreviewImages)
		}

		if req.StartupOnBoot != "" {
			// Get current value to check if it changed
			currentValue, err := h.DB.GetSetting("startup_on_boot")
			if err != nil {
				log.Printf("Failed to get startup_on_boot setting: %v", err)
				// If we can't read the current value, save the new value but don't apply it
				h.DB.SetSetting("startup_on_boot", req.StartupOnBoot)
			} else if currentValue != req.StartupOnBoot {
				// Only apply if the value changed
				h.DB.SetSetting("startup_on_boot", req.StartupOnBoot)

				// Apply the startup setting
				if req.StartupOnBoot == "true" {
					if err := utils.EnableStartup(); err != nil {
						log.Printf("Failed to enable startup: %v", err)
					}
				} else {
					if err := utils.DisableStartup(); err != nil {
						log.Printf("Failed to disable startup: %v", err)
					}
				}
			}
		}

		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
