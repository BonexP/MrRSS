/**
 * Composable for settings management
 */
import { ref, type Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import type { SettingsData } from '@/types/settings';
import type { ThemePreference } from '@/stores/app';
import { generateInitialSettings, parseSettingsData } from './useSettings.generated';

export function useSettings() {
  const { locale } = useI18n();

  // Use generated helper for initial settings (alphabetically sorted)
  const settings: Ref<SettingsData> = ref(generateInitialSettings());

  // Override language with current locale if available
  if (locale.value) {
    settings.value.language = locale.value;
  }

  /**
   * Fetch settings from backend
   */
  async function fetchSettings(): Promise<SettingsData> {
    try {
      const res = await fetch('/api/settings');
      const data = await res.json();

      // Use generated helper to parse settings (alphabetically sorted)
      settings.value = parseSettingsData(data);

      // Override language with current locale if not set
      if (!settings.value.language && locale.value) {
        settings.value.language = locale.value;
      }

      return settings.value;
    } catch (e) {
      console.error('Error fetching settings:', e);
      throw e;
    }
  }

  /**
   * Apply fetched settings to the app
   */

  function applySettings(data: SettingsData, setTheme: (preference: ThemePreference) => void) {
    // Apply the saved language
    if (data.language) {
      locale.value = data.language;
    }

    // Apply the saved theme
    if (data.theme) {
      setTheme(data.theme as ThemePreference);
    }

    // Initialize shortcuts in store
    if (data.shortcuts) {
      try {
        const parsed = JSON.parse(data.shortcuts);
        window.dispatchEvent(
          new CustomEvent('shortcuts-changed', {
            detail: { shortcuts: parsed },
          })
        );
      } catch (e) {
        console.error('Error parsing shortcuts:', e);
      }
    }
  }

  return {
    settings,
    fetchSettings,
    applySettings,
  };
}
