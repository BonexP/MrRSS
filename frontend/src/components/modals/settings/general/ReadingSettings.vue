<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import {
  PhArticle,
  PhImage,
  PhImages,
  PhArticleNyTimes,
  PhCursorClick,
  PhEyeSlash,
} from '@phosphor-icons/vue';
import type { SettingsData } from '@/types/settings';

const { t } = useI18n();

interface Props {
  settings: SettingsData;
}

defineProps<Props>();

const emit = defineEmits<{
  'update:settings': [settings: SettingsData];
}>();
</script>

<template>
  <div class="setting-group">
    <label
      class="font-semibold mb-2 sm:mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2"
    >
      <PhArticle :size="14" class="sm:w-4 sm:h-4" />
      {{ t('readingAndDisplay') }}
    </label>

    <div class="setting-item">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhArticle :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('defaultViewMode') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('defaultViewModeDesc') }}
          </div>
        </div>
      </div>
      <select
        :value="settings.default_view_mode"
        class="input-field w-24 sm:w-48 text-xs sm:text-sm"
        @change="
          (e) =>
            emit('update:settings', {
              ...settings,
              default_view_mode: (e.target as HTMLSelectElement).value,
            })
        "
      >
        <option value="original">{{ t('viewModeOriginal') }}</option>
        <option value="rendered">{{ t('viewModeRendered') }}</option>
      </select>
    </div>

    <div class="setting-item mt-2 sm:mt-3">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhImage :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('showArticlePreviewImages') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('showArticlePreviewImagesDesc') }}
          </div>
        </div>
      </div>
      <input
        :checked="settings.show_article_preview_images"
        type="checkbox"
        class="toggle"
        @change="
          (e) =>
            emit('update:settings', {
              ...settings,
              show_article_preview_images: (e.target as HTMLInputElement).checked,
            })
        "
      />
    </div>

    <div class="setting-item mt-2 sm:mt-3">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhArticleNyTimes :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('enableFullTextFetch') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('enableFullTextFetchDesc') }}
          </div>
        </div>
      </div>
      <input
        :checked="settings.full_text_fetch_enabled"
        type="checkbox"
        class="toggle"
        @change="
          (e) =>
            emit('update:settings', {
              ...settings,
              full_text_fetch_enabled: (e.target as HTMLInputElement).checked,
            })
        "
      />
    </div>

    <div class="setting-item mt-2 sm:mt-3">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhImages :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('imageGalleryEnabled') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('imageGalleryEnabledDesc') }}
          </div>
        </div>
      </div>
      <input
        :checked="settings.image_gallery_enabled"
        type="checkbox"
        class="toggle"
        @change="
          (e) =>
            emit('update:settings', {
              ...settings,
              image_gallery_enabled: (e.target as HTMLInputElement).checked,
            })
        "
      />
    </div>

    <div class="setting-item mt-2 sm:mt-3">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhCursorClick :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('hoverMarkAsRead') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('hoverMarkAsReadDesc') }}
          </div>
        </div>
      </div>
      <input
        :checked="settings.hover_mark_as_read"
        type="checkbox"
        class="toggle"
        @change="
          (e) =>
            emit('update:settings', {
              ...settings,
              hover_mark_as_read: (e.target as HTMLInputElement).checked,
            })
        "
      />
    </div>

    <div class="setting-item mt-2 sm:mt-3">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhEyeSlash :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('showHiddenArticles') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('showHiddenArticlesDesc') }}
          </div>
        </div>
      </div>
      <input
        :checked="settings.show_hidden_articles"
        type="checkbox"
        class="toggle"
        @change="
          (e) =>
            emit('update:settings', {
              ...settings,
              show_hidden_articles: (e.target as HTMLInputElement).checked,
            })
        "
      />
    </div>
  </div>
</template>

<style scoped>
@reference "../../../../style.css";

.input-field {
  @apply p-1.5 sm:p-2.5 border border-border rounded-md bg-bg-secondary text-text-primary focus:border-accent focus:outline-none transition-colors;
}
.toggle {
  @apply w-10 h-5 appearance-none bg-bg-tertiary rounded-full relative cursor-pointer border border-border transition-colors checked:bg-accent checked:border-accent shrink-0;
}
.toggle::after {
  content: '';
  @apply absolute top-0.5 left-0.5 w-3.5 h-3.5 bg-white rounded-full shadow-sm transition-transform;
}
.toggle:checked::after {
  transform: translateX(20px);
}
.setting-item {
  @apply flex items-center sm:items-start justify-between gap-2 sm:gap-4 p-2 sm:p-3 rounded-lg bg-bg-secondary border border-border;
}
</style>
