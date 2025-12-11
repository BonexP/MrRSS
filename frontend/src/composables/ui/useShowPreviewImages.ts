/**
 * Composable for managing the show preview images setting
 * This is a shared singleton that caches the setting value to avoid duplicate API calls
 */
import { ref, readonly } from 'vue';

// Shared state across all component instances
const showPreviewImages = ref(true);
let isInitialized = false;

export function useShowPreviewImages() {
  /**
   * Initialize the setting value from the API
   * Only called once for the entire application
   */
  async function initialize() {
    if (isInitialized) return;

    try {
      const res = await fetch('/api/settings');
      const data = await res.json();
      showPreviewImages.value = data.show_article_preview_images === 'true';
      isInitialized = true;
    } catch (e) {
      console.error('Error loading show preview images setting:', e);
      // Default to true on error
      showPreviewImages.value = true;
      isInitialized = true;
    }
  }

  /**
   * Update the setting value
   * This is called when the setting changes in the settings modal
   */
  function updateValue(value: boolean) {
    showPreviewImages.value = value;
  }

  return {
    showPreviewImages: readonly(showPreviewImages),
    initialize,
    updateValue,
  };
}
