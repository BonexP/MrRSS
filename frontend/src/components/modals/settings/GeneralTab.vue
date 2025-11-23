<script setup>
import { store } from '../../../store.js';
import { watch, onUnmounted, ref } from 'vue';

const props = defineProps({
    settings: { type: Object, required: true }
});

// Debounce timer to prevent excessive API calls
let saveTimeout = null;

// Track previous translation settings
const prevTranslationSettings = ref({
    enabled: props.settings.translation_enabled,
    targetLang: props.settings.target_language
});

// Auto-save function that saves settings immediately
async function autoSave() {
    try {
        // Check if translation settings changed
        const translationChanged = 
            prevTranslationSettings.value.enabled !== props.settings.translation_enabled ||
            (props.settings.translation_enabled && prevTranslationSettings.value.targetLang !== props.settings.target_language);
        
        await fetch('/api/settings', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                update_interval: props.settings.update_interval.toString(),
                translation_enabled: props.settings.translation_enabled.toString(),
                target_language: props.settings.target_language,
                translation_provider: props.settings.translation_provider,
                deepl_api_key: props.settings.deepl_api_key,
                auto_cleanup_enabled: props.settings.auto_cleanup_enabled.toString(),
                max_cache_size_mb: props.settings.max_cache_size_mb.toString(),
                max_article_age_days: props.settings.max_article_age_days.toString(),
                language: props.settings.language,
                theme: props.settings.theme,
                show_hidden_articles: props.settings.show_hidden_articles.toString()
            })
        });
        
        // Apply settings immediately
        store.i18n.setLocale(props.settings.language);
        store.setTheme(props.settings.theme);
        store.startAutoRefresh(props.settings.update_interval);
        
        // Clear and re-translate if translation settings changed
        if (translationChanged) {
            await fetch('/api/articles/clear-translations', { method: 'POST' });
            // Update tracking
            prevTranslationSettings.value = {
                enabled: props.settings.translation_enabled,
                targetLang: props.settings.target_language
            };
            // Refresh articles to show without translations, then re-translate if enabled
            store.fetchArticles();
        }
        
        // Refresh articles if show_hidden_articles changed
        if (props.settings.show_hidden_articles !== undefined) {
            store.fetchArticles();
        }
    } catch (e) {
        console.error('Error auto-saving settings:', e);
    }
}

// Debounced auto-save function
function debouncedAutoSave() {
    if (saveTimeout) {
        clearTimeout(saveTimeout);
    }
    saveTimeout = setTimeout(autoSave, 500); // Wait 500ms after last change
}

// Watch the entire settings object for changes
watch(() => props.settings, debouncedAutoSave, { deep: true });

// Clean up timeout on unmount to prevent memory leaks
onUnmounted(() => {
    if (saveTimeout) {
        clearTimeout(saveTimeout);
        saveTimeout = null;
    }
});

// Format last update time
function formatLastUpdate(timestamp) {
    if (!timestamp) return store.i18n.t('never');
    try {
        const date = new Date(timestamp);
        const now = new Date();
        const diffMs = now - date;
        const diffMins = Math.floor(diffMs / 60000);
        const diffHours = Math.floor(diffMs / 3600000);
        const diffDays = Math.floor(diffMs / 86400000);
        
        if (diffMins < 1) return store.i18n.t('justNow');
        if (diffMins < 60) return store.i18n.t('minutesAgo', { count: diffMins });
        if (diffHours < 24) return store.i18n.t('hoursAgo', { count: diffHours });
        if (diffDays < 7) return store.i18n.t('daysAgo', { count: diffDays });
        
        return date.toLocaleDateString(store.i18n.locale.value === 'zh' ? 'zh-CN' : 'en-US');
    } catch {
        return store.i18n.t('never');
    }
}
</script>

<template>
    <div class="space-y-6">
        <div class="setting-group">
            <label class="font-semibold mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2">
                <i class="ph ph-palette text-base"></i>
                {{ store.i18n.t('appearance') }}
            </label>
            <div class="setting-item">
                <div class="flex-1 flex items-start gap-3">
                    <i class="ph ph-moon text-xl text-text-secondary mt-0.5"></i>
                    <div class="flex-1">
                        <div class="font-medium mb-1">{{ store.i18n.t('theme') }}</div>
                        <div class="text-xs text-text-secondary">{{ store.i18n.t('themeDesc') }}</div>
                    </div>
                </div>
                <select v-model="settings.theme" class="input-field w-40">
                    <option value="light">{{ store.i18n.t('light') }}</option>
                    <option value="dark">{{ store.i18n.t('dark') }}</option>
                    <option value="auto">{{ store.i18n.t('auto') }}</option>
                </select>
            </div>
            <div class="setting-item mt-3">
                <div class="flex-1 flex items-start gap-3">
                    <i class="ph ph-translate text-xl text-text-secondary mt-0.5"></i>
                    <div class="flex-1">
                        <div class="font-medium mb-1">{{ store.i18n.t('language') }}</div>
                        <div class="text-xs text-text-secondary">{{ store.i18n.t('languageDesc') }}</div>
                    </div>
                </div>
                <select v-model="settings.language" class="input-field w-32">
                    <option value="en">{{ store.i18n.t('english') }}</option>
                    <option value="zh">{{ store.i18n.t('chinese') }}</option>
                </select>
            </div>
        </div>

        <div class="setting-group">
            <label class="font-semibold mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2">
                <i class="ph ph-arrow-clockwise text-base"></i>
                {{ store.i18n.t('updates') }}
            </label>
            <div class="setting-item">
                <div class="flex-1 flex items-start gap-3">
                    <i class="ph ph-clock text-xl text-text-secondary mt-0.5"></i>
                    <div class="flex-1">
                        <div class="font-medium mb-1">{{ store.i18n.t('autoUpdateInterval') }}</div>
                        <div class="text-xs text-text-secondary">{{ store.i18n.t('autoUpdateIntervalDesc') }}</div>
                    </div>
                </div>
                <input type="number" v-model="settings.update_interval" min="1" class="input-field w-20 text-center">
            </div>
            <div class="setting-item mt-3">
                <div class="flex-1 flex items-start gap-3">
                    <i class="ph ph-calendar-check text-xl text-text-secondary mt-0.5"></i>
                    <div class="flex-1">
                        <div class="font-medium mb-1">{{ store.i18n.t('lastArticleUpdate') }}</div>
                        <div class="text-xs text-text-secondary">{{ store.i18n.t('lastArticleUpdateDesc') }}</div>
                    </div>
                </div>
                <div class="text-sm font-medium text-accent">{{ formatLastUpdate(settings.last_article_update) }}</div>
            </div>
        </div>

        <div class="setting-group">
            <label class="font-semibold mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2">
                <i class="ph ph-database text-base"></i>
                {{ store.i18n.t('database') }}
            </label>
            <div class="setting-item">
                <div class="flex-1 flex items-start gap-3">
                    <i class="ph ph-broom text-xl text-text-secondary mt-0.5"></i>
                    <div class="flex-1">
                        <div class="font-medium mb-1">{{ store.i18n.t('autoCleanup') }}</div>
                        <div class="text-xs text-text-secondary">{{ store.i18n.t('autoCleanupDesc') }}</div>
                    </div>
                </div>
                <input type="checkbox" v-model="settings.auto_cleanup_enabled" class="toggle">
            </div>
            
            <div v-if="settings.auto_cleanup_enabled" class="ml-4 mt-3 space-y-3 border-l-2 border-border pl-4">
                <div class="sub-setting-item">
                    <div class="flex-1 flex items-start gap-3">
                        <i class="ph ph-hard-drive text-lg text-text-secondary mt-0.5"></i>
                        <div class="flex-1">
                            <div class="text-sm font-medium mb-1">{{ store.i18n.t('maxCacheSize') }}</div>
                            <div class="text-xs text-text-secondary">{{ store.i18n.t('maxCacheSizeDesc') }}</div>
                        </div>
                    </div>
                    <div class="flex items-center gap-2">
                        <input type="number" v-model="settings.max_cache_size_mb" min="1" max="1000" class="input-field w-20 text-center">
                        <span class="text-sm text-text-secondary">MB</span>
                    </div>
                </div>
                
                <div class="sub-setting-item">
                    <div class="flex-1 flex items-start gap-3">
                        <i class="ph ph-calendar-x text-lg text-text-secondary mt-0.5"></i>
                        <div class="flex-1">
                            <div class="text-sm font-medium mb-1">{{ store.i18n.t('maxArticleAge') }}</div>
                            <div class="text-xs text-text-secondary">{{ store.i18n.t('maxArticleAgeDesc') }}</div>
                        </div>
                    </div>
                    <div class="flex items-center gap-2">
                        <input type="number" v-model="settings.max_article_age_days" min="1" max="365" class="input-field w-20 text-center">
                        <span class="text-sm text-text-secondary">{{ store.i18n.t('days') }}</span>
                    </div>
                </div>
            </div>
            
            <div class="setting-item mt-3">
                <div class="flex-1 flex items-start gap-3">
                    <i class="ph ph-eye-slash text-xl text-text-secondary mt-0.5"></i>
                    <div class="flex-1">
                        <div class="font-medium mb-1">{{ store.i18n.t('showHiddenArticles') }}</div>
                        <div class="text-xs text-text-secondary">{{ store.i18n.t('showHiddenArticlesDesc') }}</div>
                    </div>
                </div>
                <input type="checkbox" v-model="settings.show_hidden_articles" class="toggle">
            </div>
        </div>

        <div class="setting-group">
            <label class="font-semibold mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2">
                <i class="ph ph-globe text-base"></i>
                {{ store.i18n.t('translation') }}
            </label>
            <div class="setting-item mb-4">
                <div class="flex-1 flex items-start gap-3">
                    <i class="ph ph-article text-xl text-text-secondary mt-0.5"></i>
                    <div class="flex-1">
                        <div class="font-medium mb-1">{{ store.i18n.t('enableTranslation') }}</div>
                        <div class="text-xs text-text-secondary">{{ store.i18n.t('enableTranslationDesc') }}</div>
                    </div>
                </div>
                <input type="checkbox" v-model="settings.translation_enabled" class="toggle">
            </div>
            
            <div v-if="settings.translation_enabled" class="ml-4 space-y-3 border-l-2 border-border pl-4">
                <div class="sub-setting-item">
                    <label class="text-sm font-medium mb-1 block">{{ store.i18n.t('translationProvider') }}</label>
                    <select v-model="settings.translation_provider" class="input-field w-full">
                        <option value="google">Google Translate (Free)</option>
                        <option value="deepl">DeepL API</option>
                    </select>
                </div>
                <div v-if="settings.translation_provider === 'deepl'" class="sub-setting-item">
                    <label class="text-sm font-medium mb-1 block">{{ store.i18n.t('deeplApiKey') }}</label>
                    <input type="password" v-model="settings.deepl_api_key" :placeholder="store.i18n.t('deeplApiKeyPlaceholder')" class="input-field w-full">
                </div>
                <div class="sub-setting-item">
                    <label class="text-sm font-medium mb-1 block">{{ store.i18n.t('targetLanguage') }}</label>
                    <select v-model="settings.target_language" class="input-field w-full">
                        <option value="en">{{ store.i18n.t('english') }}</option>
                        <option value="es">{{ store.i18n.t('spanish') }}</option>
                        <option value="fr">{{ store.i18n.t('french') }}</option>
                        <option value="de">{{ store.i18n.t('german') }}</option>
                        <option value="zh">{{ store.i18n.t('chinese') }}</option>
                        <option value="ja">{{ store.i18n.t('japanese') }}</option>
                    </select>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.input-field {
    @apply p-2.5 border border-border rounded-md bg-bg-secondary text-text-primary text-sm focus:border-accent focus:outline-none transition-colors;
}
.toggle {
    @apply w-10 h-5 appearance-none bg-bg-tertiary rounded-full relative cursor-pointer border border-border transition-colors checked:bg-accent checked:border-accent;
}
.toggle::after {
    content: '';
    @apply absolute top-0.5 left-0.5 w-3.5 h-3.5 bg-white rounded-full shadow-sm transition-transform;
}
.toggle:checked::after {
    transform: translateX(20px);
}
.setting-item {
    @apply flex items-start justify-between gap-4 p-3 rounded-lg bg-bg-secondary border border-border;
}
.sub-setting-item {
    @apply p-2.5 rounded-md bg-bg-tertiary;
}
</style>
