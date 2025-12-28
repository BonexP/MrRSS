<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { PhArchive, PhFolders } from '@phosphor-icons/vue';
import type { SettingsData } from '@/types/settings';

const { t } = useI18n();

interface Props {
  settings: SettingsData;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  'update:settings': [settings: SettingsData];
}>();

// Handler for checkbox change
function handleObsidianEnabledChange(event: Event) {
  const target = event.target as HTMLInputElement;
  emit('update:settings', {
    ...props.settings,
    obsidian_enabled: target.checked,
  });
}

// Handler for vault input change
function handleObsidianVaultChange(event: Event) {
  const target = event.target as HTMLInputElement;
  emit('update:settings', {
    ...props.settings,
    obsidian_vault: target.value,
  });
}

// Handler for vault path input change
function handleObsidianVaultPathChange(event: Event) {
  const target = event.target as HTMLInputElement;
  emit('update:settings', {
    ...props.settings,
    obsidian_vault_path: target.value,
  });
}
</script>

<template>
  <!-- Enable Obsidian Integration -->
  <div class="setting-item">
    <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
      <img
        src="/assets/plugin_icons/obsidian.svg"
        alt="Obsidian"
        class="w-5 h-5 sm:w-6 sm:h-6 mt-0.5 shrink-0"
      />
      <div class="flex-1 min-w-0">
        <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
          {{ t('obsidianIntegration') }}
        </div>
        <div class="text-xs text-text-secondary hidden sm:block">
          {{ t('obsidianIntegrationDescription') }}
        </div>
      </div>
    </div>
    <input
      type="checkbox"
      :checked="props.settings.obsidian_enabled"
      class="toggle"
      @change="handleObsidianEnabledChange"
    />
  </div>

  <div
    v-if="props.settings.obsidian_enabled"
    class="ml-2 sm:ml-4 space-y-2 sm:space-y-3 border-l-2 border-border pl-2 sm:pl-4"
  >
    <!-- Vault Name -->
    <div class="sub-setting-item">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhArchive :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('obsidianVaultName') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('obsidianVaultNameDesc') }}
          </div>
        </div>
      </div>
      <input
        type="text"
        :value="props.settings.obsidian_vault"
        :placeholder="t('obsidianVaultNamePlaceholder')"
        class="input-field w-32 sm:w-48 text-xs sm:text-sm"
        @input="handleObsidianVaultChange"
      />
    </div>

    <!-- Vault Path -->
    <div class="sub-setting-item">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhFolders :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('obsidianVaultPath') }} <span class="text-red-500">*</span>
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('obsidianVaultPathDesc') }}
          </div>
        </div>
      </div>
      <input
        type="text"
        :value="props.settings.obsidian_vault_path"
        :placeholder="t('obsidianVaultPathPlaceholder')"
        class="input-field w-48 sm:w-64 text-xs sm:text-sm"
        @input="handleObsidianVaultPathChange"
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
.sub-setting-item {
  @apply flex items-center sm:items-start justify-between gap-2 sm:gap-4 p-2 sm:p-2.5 rounded-md bg-bg-tertiary;
}
.setting-group {
  @apply space-y-2 sm:space-y-3;
}
</style>
