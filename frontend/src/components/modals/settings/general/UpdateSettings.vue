<script setup lang="ts">
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  PhArchiveTray,
  PhArrowClockwise,
  PhArrowsClockwise,
  PhClock,
  PhPower,
} from '@phosphor-icons/vue';
import type { SettingsData } from '@/types/settings';
import { formatRelativeTime } from '@/utils/date';

const { t } = useI18n();

interface Props {
  settings: SettingsData;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  'update:settings': [settings: SettingsData];
}>();

// Create local reactive copy
const localSettings = ref<SettingsData>({ ...props.settings });

// Watch for changes and emit updates
watch(
  localSettings,
  (newSettings) => {
    emit('update:settings', { ...newSettings });
  },
  { deep: true }
);

// Format last update time using shared utility
function formatLastUpdate(timestamp: string): string {
  return formatRelativeTime(timestamp, props.settings.language, t);
}
</script>

<template>
  <div class="setting-group">
    <label
      class="font-semibold mb-2 sm:mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2"
    >
      <PhArrowClockwise :size="14" class="sm:w-4 sm:h-4" />
      {{ t('updates') }}
    </label>
    <div class="setting-item">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhArrowsClockwise :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('refreshMode') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('refreshModeDesc') }}
          </div>
          <!-- Last update time - now as sub-item -->
          <div class="text-xs text-text-secondary mt-1 flex items-center gap-1">
            <span>{{ t('lastArticleUpdate') }}:</span>
            <span class="font-medium text-accent">{{
              formatLastUpdate(settings.last_article_update)
            }}</span>
          </div>
        </div>
      </div>
      <select
        v-model="localSettings.refresh_mode"
        class="input-field w-32 sm:w-40 text-xs sm:text-sm"
      >
        <option value="fixed">{{ t('fixedInterval') }}</option>
        <option value="intelligent">{{ t('intelligentInterval') }}</option>
      </select>
    </div>

    <!-- Auto Update Interval (shown when fixed mode is selected) -->
    <div
      v-if="settings.refresh_mode === 'fixed'"
      class="mt-2 sm:mt-3 ml-4 sm:ml-6 space-y-2 sm:space-y-3 border-l-2 border-border pl-3 sm:pl-4"
    >
      <div class="sub-setting-item">
        <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
          <PhClock :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
          <div class="flex-1 min-w-0">
            <div class="font-medium mb-0 sm:mb-1 text-xs sm:text-sm">
              {{ t('autoUpdateInterval') }}
            </div>
            <div class="text-[10px] sm:text-xs text-text-secondary hidden sm:block">
              {{ t('autoUpdateIntervalDesc') }}
            </div>
          </div>
        </div>
        <input
          v-model="localSettings.update_interval"
          type="number"
          min="1"
          class="input-field w-16 sm:w-20 text-center text-xs sm:text-sm"
        />
      </div>
    </div>

    <div class="setting-item mt-2 sm:mt-3">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhPower :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('startupOnBoot') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('startupOnBootDesc') }}
          </div>
        </div>
      </div>
      <input v-model="localSettings.startup_on_boot" type="checkbox" class="toggle" />
    </div>

    <div class="setting-item mt-2 sm:mt-3">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhArchiveTray :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('closeToTray') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('closeToTrayDesc') }}
          </div>
        </div>
      </div>
      <input v-model="localSettings.close_to_tray" type="checkbox" class="toggle" />
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
.sub-setting-item {
  @apply flex items-center sm:items-start justify-between gap-2 sm:gap-4 p-2 sm:p-2.5 rounded-md bg-bg-tertiary;
}
</style>
