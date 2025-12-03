<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { ref, type Ref } from 'vue';
import {
  PhHardDrives,
  PhUpload,
  PhDownload,
  PhBroom,
  PhMagnifyingGlass,
} from '@phosphor-icons/vue';

const { t } = useI18n();

const emit = defineEmits<{
  'import-opml': [event: Event];
  'export-opml': [];
  'cleanup-database': [];
  'discover-all': [];
}>();

const opmlInput: Ref<HTMLInputElement | null> = ref(null);

function clickFileInput() {
  opmlInput.value?.click();
}

function handleImportOPML(event: Event) {
  emit('import-opml', event);
}

function handleExportOPML() {
  emit('export-opml');
}

function handleCleanupDatabase() {
  emit('cleanup-database');
}

function handleDiscoverAll() {
  emit('discover-all');
}
</script>

<template>
  <div class="setting-group">
    <label
      class="font-semibold mb-2 sm:mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2"
    >
      <PhHardDrives :size="14" class="sm:w-4 sm:h-4" />
      {{ t('dataManagement') }}
    </label>
    <div class="flex flex-col sm:flex-row gap-2 sm:gap-3 mb-2 sm:mb-3">
      <button
        @click="clickFileInput"
        class="btn-secondary flex-1 justify-center text-sm sm:text-base"
      >
        <PhUpload :size="18" class="sm:w-5 sm:h-5" /> {{ t('importOPML') }}
      </button>
      <input type="file" ref="opmlInput" class="hidden" @change="handleImportOPML" />
      <button
        @click="handleExportOPML"
        class="btn-secondary flex-1 justify-center text-sm sm:text-base"
      >
        <PhDownload :size="18" class="sm:w-5 sm:h-5" /> {{ t('exportOPML') }}
      </button>
    </div>
    <div class="flex mb-2 sm:mb-3">
      <button
        @click="handleCleanupDatabase"
        class="btn-danger flex-1 justify-center text-sm sm:text-base"
      >
        <PhBroom :size="18" class="sm:w-5 sm:h-5" /> {{ t('cleanDatabase') }}
      </button>
    </div>
    <div class="flex flex-col sm:flex-row gap-2 sm:gap-3 mb-2 sm:mb-3">
      <button
        @click="handleDiscoverAll"
        class="btn-primary flex-1 justify-center text-sm sm:text-base"
      >
        <PhMagnifyingGlass :size="18" class="sm:w-5 sm:h-5" />
        {{ t('discoverAllFeeds') }}
      </button>
    </div>
    <p class="text-xs text-text-secondary mb-2">
      {{ t('discoverAllFeedsDesc') }}
    </p>
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
  @apply bg-bg-tertiary border border-border text-text-primary px-3 sm:px-4 py-1.5 sm:py-2 rounded-md cursor-pointer flex items-center gap-1.5 sm:gap-2 font-medium hover:bg-bg-secondary transition-colors;
}
.btn-secondary:disabled {
  @apply opacity-50 cursor-not-allowed;
}
.btn-danger {
  @apply bg-red-50 border border-red-300 text-red-600 px-3 sm:px-4 py-1.5 sm:py-2 rounded-md cursor-pointer flex items-center gap-1.5 sm:gap-2 font-semibold hover:bg-red-100 transition-colors;
}
:global(.dark-mode) .btn-danger {
  @apply bg-red-900/20 border-red-400 text-red-400 hover:bg-red-900/30;
}
.btn-danger:disabled {
  @apply opacity-50 cursor-not-allowed;
}
</style>
