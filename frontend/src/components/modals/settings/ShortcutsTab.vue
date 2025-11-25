<script setup>
import { store } from '../../../store.js';
import { ref, onMounted, onUnmounted, computed, watch } from 'vue';
import { 
    PhKeyboard, PhArrowDown, PhArrowUp, PhArrowRight, PhX, PhBookOpen,
    PhStar, PhArrowSquareOut, PhArticle, PhArrowClockwise, PhCheckCircle, 
    PhGear, PhPlus, PhMagnifyingGlass, PhListDashes, PhCircle, PhHeart,
    PhArrowCounterClockwise
} from "@phosphor-icons/vue";

const props = defineProps({
    settings: { type: Object, required: true }
});

// Default shortcuts configuration
const defaultShortcuts = {
    nextArticle: 'j',
    previousArticle: 'k',
    openArticle: 'Enter',
    closeArticle: 'Escape',
    toggleReadStatus: 'r',
    toggleFavoriteStatus: 's',
    openInBrowser: 'o',
    toggleContentView: 'v',
    refreshFeeds: 'Shift+r',
    markAllRead: 'Shift+a',
    openSettings: ',',
    addFeed: 'a',
    focusSearch: '/',
    goToAllArticles: '1',
    goToUnread: '2',
    goToFavorites: '3'
};

// Current shortcuts (loaded from settings or use defaults)
const shortcuts = ref({ ...defaultShortcuts });

// Track which shortcut is being edited
const editingShortcut = ref(null);
const recordedKey = ref('');

// Shortcut groups for display
const shortcutGroups = computed(() => [
    {
        label: store.i18n.t('shortcutNavigation'),
        items: [
            { key: 'nextArticle', label: store.i18n.t('nextArticle'), icon: PhArrowDown },
            { key: 'previousArticle', label: store.i18n.t('previousArticle'), icon: PhArrowUp },
            { key: 'openArticle', label: store.i18n.t('openArticle'), icon: PhArrowRight },
            { key: 'closeArticle', label: store.i18n.t('closeArticle'), icon: PhX },
            { key: 'goToAllArticles', label: store.i18n.t('goToAllArticles'), icon: PhListDashes },
            { key: 'goToUnread', label: store.i18n.t('goToUnread'), icon: PhCircle },
            { key: 'goToFavorites', label: store.i18n.t('goToFavorites'), icon: PhHeart }
        ]
    },
    {
        label: store.i18n.t('shortcutArticles'),
        items: [
            { key: 'toggleReadStatus', label: store.i18n.t('toggleReadStatus'), icon: PhBookOpen },
            { key: 'toggleFavoriteStatus', label: store.i18n.t('toggleFavoriteStatus'), icon: PhStar },
            { key: 'openInBrowser', label: store.i18n.t('openInBrowserShortcut'), icon: PhArrowSquareOut },
            { key: 'toggleContentView', label: store.i18n.t('toggleContentView'), icon: PhArticle }
        ]
    },
    {
        label: store.i18n.t('shortcutOther'),
        items: [
            { key: 'refreshFeeds', label: store.i18n.t('refreshFeedsShortcut'), icon: PhArrowClockwise },
            { key: 'markAllRead', label: store.i18n.t('markAllReadShortcut'), icon: PhCheckCircle },
            { key: 'openSettings', label: store.i18n.t('openSettingsShortcut'), icon: PhGear },
            { key: 'addFeed', label: store.i18n.t('addFeedShortcut'), icon: PhPlus },
            { key: 'focusSearch', label: store.i18n.t('focusSearch'), icon: PhMagnifyingGlass }
        ]
    }
]);

// Load shortcuts from settings
onMounted(() => {
    if (props.settings.shortcuts) {
        try {
            const parsed = typeof props.settings.shortcuts === 'string' 
                ? JSON.parse(props.settings.shortcuts) 
                : props.settings.shortcuts;
            shortcuts.value = { ...defaultShortcuts, ...parsed };
        } catch (e) {
            console.error('Error parsing shortcuts:', e);
            shortcuts.value = { ...defaultShortcuts };
        }
    }
    
    // Add global keyboard listener for recording
    window.addEventListener('keydown', handleKeyRecord, true);
});

onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyRecord, true);
});

// Format key for display
function formatKey(key) {
    if (!key) return '—';
    
    // Convert key combinations to display format
    const parts = key.split('+');
    return parts.map(part => {
        // Capitalize first letter and handle special keys
        if (part === 'Shift') return '⇧';
        if (part === 'Control' || part === 'Ctrl') return '⌃';
        if (part === 'Alt') return '⌥';
        if (part === 'Meta' || part === 'Cmd') return '⌘';
        if (part === 'Enter') return '↵';
        if (part === 'Escape') return 'Esc';
        if (part === 'ArrowUp') return '↑';
        if (part === 'ArrowDown') return '↓';
        if (part === 'ArrowLeft') return '←';
        if (part === 'ArrowRight') return '→';
        if (part === 'Space') return '␣';
        if (part === 'Backspace') return '⌫';
        if (part === 'Delete') return 'Del';
        if (part === 'Tab') return '⇥';
        return part.toUpperCase();
    }).join(' + ');
}

// Start editing a shortcut
function startEditing(shortcutKey) {
    editingShortcut.value = shortcutKey;
    recordedKey.value = '';
}

// Stop editing
function stopEditing() {
    editingShortcut.value = null;
    recordedKey.value = '';
}

// Handle key recording
function handleKeyRecord(e) {
    if (!editingShortcut.value) return;
    
    e.preventDefault();
    e.stopPropagation();
    
    // Handle Escape to clear the shortcut
    if (e.key === 'Escape' && !e.shiftKey && !e.ctrlKey && !e.altKey && !e.metaKey) {
        // Clear the shortcut
        shortcuts.value[editingShortcut.value] = '';
        saveShortcuts();
        window.showToast(store.i18n.t('shortcutCleared'), 'info');
        stopEditing();
        return;
    }
    
    // Build key combination
    let key = '';
    if (e.ctrlKey) key += 'Ctrl+';
    if (e.altKey) key += 'Alt+';
    if (e.shiftKey) key += 'Shift+';
    if (e.metaKey) key += 'Meta+';
    
    // Get the actual key
    let actualKey = e.key;
    
    // Skip modifier keys alone
    if (['Control', 'Alt', 'Shift', 'Meta'].includes(actualKey)) {
        return;
    }
    
    // Normalize key names
    if (actualKey === ' ') actualKey = 'Space';
    else if (actualKey.length === 1) actualKey = actualKey.toLowerCase();
    
    key += actualKey;
    
    // Check for conflicts
    const conflictKey = Object.entries(shortcuts.value).find(
        ([k, v]) => v === key && k !== editingShortcut.value
    );
    
    if (conflictKey) {
        window.showToast(store.i18n.t('shortcutConflict'), 'warning');
        stopEditing();
        return;
    }
    
    // Update the shortcut
    shortcuts.value[editingShortcut.value] = key;
    saveShortcuts();
    window.showToast(store.i18n.t('shortcutUpdated'), 'success');
    stopEditing();
}

// Save shortcuts to settings
async function saveShortcuts() {
    try {
        // Update props.settings.shortcuts
        props.settings.shortcuts = JSON.stringify(shortcuts.value);
        
        // The parent component will handle auto-save via the watcher
        // But we also dispatch an event to notify the app
        window.dispatchEvent(new CustomEvent('shortcuts-changed', {
            detail: { shortcuts: shortcuts.value }
        }));
    } catch (e) {
        console.error('Error saving shortcuts:', e);
    }
}

// Reset all shortcuts to defaults
function resetToDefaults() {
    shortcuts.value = { ...defaultShortcuts };
    saveShortcuts();
    window.showToast(store.i18n.t('shortcutUpdated'), 'success');
}

// Watch for settings changes from parent
watch(() => props.settings.shortcuts, (newVal) => {
    if (newVal) {
        try {
            const parsed = typeof newVal === 'string' ? JSON.parse(newVal) : newVal;
            shortcuts.value = { ...defaultShortcuts, ...parsed };
        } catch (e) {
            console.error('Error parsing shortcuts:', e);
        }
    }
}, { immediate: true });
</script>

<template>
    <div class="space-y-4 sm:space-y-6">
        <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
                <PhKeyboard :size="20" class="text-text-secondary sm:w-6 sm:h-6" />
                <div>
                    <h3 class="font-semibold text-sm sm:text-base">{{ store.i18n.t('shortcuts') }}</h3>
                    <p class="text-xs text-text-secondary">{{ store.i18n.t('shortcutsDesc') }}</p>
                </div>
            </div>
            <button @click="resetToDefaults" class="btn-secondary text-xs sm:text-sm py-1.5 px-3">
                <PhArrowCounterClockwise :size="14" class="sm:w-4 sm:h-4" />
                {{ store.i18n.t('resetToDefault') }}
            </button>
        </div>

        <div v-for="group in shortcutGroups" :key="group.label" class="setting-group">
            <label class="font-semibold mb-2 sm:mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2">
                {{ group.label }}
            </label>
            
            <div class="space-y-1">
                <div v-for="item in group.items" :key="item.key" class="shortcut-item">
                    <div class="flex items-center gap-2 flex-1 min-w-0">
                        <component :is="item.icon" :size="18" class="text-text-secondary shrink-0 sm:w-5 sm:h-5" />
                        <span class="text-sm sm:text-base truncate">{{ item.label }}</span>
                    </div>
                    
                    <button 
                        @click="startEditing(item.key)"
                        :class="['shortcut-key-btn', editingShortcut === item.key ? 'recording' : '']"
                    >
                        <span v-if="editingShortcut === item.key" class="text-accent animate-pulse">
                            {{ store.i18n.t('pressKey') }}
                        </span>
                        <span v-else>{{ formatKey(shortcuts[item.key]) }}</span>
                    </button>
                </div>
            </div>
        </div>
        
        <div class="text-xs text-text-secondary mt-4 p-3 bg-bg-secondary rounded-lg border border-border">
            <p class="mb-1">💡 {{ store.i18n.t('escToClear') }}</p>
        </div>
    </div>
</template>

<style scoped>
.shortcut-item {
    @apply flex items-center justify-between gap-2 p-2 sm:p-3 rounded-lg bg-bg-secondary border border-border;
}

.shortcut-key-btn {
    @apply px-3 py-1.5 rounded-md bg-bg-tertiary border border-border text-sm font-mono cursor-pointer transition-all min-w-[80px] text-center;
}

.shortcut-key-btn:hover {
    @apply border-accent bg-bg-primary;
}

.shortcut-key-btn.recording {
    @apply border-accent;
    background-color: rgba(59, 130, 246, 0.1);
}

.btn-secondary {
    @apply bg-transparent border border-border text-text-primary rounded-md cursor-pointer flex items-center gap-1.5 font-medium hover:bg-bg-tertiary transition-colors;
}

.animate-pulse {
    animation: pulse 1.5s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.5; }
}
</style>
