<script setup lang="ts">
import { PhWarningCircle } from '@phosphor-icons/vue';
import type { Feed } from '@/types/models';

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
    @click="emit('click')"
    @contextmenu="(e) => emit('contextmenu', e)"
    :class="['feed-item', isActive ? 'active' : '']"
  >
    <div class="w-4 h-4 flex items-center justify-center shrink-0">
      <img
        :src="feed.image_url || getFavicon(feed.url)"
        class="w-full h-full object-contain"
        @error="($event.target as HTMLElement).style.display = 'none'"
      />
    </div>
    <span class="truncate flex-1">{{ feed.title }}</span>
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
.feed-item {
  @apply px-2 sm:px-3 py-1.5 sm:py-2 cursor-pointer rounded-md text-xs sm:text-sm text-text-primary flex items-center gap-1.5 sm:gap-2.5 hover:bg-bg-tertiary transition-colors;
}
.feed-item.active {
  @apply bg-bg-tertiary text-accent font-medium;
}
.unread-badge {
  @apply text-[9px] sm:text-[10px] font-semibold rounded-full min-w-[14px] sm:min-w-[16px] h-[14px] sm:h-[16px] px-0.5 sm:px-1 flex items-center justify-center;
  background-color: rgba(200, 200, 200, 0.3);
  color: #666666;
}
.dark-mode .unread-badge {
  background-color: rgba(160, 160, 160, 0.25);
  color: #cccccc;
}
</style>
