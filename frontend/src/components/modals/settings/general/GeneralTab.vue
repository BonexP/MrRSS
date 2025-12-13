<script setup lang="ts">
import { toRef } from 'vue';
import type { SettingsData } from '@/types/settings';
import { useSettingsAutoSave } from '@/composables/core/useSettingsAutoSave';
import { useSettingsValidation } from '@/composables/core/useSettingsValidation';
import { useI18n } from 'vue-i18n';
import { PhWarning } from '@phosphor-icons/vue';
import AppearanceSettings from './AppearanceSettings.vue';
import UpdateSettings from './UpdateSettings.vue';
import ProxySettings from './ProxySettings.vue';
import DatabaseSettings from './DatabaseSettings.vue';
import TranslationSettings from './TranslationSettings.vue';
import SummarySettings from './SummarySettings.vue';

interface Props {
  settings: SettingsData;
}

const props = defineProps<Props>();
const { t } = useI18n();

const emit = defineEmits<{
  'update:settings': [settings: SettingsData];
}>();

// Use composable for auto-save with reactivity
const settingsRef = toRef(props, 'settings');
useSettingsAutoSave(settingsRef);

// Use validation composable
const { isValid, isTranslationValid, isSummaryValid, isProxyValid } =
  useSettingsValidation(settingsRef);
</script>

<template>
  <div class="space-y-4 sm:space-y-6">
    <!-- Validation Warning -->
    <div
      v-if="!isValid"
      class="p-3 sm:p-4 rounded-lg border-2 border-red-500 bg-red-500/10 flex items-start gap-3"
    >
      <PhWarning :size="20" class="text-red-500 shrink-0 mt-0.5" :weight="'fill'" />
      <div class="flex-1">
        <div class="font-semibold text-red-500 text-sm sm:text-base mb-1">
          {{ t('requiredField') }}
        </div>
        <div class="text-xs sm:text-sm text-text-secondary">
          <span v-if="!isTranslationValid">
            {{ t('translationCredentialsRequired') }}
          </span>
          <span v-if="!isTranslationValid && !isSummaryValid"> • </span>
          <span v-if="!isSummaryValid">
            {{ t('summaryCredentialsRequired') }}
          </span>
          <span v-if="(!isTranslationValid || !isSummaryValid) && !isProxyValid"> • </span>
          <span v-if="!isProxyValid">
            {{ t('proxyCredentialsRequired') }}
          </span>
        </div>
      </div>
    </div>

    <AppearanceSettings :settings="settings" @update:settings="emit('update:settings', $event)" />

    <UpdateSettings :settings="settings" @update:settings="emit('update:settings', $event)" />

    <ProxySettings :settings="settings" @update:settings="emit('update:settings', $event)" />

    <DatabaseSettings :settings="settings" @update:settings="emit('update:settings', $event)" />

    <TranslationSettings :settings="settings" @update:settings="emit('update:settings', $event)" />

    <SummarySettings :settings="settings" @update:settings="emit('update:settings', $event)" />
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
.info-display {
  @apply px-2 sm:px-3 py-1.5 sm:py-2 rounded-lg border border-border;
  background-color: rgba(233, 236, 239, 0.3);
}
.dark-mode .info-display {
  background-color: rgba(45, 45, 45, 0.3);
}
</style>
