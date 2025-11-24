<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import * as PhosphorIcons from "@phosphor-icons/vue";

const props = defineProps({
    items: { type: Array, required: true }, // [{ label: 'Edit', action: 'edit', icon: 'PhPencil' }, { separator: true }]
    x: { type: Number, required: true },
    y: { type: Number, required: true }
});

const emit = defineEmits(['close', 'action']);
const menuRef = ref(null);

// Map old icon names to new component names
const iconMap = {
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
    'ph-arrow-square-out': 'PhArrowSquareOut'
};

// Get icon component from icon string
function getIconComponent(iconName) {
    if (!iconName) return null;
    const componentName = iconMap[iconName] || iconName;
    return PhosphorIcons[componentName] || null;
}

function handleClickOutside(event) {
    if (menuRef.value && !menuRef.value.contains(event.target)) {
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

function handleAction(item) {
    if (item.disabled) return;
    emit('action', item.action);
    emit('close');
}
</script>

<template>
    <div ref="menuRef" class="fixed z-50 bg-bg-primary border border-border rounded-lg shadow-xl py-1 min-w-[180px] animate-fade-in"
         :style="{ top: `${y}px`, left: `${x}px` }">
        <template v-for="(item, index) in items" :key="index">
            <div v-if="item.separator" class="h-px bg-border my-1"></div>
            <div v-else 
                 @click="handleAction(item)"
                 class="px-4 py-2 flex items-center gap-3 cursor-pointer hover:bg-bg-tertiary text-sm transition-colors"
                 :class="[
                     item.disabled ? 'opacity-50 cursor-not-allowed' : '',
                     item.danger ? 'text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20' : 'text-text-primary'
                 ]">
                <component v-if="item.icon && getIconComponent(item.icon)" 
                           :is="getIconComponent(item.icon)" 
                           :size="20" 
                           :weight="item.iconWeight || 'regular'"
                           :class="item.iconColor || (item.danger ? 'text-red-600 dark:text-red-400' : 'text-text-secondary')" />
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
    from { opacity: 0; transform: scale(0.95); }
    to { opacity: 1; transform: scale(1); }
}
</style>
