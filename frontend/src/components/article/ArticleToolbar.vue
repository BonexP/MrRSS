<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import {
  PhArrowLeft,
  PhGlobe,
  PhArticle,
  PhEnvelopeOpen,
  PhEnvelope,
  PhStar,
  PhArrowSquareOut,
} from '@phosphor-icons/vue';
import type { Article } from '@/types/models';

const { t } = useI18n();

interface Props {
  article: Article;
  showContent: boolean;
}

defineProps<Props>();
</script>

<template>
  <div
    class="h-[44px] sm:h-[50px] px-3 sm:px-5 border-b border-border flex justify-between items-center bg-bg-primary shrink-0"
  >
    <button
      @click="$emit('close')"
      class="md:hidden flex items-center gap-1.5 sm:gap-2 text-text-secondary hover:text-text-primary text-sm sm:text-base"
    >
      <PhArrowLeft :size="18" class="sm:w-5 sm:h-5" />
      <span class="hidden xs:inline">{{ t('back') }}</span>
    </button>
    <div class="flex gap-1 sm:gap-2 ml-auto">
      <button
        @click="$emit('toggleContentView')"
        class="action-btn"
        :title="showContent ? t('viewOriginal') : t('viewContent')"
      >
        <PhGlobe v-if="showContent" :size="18" class="sm:w-5 sm:h-5" />
        <PhArticle v-else :size="18" class="sm:w-5 sm:h-5" />
      </button>
      <button
        @click="$emit('toggleRead')"
        class="action-btn"
        :title="article.is_read ? t('markAsUnread') : t('markAsRead')"
      >
        <PhEnvelopeOpen v-if="article.is_read" :size="18" class="sm:w-5 sm:h-5" />
        <PhEnvelope v-else :size="18" class="sm:w-5 sm:h-5" />
      </button>
      <button
        @click="$emit('toggleFavorite')"
        :class="[
          'action-btn',
          article.is_favorite ? 'text-yellow-500 hover:text-yellow-600' : 'hover:text-yellow-500',
        ]"
        :title="article.is_favorite ? t('removeFromFavorite') : t('addToFavorite')"
      >
        <PhStar
          :size="18"
          class="sm:w-5 sm:h-5"
          :weight="article.is_favorite ? 'fill' : 'regular'"
        />
      </button>
      <button @click="$emit('openOriginal')" class="action-btn" :title="t('openInBrowser')">
        <PhArrowSquareOut :size="18" class="sm:w-5 sm:h-5" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.action-btn {
  @apply text-lg sm:text-xl cursor-pointer text-text-secondary p-1 sm:p-1.5 rounded-md transition-colors hover:bg-bg-tertiary hover:text-text-primary;
}
</style>
