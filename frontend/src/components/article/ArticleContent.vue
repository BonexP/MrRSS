<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { ref, watch, onMounted } from 'vue';
import { PhArticle, PhTextAlignLeft, PhListBullets, PhSpinnerGap } from '@phosphor-icons/vue';
import type { Article } from '@/types/models';
import { formatDate } from '@/utils/date';

const { t, locale } = useI18n();

interface SummaryResult {
  summary: string;
  key_points: string[];
  sentence_count: number;
  is_too_short: boolean;
  error?: string;
}

interface Props {
  article: Article;
  articleContent: string;
  isLoadingContent: boolean;
}

const props = defineProps<Props>();

// Summary state
const summaryEnabled = ref(false);
const summaryLength = ref('medium');
const summaryResult = ref<SummaryResult | null>(null);
const isLoadingSummary = ref(false);
const showSummary = ref(true);

// Load summary settings
async function loadSummarySettings() {
  try {
    const res = await fetch('/api/settings');
    const data = await res.json();
    summaryEnabled.value = data.summary_enabled === 'true';
    summaryLength.value = data.summary_length || 'medium';
  } catch (e) {
    console.error('Error loading summary settings:', e);
  }
}

// Generate summary for the current article
async function generateSummary() {
  if (!summaryEnabled.value || !props.article) {
    return;
  }

  isLoadingSummary.value = true;
  summaryResult.value = null;

  try {
    const res = await fetch('/api/articles/summarize', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        article_id: props.article.id,
        length: summaryLength.value,
      }),
    });

    if (res.ok) {
      summaryResult.value = await res.json();
    }
  } catch (e) {
    console.error('Error generating summary:', e);
  } finally {
    isLoadingSummary.value = false;
  }
}

// Watch for article changes and regenerate summary
watch(
  () => props.article?.id,
  () => {
    summaryResult.value = null;
    if (summaryEnabled.value && props.article) {
      generateSummary();
    }
  }
);

// Watch for content loading completion
watch(
  () => props.isLoadingContent,
  (isLoading, wasLoading) => {
    if (wasLoading && !isLoading && summaryEnabled.value && props.article) {
      generateSummary();
    }
  }
);

onMounted(async () => {
  await loadSummarySettings();
  if (summaryEnabled.value && props.article && props.articleContent) {
    generateSummary();
  }
});
</script>

<template>
  <div class="flex-1 overflow-y-auto bg-bg-primary p-3 sm:p-6">
    <div class="max-w-3xl mx-auto bg-bg-primary">
      <h1 class="text-xl sm:text-3xl font-bold mb-3 sm:mb-4 text-text-primary leading-tight">
        {{ article.title }}
      </h1>
      <div
        class="text-xs sm:text-sm text-text-secondary mb-4 sm:mb-6 flex flex-wrap items-center gap-2 sm:gap-4"
      >
        <span>{{ article.feed_title }}</span>
        <span class="hidden sm:inline">•</span>
        <span>{{ formatDate(article.published_at, locale === 'zh-CN' ? 'zh-CN' : 'en-US') }}</span>
      </div>

      <!-- Summary Section -->
      <div
        v-if="summaryEnabled && (isLoadingSummary || summaryResult)"
        class="mb-6 p-4 rounded-lg border border-border bg-bg-secondary"
      >
        <!-- Summary Header -->
        <div
          class="flex items-center justify-between cursor-pointer"
          @click="showSummary = !showSummary"
        >
          <div class="flex items-center gap-2 text-accent font-medium">
            <PhTextAlignLeft :size="20" />
            <span>{{ t('articleSummary') }}</span>
          </div>
          <span class="text-xs text-text-secondary">
            {{ showSummary ? '▲' : '▼' }}
          </span>
        </div>

        <!-- Summary Content -->
        <div v-if="showSummary" class="mt-3">
          <!-- Loading State -->
          <div v-if="isLoadingSummary" class="flex items-center gap-2 text-text-secondary">
            <PhSpinnerGap :size="16" class="animate-spin" />
            <span class="text-sm">{{ t('generatingSummary') }}</span>
          </div>

          <!-- Too Short Warning -->
          <div v-else-if="summaryResult?.is_too_short" class="text-sm text-text-secondary italic">
            {{ t('summaryTooShort') }}
          </div>

          <!-- Summary Display -->
          <div v-else-if="summaryResult?.summary">
            <p class="text-sm text-text-primary leading-relaxed mb-3">
              {{ summaryResult.summary }}
            </p>

            <!-- Key Points -->
            <div
              v-if="summaryResult.key_points && summaryResult.key_points.length > 1"
              class="mt-3"
            >
              <div class="flex items-center gap-2 text-text-secondary text-xs font-medium mb-2">
                <PhListBullets :size="14" />
                <span>{{ t('keyPoints') }}</span>
              </div>
              <ul class="space-y-1">
                <li
                  v-for="(point, index) in summaryResult.key_points.slice(0, 3)"
                  :key="index"
                  class="text-xs text-text-secondary pl-3 border-l-2 border-accent/30"
                >
                  {{ point }}
                </li>
              </ul>
            </div>
          </div>

          <!-- No Summary Available -->
          <div v-else class="text-sm text-text-secondary italic">
            {{ t('noSummaryAvailable') }}
          </div>
        </div>
      </div>

      <!-- Loading state with proper background -->
      <div
        v-if="isLoadingContent"
        class="flex flex-col items-center justify-center py-12 sm:py-16 bg-bg-primary"
      >
        <div class="relative mb-4 sm:mb-6">
          <!-- Outer pulsing ring -->
          <div
            class="absolute inset-0 rounded-full border-2 sm:border-4 border-accent animate-ping opacity-20"
          ></div>
          <!-- Middle spinning ring -->
          <div
            class="absolute inset-0 rounded-full border-2 sm:border-4 border-t-accent border-r-transparent border-b-transparent border-l-transparent animate-spin"
          ></div>
          <!-- Inner icon -->
          <div class="relative bg-bg-secondary rounded-full p-4 sm:p-6">
            <PhArticle :size="48" class="text-accent sm:w-16 sm:h-16" />
          </div>
        </div>
        <p class="text-base sm:text-lg font-medium text-text-primary mb-1 sm:mb-2">
          {{ t('loadingContent') }}
        </p>
        <p class="text-xs sm:text-sm text-text-secondary px-4 text-center">
          {{ t('fetchingArticleContent') }}
        </p>
      </div>

      <!-- Content display -->
      <div
        v-else-if="articleContent"
        class="prose prose-sm sm:prose-lg max-w-none text-text-primary"
        v-html="articleContent"
      ></div>

      <!-- No content available -->
      <div v-else class="text-center text-text-secondary py-6 sm:py-8">
        <PhArticle :size="48" class="mb-2 sm:mb-3 opacity-50 mx-auto sm:w-16 sm:h-16" />
        <p class="text-sm sm:text-base">{{ t('noContent') }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Prose styling for article content */
.prose {
  color: var(--text-primary);
}
.prose :deep(h1),
.prose :deep(h2),
.prose :deep(h3),
.prose :deep(h4),
.prose :deep(h5),
.prose :deep(h6) {
  color: var(--text-primary);
  font-weight: 600;
  margin-top: 1.5em;
  margin-bottom: 0.75em;
}
.prose :deep(p) {
  margin-bottom: 1em;
  line-height: 1.7;
}
.prose :deep(a) {
  color: var(--accent-color);
  text-decoration: underline;
  cursor: pointer;
}
.prose :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.5rem;
  margin: 1.5em 0;
  cursor: pointer;
  transition: opacity 0.2s;
}
.prose :deep(img:hover) {
  opacity: 0.9;
}
.prose :deep(pre) {
  background-color: var(--bg-secondary);
  padding: 1em;
  border-radius: 0.5rem;
  overflow-x: auto;
  margin: 1em 0;
}
.prose :deep(code) {
  background-color: var(--bg-secondary);
  padding: 0.2em 0.4em;
  border-radius: 0.25rem;
  font-size: 0.9em;
}
.prose :deep(blockquote) {
  border-left: 4px solid var(--accent-color);
  padding-left: 1em;
  margin: 1em 0;
  font-style: italic;
  color: var(--text-secondary);
}
.prose :deep(ul),
.prose :deep(ol) {
  margin: 1em 0;
  padding-left: 2em;
}
.prose :deep(li) {
  margin-bottom: 0.5em;
}
</style>
