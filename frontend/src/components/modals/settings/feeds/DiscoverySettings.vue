<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { PhPlay, PhBinoculars, PhInfo } from '@phosphor-icons/vue';

const { t } = useI18n();

const emit = defineEmits<{
  'discover-all': [];
}>();

function handleDiscoverAll() {
  emit('discover-all');
}
</script>

<template>
  <div class="setting-group">
    <label
      class="font-semibold mb-2 sm:mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2"
    >
      <PhBinoculars :size="14" class="sm:w-4 sm:h-4" />
      {{ t('feedDiscovery') }}
    </label>

    <div class="tip-box">
      <PhInfo :size="16" class="text-accent shrink-0 sm:w-5 sm:h-5" />
      <span class="text-xs sm:text-sm">{{ t('discoveryLongRunningWarning') }}</span>
    </div>

    <!-- Discover All Feeds -->
    <div class="setting-item mt-2 sm:mt-3">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhBinoculars :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('discoverAllFeeds') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('discoverAllFeedsDesc') }}
          </div>
        </div>
      </div>
      <button class="btn-secondary" @click="handleDiscoverAll">
        <PhPlay :size="16" class="sm:w-5 sm:h-5" />
        <span class="hidden sm:inline">{{ t('startDiscovery') }}</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
@reference "../../../../style.css";

.setting-item {
  @apply flex items-center sm:items-start justify-between gap-2 sm:gap-4 p-2 sm:p-3 rounded-lg bg-bg-secondary border border-border;
}

.btn-secondary {
  @apply bg-bg-tertiary border border-border text-text-primary px-3 sm:px-4 py-1.5 sm:py-2 rounded-md cursor-pointer flex items-center gap-1.5 sm:gap-2 font-medium hover:bg-bg-secondary transition-colors;
}

.tip-box {
  @apply flex items-center gap-2 sm:gap-3 py-2 sm:py-2.5 px-2.5 sm:px-3 rounded-lg;
  background-color: rgba(59, 130, 246, 0.05);
  border: 1px solid rgba(59, 130, 246, 0.3);
}
</style>
