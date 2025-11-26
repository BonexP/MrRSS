<script setup lang="ts">
import { PhListDashes, PhCircle, PhStar } from '@phosphor-icons/vue';
import { Component } from 'vue';

interface Props {
  label: string;
  isActive: boolean;
  icon: 'all' | 'unread' | 'favorites';
  unreadCount?: number;
}

defineProps<Props>();

const emit = defineEmits<{
  click: [];
}>();

const iconMap: Record<string, Component> = {
  all: PhListDashes,
  unread: PhCircle,
  favorites: PhStar,
};
</script>

<template>
  <button @click="emit('click')" :class="['nav-item', isActive ? 'active' : '']">
    <component :is="iconMap[icon]" :size="20" />
    <span class="flex-1 text-left">{{ label }}</span>
    <span v-if="unreadCount && unreadCount > 0" class="unread-badge">{{ unreadCount }}</span>
  </button>
</template>

<style scoped>
.nav-item {
  @apply flex items-center gap-2 sm:gap-3 w-full px-2 sm:px-3 py-2 sm:py-2.5 text-text-secondary rounded-lg font-medium transition-colors hover:bg-bg-tertiary hover:text-text-primary text-left text-sm sm:text-base;
}
.nav-item.active {
  @apply bg-bg-tertiary text-accent font-semibold;
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
