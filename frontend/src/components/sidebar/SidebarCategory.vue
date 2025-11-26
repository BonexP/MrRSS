<script setup lang="ts">
import { PhFolder, PhFolderDashed, PhCaretDown } from '@phosphor-icons/vue';
import type { Feed } from '@/types/models';
import SidebarFeed from './SidebarFeed.vue';

interface Props {
  name: string;
  feeds: Feed[];
  isOpen: boolean;
  isActive: boolean;
  isUncategorized?: boolean;
  unreadCount: number;
  currentFeedId: number | null;
  feedUnreadCounts: Record<number, number>;
}

defineProps<Props>();

const emit = defineEmits<{
  toggle: [];
  selectCategory: [];
  selectFeed: [feedId: number];
  categoryContextMenu: [event: MouseEvent];
  feedContextMenu: [event: MouseEvent, feed: Feed];
}>();
</script>

<template>
  <div class="mb-1">
    <div
      :class="['category-header', isActive ? 'active' : '']"
      @contextmenu="(e) => emit('categoryContextMenu', e)"
    >
      <span class="flex-1 flex items-center gap-2" @click="emit('selectCategory')">
        <PhFolderDashed v-if="isUncategorized" :size="20" />
        <PhFolder v-else :size="20" />
        {{ name }}
      </span>
      <span v-if="unreadCount > 0" class="unread-badge mr-1">{{ unreadCount }}</span>
      <PhCaretDown
        :size="20"
        class="p-1 cursor-pointer transition-transform"
        :class="{ 'rotate-180': isOpen }"
        @click.stop="emit('toggle')"
      />
    </div>
    <div v-show="isOpen" class="pl-2">
      <SidebarFeed
        v-for="feed in feeds"
        :key="feed.id"
        :feed="feed"
        :is-active="currentFeedId === feed.id"
        :unread-count="feedUnreadCounts[feed.id] || 0"
        @click="emit('selectFeed', feed.id)"
        @contextmenu="(e) => emit('feedContextMenu', e, feed)"
      />
    </div>
  </div>
</template>

<style scoped>
.category-header {
  @apply px-2 sm:px-3 py-1.5 sm:py-2 cursor-pointer font-semibold text-xs sm:text-sm text-text-secondary flex items-center justify-between rounded-md hover:bg-bg-tertiary hover:text-text-primary transition-colors;
}
.category-header.active {
  @apply bg-bg-tertiary text-accent;
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
