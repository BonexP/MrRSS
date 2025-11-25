<script setup>
import { store } from '../../../store.js';
import { ref, computed, onUnmounted } from 'vue';
import { 
    PhHardDrives, PhUpload, PhDownload, PhBroom, PhRss, PhPlus, 
    PhTrash, PhFolder, PhPencil, PhMagnifyingGlass, PhCircleNotch
} from "@phosphor-icons/vue";

const emit = defineEmits(['import-opml', 'export-opml', 'cleanup-database', 'add-feed', 'edit-feed', 'delete-feed', 'batch-delete', 'batch-move']);

const selectedFeeds = ref([]);
const isDiscoveringAll = ref(false);
const discoveryProgress = ref({ message: '', detail: '', current: 0, total: 0, feedName: '', foundCount: 0 });
let pollInterval = null;

const isAllSelected = computed(() => {
    return store.feeds && store.feeds.length > 0 && selectedFeeds.value.length === store.feeds.length;
});

function toggleSelectAll(e) {
    if (!store.feeds) return;
    if (e.target.checked) {
        selectedFeeds.value = store.feeds.map(f => f.id);
    } else {
        selectedFeeds.value = [];
    }
}

function handleImportOPML(event) {
    emit('import-opml', event);
}

function handleExportOPML() {
    emit('export-opml');
}

function handleCleanupDatabase() {
    emit('cleanup-database');
}

function getHostname(url) {
    try {
        return new URL(url).hostname;
    } catch {
        return url;
    }
}

async function handleDiscoverAll() {
    isDiscoveringAll.value = true;
    discoveryProgress.value = { message: store.i18n.t('preparingDiscovery'), detail: '', current: 0, total: 0, feedName: '', foundCount: 0 };
    
    // Clear any existing poll interval
    if (pollInterval) {
        clearInterval(pollInterval);
        pollInterval = null;
    }

    try {
        // Clear any previous discovery state
        await fetch('/api/feeds/discover-all/clear', { method: 'POST' });

        // Start batch discovery in background
        const startResponse = await fetch('/api/feeds/discover-all/start', {
            method: 'POST'
        });

        if (!startResponse.ok) {
            const errorText = await startResponse.text();
            throw new Error(errorText || 'Failed to start batch discovery');
        }

        const startResult = await startResponse.json();
        
        // Check if already complete (all feeds discovered)
        if (startResult.status === 'complete') {
            window.showToast(startResult.message || store.i18n.t('noFriendLinksFound'), 'info');
            isDiscoveringAll.value = false;
            return;
        }

        discoveryProgress.value.total = startResult.total || 0;

        // Start polling for progress
        pollInterval = setInterval(async () => {
            try {
                const progressResponse = await fetch('/api/feeds/discover-all/progress');
                if (!progressResponse.ok) {
                    throw new Error('Failed to get progress');
                }

                const state = await progressResponse.json();
                console.log('Batch progress state:', state);

                // Update progress display
                if (state.progress) {
                    const progress = state.progress;
                    discoveryProgress.value = {
                        message: progress.message || store.i18n.t('processingFeed', { current: progress.current || 0, total: progress.total || 0 }),
                        detail: progress.feed_name || (progress.detail ? getHostname(progress.detail) : ''),
                        current: progress.current || 0,
                        total: progress.total || 0,
                        feedName: progress.feed_name || '',
                        foundCount: progress.found_count || 0
                    };
                }

                // Check if complete
                if (state.is_complete) {
                    clearInterval(pollInterval);
                    pollInterval = null;

                    // Refresh feeds to show updated discovery status
                    await store.fetchFeeds();

                    if (state.error) {
                        window.showToast(store.i18n.t('discoveryFailed') + ': ' + state.error, 'error');
                    } else if (state.feeds && state.feeds.length > 0) {
                        window.showToast(
                            store.i18n.t('discoveryComplete') + ': ' + 
                            store.i18n.t('foundFeeds', { count: state.feeds.length }),
                            'success'
                        );
                    } else {
                        window.showToast(store.i18n.t('noFriendLinksFound'), 'info');
                    }

                    isDiscoveringAll.value = false;
                    discoveryProgress.value = { message: '', detail: '', current: 0, total: 0, feedName: '', foundCount: 0 };

                    // Clear the discovery state
                    await fetch('/api/feeds/discover-all/clear', { method: 'POST' });
                }
            } catch (pollError) {
                console.error('Polling error:', pollError);
                // Don't stop polling on transient errors
            }
        }, 500); // Poll every 500ms

    } catch (error) {
        console.error('Discovery error:', error);
        window.showToast(store.i18n.t('discoveryFailed'), 'error');
        isDiscoveringAll.value = false;
        discoveryProgress.value = { message: '', detail: '', current: 0, total: 0, feedName: '', foundCount: 0 };
    }
}

function handleAddFeed() {
    emit('add-feed');
}

function handleEditFeed(feed) {
    emit('edit-feed', feed);
}

function handleDeleteFeed(id) {
    emit('delete-feed', id);
}

function handleBatchDelete() {
    if (selectedFeeds.value.length === 0) return;
    emit('batch-delete', selectedFeeds.value);
    selectedFeeds.value = [];
}

function handleBatchMove() {
    if (selectedFeeds.value.length === 0) return;
    emit('batch-move', selectedFeeds.value);
    selectedFeeds.value = [];
}

function getFavicon(url) {
    try {
        return `https://www.google.com/s2/favicons?domain=${new URL(url).hostname}`;
    } catch {
        return '';
    }
}

// Cleanup on unmount
onUnmounted(() => {
    if (pollInterval) {
        clearInterval(pollInterval);
        pollInterval = null;
    }
    // Clear discovery state on server
    fetch('/api/feeds/discover-all/clear', { method: 'POST' }).catch(() => {});
});
</script>

<template>
    <div class="space-y-4 sm:space-y-6">
        <div class="setting-group">
            <label class="font-semibold mb-2 sm:mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2">
                <PhHardDrives :size="14" class="sm:w-4 sm:h-4" />
                {{ store.i18n.t('dataManagement') }}
            </label>
            <div class="flex flex-col sm:flex-row gap-2 sm:gap-3 mb-2 sm:mb-3">
                <button @click="$refs.opmlInput.click()" class="btn-secondary flex-1 justify-center text-sm sm:text-base">
                    <PhUpload :size="18" class="sm:w-5 sm:h-5" /> {{ store.i18n.t('importOPML') }}
                </button>
                <input type="file" ref="opmlInput" class="hidden" @change="handleImportOPML">
                <button @click="handleExportOPML" class="btn-secondary flex-1 justify-center text-sm sm:text-base">
                    <PhDownload :size="18" class="sm:w-5 sm:h-5" /> {{ store.i18n.t('exportOPML') }}
                </button>
            </div>
            <div class="flex flex-col sm:flex-row gap-2 sm:gap-3 mb-2 sm:mb-3">
                <button @click="handleDiscoverAll" :disabled="isDiscoveringAll" 
                        :class="['btn-primary flex-1 justify-center text-sm sm:text-base', isDiscoveringAll && 'opacity-50 cursor-not-allowed']">
                    <PhCircleNotch v-if="isDiscoveringAll" :size="18" class="sm:w-5 sm:h-5 animate-spin" />
                    <PhMagnifyingGlass v-else :size="18" class="sm:w-5 sm:h-5" />
                    {{ isDiscoveringAll ? store.i18n.t('discoveringAllFeeds') : store.i18n.t('discoverAllFeeds') }}
                </button>
            </div>
            <p v-if="!isDiscoveringAll" class="text-xs text-text-secondary mb-2">
                {{ store.i18n.t('discoverAllFeedsDesc') }}
            </p>
            <div v-if="isDiscoveringAll" class="bg-gradient-to-r from-accent/10 to-accent/5 border-l-4 border-accent rounded-lg p-4 mb-2 space-y-2">
                <div class="flex items-center gap-2">
                    <PhCircleNotch :size="20" class="animate-spin text-accent shrink-0" />
                    <div class="flex-1 min-w-0">
                        <p class="text-sm font-semibold text-accent">
                            {{ discoveryProgress.message }}
                        </p>
                        <p v-if="discoveryProgress.detail" class="text-xs text-text-secondary mt-1 truncate">
                            {{ discoveryProgress.detail }}
                        </p>
                    </div>
                </div>
                <div v-if="discoveryProgress.total > 0" class="w-full bg-bg-tertiary rounded-full h-1.5 overflow-hidden">
                    <div class="bg-accent h-full transition-all duration-300" :style="{ width: (discoveryProgress.current / discoveryProgress.total * 100) + '%' }"></div>
                </div>
                <div class="flex items-center justify-between text-xs text-text-tertiary pt-2 border-t border-accent/20">
                    <div class="flex items-center gap-2">
                        <PhRss :size="12" />
                        <span v-if="discoveryProgress.total > 0">{{ discoveryProgress.current }}/{{ discoveryProgress.total }}</span>
                        <span v-else>{{ store.i18n.t('pleaseWait') }}</span>
                    </div>
                    <span v-if="discoveryProgress.foundCount > 0">{{ store.i18n.t('foundSoFar', { count: discoveryProgress.foundCount }) }}</span>
                </div>
            </div>
            <div class="flex">
                <button @click="handleCleanupDatabase" class="btn-danger flex-1 justify-center text-sm sm:text-base">
                    <PhBroom :size="18" class="sm:w-5 sm:h-5" /> {{ store.i18n.t('cleanDatabase') }}
                </button>
            </div>
            <p class="text-xs text-text-secondary mt-2">
                {{ store.i18n.t('cleanDatabaseDesc') }}
            </p>
        </div>
        
        <div class="setting-group">
            <label class="font-semibold mb-2 sm:mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2">
                <PhRss :size="14" class="sm:w-4 sm:h-4" />
                {{ store.i18n.t('manageFeeds') }}
            </label>
            
            <div class="flex flex-wrap gap-1.5 sm:gap-2 mb-2 text-xs sm:text-sm">
                <button @click="handleAddFeed" class="btn-secondary py-1.5 px-2.5 sm:px-3">
                    <PhPlus :size="14" class="sm:w-4 sm:h-4" /> <span class="hidden sm:inline">{{ store.i18n.t('addFeed') }}</span><span class="sm:hidden">{{ store.i18n.t('addFeed').split(' ')[0] }}</span>
                </button>
                <button @click="handleBatchDelete" class="btn-danger py-1.5 px-2.5 sm:px-3" :disabled="selectedFeeds.length === 0">
                    <PhTrash :size="14" class="sm:w-4 sm:h-4" /> <span class="hidden sm:inline">{{ store.i18n.t('deleteSelected') }}</span><span class="sm:hidden">{{ store.i18n.t('delete') }}</span>
                </button>
                <button @click="handleBatchMove" class="btn-secondary py-1.5 px-2.5 sm:px-3" :disabled="selectedFeeds.length === 0">
                    <PhFolder :size="14" class="sm:w-4 sm:h-4" /> <span class="hidden sm:inline">{{ store.i18n.t('moveSelected') }}</span><span class="sm:hidden">{{ store.i18n.t('move') }}</span>
                </button>
                <div class="flex-1 min-w-0"></div>
                <label class="flex items-center gap-1.5 sm:gap-2 cursor-pointer select-none whitespace-nowrap">
                    <input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll" class="w-3.5 h-3.5 sm:w-4 sm:h-4 rounded border-border text-accent focus:ring-2 focus:ring-accent cursor-pointer">
                    <span class="hidden sm:inline">{{ store.i18n.t('selectAll') }}</span>
                </label>
            </div>

            <div class="border border-border rounded-lg bg-bg-secondary overflow-y-auto max-h-60 sm:max-h-96">
                <div v-for="feed in store.feeds" :key="feed.id" class="flex items-center p-1.5 sm:p-2 border-b border-border last:border-0 bg-bg-primary hover:bg-bg-secondary gap-1.5 sm:gap-2">
                    <input type="checkbox" :value="feed.id" v-model="selectedFeeds" class="w-3.5 h-3.5 sm:w-4 sm:h-4 shrink-0 rounded border-border text-accent focus:ring-2 focus:ring-accent cursor-pointer">
                    <div class="w-4 h-4 flex items-center justify-center shrink-0">
                        <img :src="feed.image_url || getFavicon(feed.url)" class="w-full h-full object-contain" @error="$event.target.style.display='none'">
                    </div>
                    <div class="truncate flex-1 min-w-0">
                        <div class="font-medium truncate text-xs sm:text-sm">{{ feed.title }}</div>
                        <div class="text-xs text-text-secondary truncate hidden sm:block">
                            <span v-if="feed.category" class="inline-flex items-center gap-1">
                                <PhFolder :size="10" class="inline" />
                                {{ feed.category }}
                                <span class="mx-1">•</span>
                            </span>
                            <span>{{ feed.url }}</span>
                        </div>
                    </div>
                    <div class="flex gap-0.5 sm:gap-1 shrink-0">
                        <button @click="handleEditFeed(feed)" class="text-accent hover:bg-bg-tertiary p-1 rounded text-sm" :title="store.i18n.t('edit')"><PhPencil :size="14" class="sm:w-4 sm:h-4" /></button>
                        <button @click="handleDeleteFeed(feed.id)" class="text-red-500 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 p-1 rounded text-sm" :title="store.i18n.t('delete')"><PhTrash :size="14" class="sm:w-4 sm:h-4" /></button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.btn-primary {
    @apply bg-accent text-white px-3 sm:px-4 py-1.5 sm:py-2 rounded-md cursor-pointer flex items-center gap-1.5 sm:gap-2 font-semibold hover:bg-accent-hover transition-colors shadow-sm;
}
.btn-primary:disabled {
    @apply opacity-50 cursor-not-allowed;
}
.btn-secondary {
    @apply bg-transparent border border-border text-text-primary px-3 sm:px-4 py-1.5 sm:py-2 rounded-md cursor-pointer flex items-center gap-1.5 sm:gap-2 font-medium hover:bg-bg-tertiary transition-colors;
}
.btn-secondary:disabled {
    @apply opacity-50 cursor-not-allowed;
}
.btn-danger {
    @apply bg-transparent border border-red-300 text-red-600 px-3 sm:px-4 py-1.5 sm:py-2 rounded-md cursor-pointer flex items-center gap-1.5 sm:gap-2 font-semibold hover:bg-red-50 dark:hover:bg-red-900/20 dark:border-red-400 dark:text-red-400 transition-colors;
}
.btn-danger:disabled {
    @apply opacity-50 cursor-not-allowed;
}
</style>
