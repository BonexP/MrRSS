<script setup lang="ts">
import { PhSpinnerGap, PhArticle } from '@phosphor-icons/vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

interface Props {
  articleContent: string;
  isTranslatingContent: boolean;
  hasMediaContent?: boolean; // Whether article has audio/video content
}

withDefaults(defineProps<Props>(), {
  hasMediaContent: false,
});
</script>

<template>
  <!-- Content display with inline translations -->
  <div v-if="articleContent">
    <!-- eslint-disable-next-line vue/no-v-html -->
    <div
      class="prose prose-sm sm:prose-lg max-w-none text-text-primary prose-content"
      v-html="articleContent"
    ></div>
    <!-- Translation loading indicator -->
    <div v-if="isTranslatingContent" class="flex items-center gap-2 mt-4 text-text-secondary">
      <PhSpinnerGap :size="16" class="animate-spin" />
      <span class="text-sm">{{ t('translatingContent') }}</span>
    </div>
  </div>

  <!-- No content available (only show if there's no media content either) -->
  <div v-else-if="!hasMediaContent" class="text-center text-text-secondary py-6 sm:py-8">
    <PhArticle :size="48" class="mb-2 sm:mb-3 opacity-50 mx-auto sm:w-16 sm:h-16" />
    <p class="text-sm sm:text-base">{{ t('noContentAvailable') }}</p>
  </div>
</template>
