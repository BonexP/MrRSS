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
        <PhFolder v-else :size="20" :weight="'fill'" />
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
@reference "../../style.css";

.category-header {
  @apply px-2 sm:px-3 py-1.5 sm:py-2 cursor-pointer font-semibold text-xs sm:text-sm text-text-secondary flex items-center justify-between rounded-md hover:bg-bg-tertiary hover:text-text-primary transition-colors;
  @apply sticky z-10 bg-bg-secondary;
  top: -0.375rem; /* matches container's p-1.5 */
  margin-left: -0.375rem;
  margin-right: -0.375rem;
  padding-left: calc(0.5rem + 0.375rem);
  padding-right: calc(0.75rem + 0.375rem);
}
@media (min-width: 640px) {
  .category-header {
    top: -0.5rem; /* matches container's sm:p-2 */
    margin-left: -0.5rem;
    margin-right: -0.5rem;
    padding-left: calc(0.75rem + 0.5rem);
    padding-right: calc(0.75rem + 0.5rem);
  }
}
.category-header.active {
  @apply bg-bg-tertiary text-accent;
}
.unread-badge {
  @apply text-[9px] sm:text-[10px] font-semibold rounded-full min-w-[14px] sm:min-w-[16px] h-[14px] sm:h-[16px] px-0.5 sm:px-1 flex items-center justify-center;
  background-color: rgba(120, 120, 120, 0.25);
  color: #444444;
}
</style>

<style>
.dark-mode .unread-badge {
  /* This style will be applied to child components, so it can not use scoped */
  background-color: rgba(100, 100, 100, 0.6) !important;
  color: #f0f0f0 !important;
}
</style>
