<script setup>
import { ref, computed } from 'vue';
import { store } from '../../store.js';
import { PhX, PhCheck, PhGlobe, PhRss, PhCircleNotch } from "@phosphor-icons/vue";

const props = defineProps({
    feed: { type: Object, required: true },
    show: { type: Boolean, required: true }
});

const emit = defineEmits(['close']);

const isDiscovering = ref(false);
const discoveredBlogs = ref([]);
const selectedBlogs = ref(new Set());
const errorMessage = ref('');

async function startDiscovery() {
    isDiscovering.value = true;
    errorMessage.value = '';
    discoveredBlogs.value = [];
    selectedBlogs.value.clear();

    try {
        const response = await fetch('/api/feeds/discover', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ feed_id: props.feed.id })
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to discover blogs');
        }

        const blogs = await response.json();
        discoveredBlogs.value = blogs || [];

        if (discoveredBlogs.value.length === 0) {
            errorMessage.value = store.i18n.t('noFriendLinksFound');
        }
    } catch (error) {
        console.error('Discovery error:', error);
        errorMessage.value = store.i18n.t('discoveryFailed') + ': ' + error.message;
    } finally {
        isDiscovering.value = false;
    }
}

function toggleBlogSelection(index) {
    if (selectedBlogs.value.has(index)) {
        selectedBlogs.value.delete(index);
    } else {
        selectedBlogs.value.add(index);
    }
}

function selectAll() {
    if (selectedBlogs.value.size === discoveredBlogs.value.length) {
        selectedBlogs.value.clear();
    } else {
        discoveredBlogs.value.forEach((_, index) => selectedBlogs.value.add(index));
    }
}

const hasSelection = computed(() => selectedBlogs.value.size > 0);
const allSelected = computed(() => discoveredBlogs.value.length > 0 && selectedBlogs.value.size === discoveredBlogs.value.length);

async function subscribeSelected() {
    if (!hasSelection.value) return;

    const subscribePromises = [];
    
    for (const index of selectedBlogs.value) {
        const blog = discoveredBlogs.value[index];
        const promise = fetch('/api/feeds/add', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                url: blog.rss_feed,
                category: props.feed.category || '',
                title: blog.name
            })
        });
        subscribePromises.push(promise);
    }

    try {
        await Promise.all(subscribePromises);
        await store.fetchFeeds();
        window.showToast(store.i18n.t('blogsSubscribedSuccess', { count: selectedBlogs.value.size }), 'success');
        emit('close');
    } catch (error) {
        console.error('Subscription error:', error);
        window.showToast(store.i18n.t('errorSubscribingBlogs'), 'error');
    }
}

function close() {
    emit('close');
}

// Start discovery when modal is shown
if (props.show) {
    startDiscovery();
}
</script>

<template>
    <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4" @click.self="close">
        <div class="bg-bg-primary w-full max-w-4xl max-h-[90vh] rounded-2xl shadow-2xl border border-border flex flex-col">
            <!-- Header -->
            <div class="flex justify-between items-center p-6 border-b border-border">
                <div>
                    <h2 class="text-xl font-bold text-text-primary">{{ store.i18n.t('discoverBlogs') }}</h2>
                    <p class="text-sm text-text-secondary mt-1">{{ store.i18n.t('fromFeed') }}: {{ feed.title }}</p>
                </div>
                <button @click="close" class="p-2 hover:bg-bg-tertiary rounded-lg transition-colors">
                    <PhX :size="24" class="text-text-secondary" />
                </button>
            </div>

            <!-- Content -->
            <div class="flex-1 overflow-y-auto p-6">
                <!-- Loading State -->
                <div v-if="isDiscovering" class="flex flex-col items-center justify-center py-12">
                    <PhCircleNotch :size="48" class="text-accent animate-spin mb-4" />
                    <p class="text-text-secondary">{{ store.i18n.t('discovering') }}</p>
                </div>

                <!-- Error State -->
                <div v-else-if="errorMessage" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 text-red-600 dark:text-red-400">
                    {{ errorMessage }}
                </div>

                <!-- Results -->
                <div v-else-if="discoveredBlogs.length > 0">
                    <div class="mb-4 flex items-center justify-between">
                        <p class="text-sm text-text-secondary">
                            {{ store.i18n.t('foundBlogs', { count: discoveredBlogs.length }) }}
                        </p>
                        <button @click="selectAll" class="text-sm text-accent hover:text-accent-hover font-medium">
                            {{ allSelected ? store.i18n.t('deselectAll') : store.i18n.t('selectAll') }}
                        </button>
                    </div>

                    <div class="space-y-3">
                        <div v-for="(blog, index) in discoveredBlogs" :key="index" 
                             @click="toggleBlogSelection(index)"
                             :class="[
                                 'border border-border rounded-lg p-4 cursor-pointer transition-all',
                                 selectedBlogs.has(index) 
                                     ? 'bg-accent/10 border-accent ring-2 ring-accent/20' 
                                     : 'bg-bg-secondary hover:bg-bg-tertiary'
                             ]">
                            <div class="flex items-start gap-4">
                                <!-- Checkbox -->
                                <div class="mt-1 shrink-0">
                                    <div :class="[
                                        'w-5 h-5 rounded border-2 flex items-center justify-center transition-all',
                                        selectedBlogs.has(index) 
                                            ? 'bg-accent border-accent' 
                                            : 'border-border bg-bg-primary'
                                    ]">
                                        <PhCheck v-if="selectedBlogs.has(index)" :size="14" weight="bold" class="text-white" />
                                    </div>
                                </div>

                                <!-- Blog Info -->
                                <div class="flex-1 min-w-0">
                                    <div class="flex items-start gap-3 mb-2">
                                        <img :src="blog.icon_url" class="w-8 h-8 rounded" :alt="blog.name" @error="$event.target.style.display='none'">
                                        <div class="flex-1 min-w-0">
                                            <h3 class="font-semibold text-text-primary truncate">{{ blog.name }}</h3>
                                            <a :href="blog.homepage" 
                                               target="_blank" 
                                               @click.stop
                                               class="text-xs text-accent hover:text-accent-hover flex items-center gap-1 mt-1 truncate">
                                                <PhGlobe :size="14" />
                                                <span class="truncate">{{ blog.homepage }}</span>
                                            </a>
                                        </div>
                                    </div>

                                    <!-- Recent Articles -->
                                    <div v-if="blog.recent_articles && blog.recent_articles.length > 0" class="mt-3 space-y-1.5">
                                        <p class="text-xs font-medium text-text-secondary mb-2">{{ store.i18n.t('recentArticles') }}:</p>
                                        <div v-for="(article, aIndex) in blog.recent_articles" :key="aIndex" 
                                             class="text-sm text-text-secondary pl-4 truncate border-l-2 border-border">
                                            {{ article }}
                                        </div>
                                    </div>

                                    <!-- RSS Feed URL -->
                                    <div class="mt-3 flex items-center gap-1 text-xs text-text-tertiary">
                                        <PhRss :size="12" />
                                        <span class="truncate">{{ blog.rss_feed }}</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Empty State -->
                <div v-else class="text-center py-12">
                    <p class="text-text-secondary">{{ store.i18n.t('startDiscovery') }}</p>
                </div>
            </div>

            <!-- Footer -->
            <div class="flex justify-between items-center p-6 border-t border-border">
                <button @click="close" class="btn-secondary">
                    {{ store.i18n.t('cancel') }}
                </button>
                <button @click="subscribeSelected" 
                        :disabled="!hasSelection" 
                        :class="['btn-primary', !hasSelection && 'opacity-50 cursor-not-allowed']">
                    {{ store.i18n.t('subscribeSelected') }} 
                    <span v-if="hasSelection">({{ selectedBlogs.size }})</span>
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.btn-primary {
    @apply px-4 py-2 bg-accent text-white rounded-lg hover:bg-accent-hover transition-colors font-medium;
}

.btn-secondary {
    @apply px-4 py-2 bg-bg-tertiary text-text-primary rounded-lg hover:opacity-80 transition-colors;
}
</style>
