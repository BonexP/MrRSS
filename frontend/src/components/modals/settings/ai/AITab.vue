<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { PhInfo } from '@phosphor-icons/vue';
import { computed } from 'vue';
import type { SettingsData } from '@/types/settings';
import { useSettingsAutoSave } from '@/composables/core/useSettingsAutoSave';
import AISettings from './AISettings.vue';
import AITestSettings from './AITestSettings.vue';
import AIUsageSettings from './AIUsageSettings.vue';
import AIFeatureSettings from './AIFeatureSettings.vue';

const { t } = useI18n();

interface Props {
  settings: SettingsData;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  'update:settings': [settings: SettingsData];
}>();

// Create a computed ref that returns the settings object
const settingsRef = computed(() => props.settings);

// Use composable for auto-save
useSettingsAutoSave(settingsRef);

// Handler for settings updates from child components
function handleUpdateSettings(updatedSettings: SettingsData) {
  emit('update:settings', updatedSettings);
}
</script>

<template>
  <div class="space-y-4 sm:space-y-6">
    <div class="tip-box">
      <PhInfo :size="16" class="text-accent shrink-0 sm:w-5 sm:h-5" />
      <span class="text-xs sm:text-sm">{{ t('aiIsDanger') }}</span>
    </div>
    <AISettings :settings="settings" @update:settings="handleUpdateSettings" />
    <AITestSettings :settings="settings" @update:settings="handleUpdateSettings" />
    <AIUsageSettings :settings="settings" @update:settings="handleUpdateSettings" />
    <AIFeatureSettings :settings="settings" @update:settings="handleUpdateSettings" />
  </div>
</template>

<style scoped>
.tip-box {
  @apply flex items-center gap-2 sm:gap-3 py-2 sm:py-2.5 px-2.5 sm:px-3 rounded-lg;
  background-color: rgba(59, 130, 246, 0.05);
  border: 1px solid rgba(59, 130, 246, 0.3);
}
</style>
