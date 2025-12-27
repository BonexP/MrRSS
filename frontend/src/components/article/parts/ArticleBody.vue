<script setup lang="ts">
/* eslint-disable vue/no-v-html */
/* eslint-disable no-undef */
import { ref, onMounted, watch, computed, nextTick } from 'vue';
import { PhSpinnerGap, PhArticle, PhArrowClockwise } from '@phosphor-icons/vue';
import { useI18n } from 'vue-i18n';
import { useSettings } from '@/composables/core/useSettings';

const { t } = useI18n();
const { settings, fetchSettings } = useSettings();

interface Props {
  articleContent: string;
  isTranslatingContent: boolean;
  hasMediaContent?: boolean; // Whether article has audio/video content
  isLoadingContent?: boolean; // Whether content is currently loading
}

const props = withDefaults(defineProps<Props>(), {
  hasMediaContent: false,
  isLoadingContent: false,
});

// Emits
const emit = defineEmits<{
  retryLoad: [];
}>();

const customCSS = ref('');
let styleElement: HTMLStyleElement | null = null;

const hasCustomCSS = computed(() => !!settings.value.custom_css_file);

const injectCustomCSS = (css: string) => {
  // Remove existing style element if any
  if (styleElement && styleElement.parentNode) {
    styleElement.parentNode.removeChild(styleElement);
  }

  if (!css) return;

  // Create new style element
  styleElement = document.createElement('style');
  styleElement.textContent = css;
  styleElement.setAttribute('data-custom-css', 'article');

  // Inject to document head
  document.head.appendChild(styleElement);

  console.log('Custom CSS injected to head, length:', css.length);
};

const removeCustomCSS = () => {
  if (styleElement && styleElement.parentNode) {
    styleElement.parentNode.removeChild(styleElement);
    styleElement = null;
    console.log('Custom CSS removed');
  }
};

const loadCustomCSS = async () => {
  console.log('Loading custom CSS, custom_css_file:', settings.value.custom_css_file);

  // First, refresh settings from backend to get latest custom_css_file value
  try {
    await fetchSettings();
    console.log('Settings refreshed, custom_css_file:', settings.value.custom_css_file);
  } catch (error) {
    console.error('Failed to refresh settings:', error);
  }

  if (!settings.value.custom_css_file) {
    customCSS.value = '';
    removeCustomCSS();
    return;
  }

  try {
    const response = await fetch('/api/custom-css');
    if (response.ok) {
      customCSS.value = await response.text();
      console.log('Custom CSS loaded successfully, length:', customCSS.value.length);
      console.log('CSS preview:', customCSS.value.substring(0, 200));

      // Inject CSS to document head
      await nextTick();
      injectCustomCSS(customCSS.value);
    } else {
      console.warn('Failed to load custom CSS:', response.statusText);
      customCSS.value = '';
    }
  } catch (error) {
    console.error('Error loading custom CSS:', error);
    customCSS.value = '';
  }
};

onMounted(() => {
  console.log('ArticleBody mounted');
  loadCustomCSS();

  // Listen for custom CSS change events
  window.addEventListener('custom-css-changed', loadCustomCSS);
});

// Watch for changes in custom_css_file setting
watch(
  () => settings.value.custom_css_file,
  () => {
    console.log('custom_css_file changed:', settings.value.custom_css_file);
    loadCustomCSS();
  }
);

// Clean up on unmount
import { onUnmounted } from 'vue';
onUnmounted(() => {
  removeCustomCSS();
  window.removeEventListener('custom-css-changed', loadCustomCSS);
});
</script>

<template>
  <!-- Content display with inline translations -->
  <div v-if="articleContent">
    <div
      class="prose prose-sm sm:prose-lg max-w-none text-text-primary prose-content"
      :class="{ 'custom-css-active': hasCustomCSS }"
      v-html="articleContent"
    ></div>
    <!-- Translation loading indicator -->
    <div v-if="isTranslatingContent" class="flex items-center gap-2 mt-4 text-text-secondary">
      <PhSpinnerGap :size="16" class="animate-spin" />
      <span class="text-sm">{{ t('translatingContent') }}</span>
    </div>
  </div>

  <!-- No content available with retry option -->
  <div v-else-if="!hasMediaContent" class="text-center text-text-secondary py-6 sm:py-8">
    <PhArticle :size="48" class="mb-2 sm:mb-3 opacity-50 mx-auto sm:w-16 sm:h-16" />
    <p class="text-sm sm:text-base mb-4">{{ t('noContentAvailable') }}</p>
    <button
      v-if="!props.isLoadingContent"
      class="btn-secondary-compact flex items-center gap-1.5 mx-auto"
      @click="emit('retryLoad')"
    >
      <PhArrowClockwise :size="12" />
      <span class="text-xs">{{ t('retrySummary') }}</span>
    </button>
  </div>
</template>
