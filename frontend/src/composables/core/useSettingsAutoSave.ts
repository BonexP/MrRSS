/**
 * Composable for auto-saving settings with debouncing
 */
import { ref, watch, onMounted, onUnmounted, type Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAppStore } from '@/stores/app';
import type { SettingsData } from '@/types/settings';
import { settingsDefaults } from '@/config/defaults';
import { useSettingsValidation } from './useSettingsValidation';

export function useSettingsAutoSave(settings: Ref<SettingsData>) {
  const { locale } = useI18n();
  const store = useAppStore();

  let saveTimeout: ReturnType<typeof setTimeout> | null = null;
  let isInitialLoad = true;

  // Use validation composable
  const { isValid } = useSettingsValidation(settings);

  // Track previous translation settings
  const prevTranslationSettings: Ref<{
    enabled: boolean;
    targetLang: string;
    provider: string;
  }> = ref({
    enabled: settingsDefaults.translation_enabled,
    targetLang: settingsDefaults.target_language,
    provider: settingsDefaults.translation_provider,
  });

  /**
   * Initialize translation tracking
   */
  onMounted(() => {
    setTimeout(() => {
      prevTranslationSettings.value = {
        enabled: settings.value.translation_enabled,
        targetLang: settings.value.target_language,
        provider: settings.value.translation_provider,
      };
      isInitialLoad = false;
    }, 100);
  });

  /**
   * Save settings to backend
   */
  async function autoSave() {
    try {
      // Skip translation clearing on initial load
      if (isInitialLoad) {
        return;
      }

      // Validate settings before saving
      if (!isValid.value) {
        // Skip save silently - validation errors are shown in UI
        return;
      }

      // Check if translation settings changed
      const translationChanged =
        prevTranslationSettings.value.enabled !== settings.value.translation_enabled ||
        prevTranslationSettings.value.provider !== settings.value.translation_provider ||
        (settings.value.translation_enabled &&
          prevTranslationSettings.value.targetLang !== settings.value.target_language);

      await fetch('/api/settings', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          update_interval: (
            settings.value.update_interval ?? settingsDefaults.update_interval
          ).toString(),
          refresh_mode: settings.value.refresh_mode ?? settingsDefaults.refresh_mode,
          translation_enabled: (
            settings.value.translation_enabled ?? settingsDefaults.translation_enabled
          ).toString(),
          target_language: settings.value.target_language ?? settingsDefaults.target_language,
          translation_provider:
            settings.value.translation_provider ?? settingsDefaults.translation_provider,
          deepl_api_key: settings.value.deepl_api_key ?? settingsDefaults.deepl_api_key,
          baidu_app_id: settings.value.baidu_app_id ?? settingsDefaults.baidu_app_id,
          baidu_secret_key: settings.value.baidu_secret_key ?? settingsDefaults.baidu_secret_key,
          ai_api_key: settings.value.ai_api_key ?? settingsDefaults.ai_api_key,
          ai_endpoint: settings.value.ai_endpoint ?? settingsDefaults.ai_endpoint,
          ai_model: settings.value.ai_model ?? settingsDefaults.ai_model,
          ai_system_prompt: settings.value.ai_system_prompt ?? settingsDefaults.ai_system_prompt,
          auto_cleanup_enabled: (
            settings.value.auto_cleanup_enabled ?? settingsDefaults.auto_cleanup_enabled
          ).toString(),
          max_cache_size_mb: (
            settings.value.max_cache_size_mb ?? settingsDefaults.max_cache_size_mb
          ).toString(),
          max_article_age_days: (
            settings.value.max_article_age_days ?? settingsDefaults.max_article_age_days
          ).toString(),
          language: settings.value.language ?? settingsDefaults.language,
          theme: settings.value.theme ?? settingsDefaults.theme,
          show_hidden_articles: (
            settings.value.show_hidden_articles ?? settingsDefaults.show_hidden_articles
          ).toString(),
          default_view_mode: settings.value.default_view_mode ?? settingsDefaults.default_view_mode,
          media_cache_enabled: (
            settings.value.media_cache_enabled ?? settingsDefaults.media_cache_enabled
          ).toString(),
          media_cache_max_size_mb: (
            settings.value.media_cache_max_size_mb ?? settingsDefaults.media_cache_max_size_mb
          ).toString(),
          media_cache_max_age_days: (
            settings.value.media_cache_max_age_days ?? settingsDefaults.media_cache_max_age_days
          ).toString(),
          startup_on_boot: (
            settings.value.startup_on_boot ?? settingsDefaults.startup_on_boot
          ).toString(),
          shortcuts: settings.value.shortcuts ?? settingsDefaults.shortcuts,
          summary_enabled: (
            settings.value.summary_enabled ?? settingsDefaults.summary_enabled
          ).toString(),
          summary_length: settings.value.summary_length ?? settingsDefaults.summary_length,
          summary_provider: settings.value.summary_provider ?? settingsDefaults.summary_provider,
          summary_ai_api_key:
            settings.value.summary_ai_api_key ?? settingsDefaults.summary_ai_api_key,
          summary_ai_endpoint:
            settings.value.summary_ai_endpoint ?? settingsDefaults.summary_ai_endpoint,
          summary_ai_model: settings.value.summary_ai_model ?? settingsDefaults.summary_ai_model,
          summary_ai_system_prompt:
            settings.value.summary_ai_system_prompt ?? settingsDefaults.summary_ai_system_prompt,
          proxy_enabled: (
            settings.value.proxy_enabled ?? settingsDefaults.proxy_enabled
          ).toString(),
          proxy_type: settings.value.proxy_type ?? settingsDefaults.proxy_type,
          proxy_host: settings.value.proxy_host ?? settingsDefaults.proxy_host,
          proxy_port: settings.value.proxy_port ?? settingsDefaults.proxy_port,
          proxy_username: settings.value.proxy_username ?? settingsDefaults.proxy_username,
          proxy_password: settings.value.proxy_password ?? settingsDefaults.proxy_password,
          google_translate_endpoint:
            settings.value.google_translate_endpoint ?? settingsDefaults.google_translate_endpoint,
          show_article_preview_images: (
            settings.value.show_article_preview_images ??
            settingsDefaults.show_article_preview_images
          ).toString(),
        }),
      });

      // Apply settings immediately
      locale.value = settings.value.language;
      store.setTheme(settings.value.theme as 'light' | 'dark' | 'auto');
      store.startAutoRefresh(settings.value.update_interval);

      // Notify components about default view mode change
      window.dispatchEvent(
        new CustomEvent('default-view-mode-changed', {
          detail: {
            mode: settings.value.default_view_mode,
          },
        })
      );

      // Clear and re-translate if translation settings changed
      if (translationChanged) {
        await fetch('/api/articles/clear-translations', { method: 'POST' });
        // Update tracking
        prevTranslationSettings.value = {
          enabled: settings.value.translation_enabled,
          targetLang: settings.value.target_language,
          provider: settings.value.translation_provider,
        };
        // Notify ArticleList about translation settings change
        window.dispatchEvent(
          new CustomEvent('translation-settings-changed', {
            detail: {
              enabled: settings.value.translation_enabled,
              targetLang: settings.value.target_language,
            },
          })
        );
        // Refresh articles to show without translations, then re-translate if enabled
        store.fetchArticles();
      }

      // Refresh articles if show_hidden_articles changed
      if (settings.value.show_hidden_articles !== undefined) {
        store.fetchArticles();
      }

      // Notify about show_article_preview_images change
      window.dispatchEvent(
        new CustomEvent('show-preview-images-changed', {
          detail: {
            value: settings.value.show_article_preview_images,
          },
        })
      );
    } catch (e) {
      console.error('Error auto-saving settings:', e);
    }
  }

  /**
   * Debounced auto-save function
   */
  function debouncedAutoSave() {
    if (saveTimeout) {
      clearTimeout(saveTimeout);
    }
    saveTimeout = setTimeout(autoSave, 500); // Wait 500ms after last change
  }

  // Watch the entire settings object for changes
  watch(() => settings.value, debouncedAutoSave, { deep: true });

  // Clean up timeout on unmount to prevent memory leaks
  onUnmounted(() => {
    if (saveTimeout) {
      clearTimeout(saveTimeout);
      saveTimeout = null;
    }
  });

  return {
    autoSave,
    debouncedAutoSave,
  };
}
