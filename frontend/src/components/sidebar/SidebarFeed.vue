<script setup lang="ts">
import { PhWarningCircle, PhEyeSlash } from '@phosphor-icons/vue';
import type { Feed } from '@/types/models';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

interface Props {
  feed: Feed;
  isActive: boolean;
  unreadCount: number;
}

defineProps<Props>();

const emit = defineEmits<{
  click: [];
  contextmenu: [event: MouseEvent];
}>();

function getFavicon(url: string): string {
  try {
    return `https://www.google.com/s2/favicons?domain=${new URL(url).hostname}`;
  } catch {
    return '';
  }
}
</script>

<template>
  <div
    :class="['feed-item', isActive ? 'active' : '']"
    @click="emit('click')"
    @contextmenu="(e) => emit('contextmenu', e)"
  >
    <div class="w-4 h-4 flex items-center justify-center shrink-0">
      <img
        :src="feed.image_url || getFavicon(feed.url)"
        class="w-full h-full object-contain"
        @error="($event.target as HTMLElement).style.display = 'none'"
      />
    </div>
    <span class="truncate flex-1">{{ feed.title }}</span>
    <PhEyeSlash
      v-if="feed.hide_from_timeline"
      :size="16"
      class="text-text-secondary shrink-0"
      :title="t('hideFromTimeline')"
    />
    <PhWarningCircle
      v-if="feed.last_error"
      :size="16"
      class="text-yellow-500 shrink-0"
      :title="feed.last_error"
    />
    <span v-if="unreadCount > 0" class="unread-badge">{{ unreadCount }}</span>
  </div>
</template>

<style scoped>
@reference "../../style.css";

.feed-item {
  @apply px-2 sm:px-3 py-1.5 sm:py-2 cursor-pointer rounded-md text-xs sm:text-sm text-text-primary flex items-center gap-1.5 sm:gap-2.5 hover:bg-bg-tertiary transition-colors;
}
.feed-item.active {
  @apply bg-bg-tertiary text-accent font-medium;
}
.unread-badge {
  @apply text-[9px] sm:text-[10px] font-semibold rounded-full min-w-[14px] sm:min-w-[16px] h-[14px] sm:h-[16px] px-0.5 sm:px-1 flex items-center justify-center;
  background-color: rgba(120, 120, 120, 0.25);
  color: #444444;
}
</style>

<style>
.dark-mode .unread-badge {
  background-color: rgba(100, 100, 100, 0.6) !important;
  color: #f0f0f0 !important;
}
</style>
