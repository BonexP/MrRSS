<script setup lang="ts">
import { ref, onMounted, onUnmounted, type Ref, type Component } from 'vue';
import * as PhosphorIcons from '@phosphor-icons/vue';

export interface ContextMenuItem {
  label?: string;
  action?: string;
  icon?: string;
  iconWeight?: 'regular' | 'bold' | 'light' | 'fill' | 'duotone' | 'thin';
  iconColor?: string;
  disabled?: boolean;
  danger?: boolean;
  separator?: boolean;
}

interface Props {
  items: ContextMenuItem[];
  x: number;
  y: number;
}

defineProps<Props>();

const emit = defineEmits<{
  close: [];
  action: [action: string];
}>();

const menuRef: Ref<HTMLDivElement | null> = ref(null);

// Map old icon names to new component names
const iconMap: Record<string, string> = {
  'ph-check-circle': 'PhCheckCircle',
  'ph-globe': 'PhGlobe',
  'ph-pencil': 'PhPencil',
  'ph-trash': 'PhTrash',
  'ph-envelope': 'PhEnvelope',
  'ph-envelope-open': 'PhEnvelopeOpen',
  'ph-star': 'PhStar',
  'ph-article': 'PhArticle',
  'ph-eye': 'PhEye',
  'ph-eye-slash': 'PhEyeSlash',
  'ph-arrow-square-out': 'PhArrowSquareOut',
  PhMagnifyingGlass: 'PhMagnifyingGlass',
};

// Get icon component from icon string
function getIconComponent(iconName?: string): Component | null {
  if (!iconName) return null;
  const componentName = iconMap[iconName] || iconName;
  return (PhosphorIcons as Record<string, Component>)[componentName] || null;
}

function handleClickOutside(event: MouseEvent) {
  if (menuRef.value && event.target instanceof Node && !menuRef.value.contains(event.target)) {
    emit('close');
  }
}

onMounted(() => {
  // Use setTimeout to avoid catching the event that opened the menu
  setTimeout(() => {
    document.addEventListener('click', handleClickOutside);
    document.addEventListener('contextmenu', handleClickOutside);
  }, 0);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
  document.removeEventListener('contextmenu', handleClickOutside);
});

function handleAction(item: ContextMenuItem) {
  if (item.disabled) return;
  if (item.action) {
    emit('action', item.action);
  }
  emit('close');
}
</script>

<template>
  <div
    ref="menuRef"
    class="fixed z-50 bg-bg-primary border border-border rounded-lg shadow-xl py-1 min-w-[180px] animate-fade-in"
    :style="{ top: `${y}px`, left: `${x}px` }"
  >
    <template v-for="(item, index) in items" :key="index">
      <div v-if="item.separator" class="h-px bg-border my-1"></div>
      <div
        v-else
        @click="handleAction(item)"
        class="px-4 py-2 flex items-center gap-3 cursor-pointer hover:bg-bg-tertiary text-sm transition-colors"
        :class="[
          item.disabled ? 'opacity-50 cursor-not-allowed' : '',
          item.danger
            ? 'text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20'
            : 'text-text-primary',
        ]"
      >
        <component
          v-if="item.icon && getIconComponent(item.icon)"
          :is="getIconComponent(item.icon)"
          :size="20"
          :weight="item.iconWeight || 'regular'"
          :class="
            item.iconColor ||
            (item.danger ? 'text-red-600 dark:text-red-400' : 'text-text-secondary')
          "
        />
        <span>{{ item.label }}</span>
      </div>
    </template>
  </div>
</template>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.1s ease-out;
}
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>
