/**
 * Composable for auto-saving settings with debouncing
 */
import { ref, watch, onMounted, onUnmounted, type Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAppStore } from '@/stores/app';
import type { SettingsData } from '@/types/settings';

export function useSettingsAutoSave(settings: Ref<SettingsData>) {
  const { locale } = useI18n();
  const store = useAppStore();

  let saveTimeout: ReturnType<typeof setTimeout> | null = null;
  let isInitialLoad = true;

  // Track previous translation settings
  const prevTranslationSettings: Ref<{ enabled: boolean; targetLang: string }> = ref({
    enabled: false,
    targetLang: 'zh',
  });

  /**
   * Initialize translation tracking
   */
  onMounted(() => {
    setTimeout(() => {
      prevTranslationSettings.value = {
        enabled: settings.value.translation_enabled,
        targetLang: settings.value.target_language,
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

      // Check if translation settings changed
      const translationChanged =
        prevTranslationSettings.value.enabled !== settings.value.translation_enabled ||
        (settings.value.translation_enabled &&
          prevTranslationSettings.value.targetLang !== settings.value.target_language);

      await fetch('/api/settings', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          update_interval: settings.value.update_interval.toString(),
          translation_enabled: settings.value.translation_enabled.toString(),
          target_language: settings.value.target_language,
          translation_provider: settings.value.translation_provider,
          deepl_api_key: settings.value.deepl_api_key,
          auto_cleanup_enabled: settings.value.auto_cleanup_enabled.toString(),
          max_cache_size_mb: settings.value.max_cache_size_mb.toString(),
          max_article_age_days: settings.value.max_article_age_days.toString(),
          language: settings.value.language,
          theme: settings.value.theme,
          show_hidden_articles: settings.value.show_hidden_articles.toString(),
          default_view_mode: settings.value.default_view_mode,
          startup_on_boot: settings.value.startup_on_boot.toString(),
          shortcuts: settings.value.shortcuts || '',
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
