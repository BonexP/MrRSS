<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { PhRobot, PhKey, PhLink, PhBrain, PhPlus, PhSliders, PhTrash } from '@phosphor-icons/vue';
import type { SettingsData } from '@/types/settings';

const { t } = useI18n();

// Custom headers management
interface CustomHeader {
  id: string;
  name: string;
  value: string;
}

const customHeaders = ref<CustomHeader[]>([]);

// Parse custom headers from JSON string
function parseCustomHeaders(jsonString: string): CustomHeader[] {
  if (!jsonString || jsonString.trim() === '') {
    return [];
  }
  try {
    const headers = JSON.parse(jsonString);
    return Object.entries(headers).map(([name, value], index) => ({
      id: `${Date.now()}-${index}`,
      name,
      value: String(value),
    }));
  } catch (e) {
    console.error('Failed to parse custom headers:', e);
    return [];
  }
}

// Convert custom headers to JSON string
function stringifyCustomHeaders(headers: CustomHeader[]): string {
  const validHeaders = headers.filter((h) => h.name.trim() !== '');
  if (validHeaders.length === 0) {
    return '';
  }
  const headersObj: Record<string, string> = {};
  validHeaders.forEach((h) => {
    headersObj[h.name] = h.value;
  });
  return JSON.stringify(headersObj);
}

// Load custom headers from settings
function loadCustomHeaders() {
  const jsonString = props.settings.ai_custom_headers || '';
  customHeaders.value = parseCustomHeaders(jsonString);
}

// Add a new custom header
function addCustomHeader() {
  customHeaders.value.push({
    id: `${Date.now()}`,
    name: '',
    value: '',
  });
}

// Remove a custom header
function removeCustomHeader(id: string) {
  const index = customHeaders.value.findIndex((h) => h.id === id);
  if (index !== -1) {
    customHeaders.value.splice(index, 1);
    saveCustomHeaders();
  }
}

// Save custom headers to settings (debounced)
let saveTimeout: ReturnType<typeof setTimeout> | null = null;
function saveCustomHeaders() {
  if (saveTimeout) {
    clearTimeout(saveTimeout);
  }
  saveTimeout = setTimeout(() => {
    const jsonString = stringifyCustomHeaders(customHeaders.value);
    emit('update:settings', {
      ...props.settings,
      ai_custom_headers: jsonString,
    });
    saveTimeout = null;
  }, 500);
}

interface Props {
  settings: SettingsData;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  'update:settings': [settings: SettingsData];
}>();

onMounted(() => {
  loadCustomHeaders();
});

// Watch for external changes to ai_custom_headers setting
// This ensures that if the setting is updated from elsewhere, we reload it
watch(
  () => props.settings.ai_custom_headers,
  (newValue, oldValue) => {
    // Only reload if the value actually changed and we're not the ones who changed it
    if (newValue !== oldValue) {
      const parsed = parseCustomHeaders(newValue || '');
      // Check if the parsed headers are different from current state
      const currentJSON = stringifyCustomHeaders(customHeaders.value);
      if (currentJSON !== newValue) {
        customHeaders.value = parsed;
      }
    }
  },
  { immediate: false }
);
</script>

<template>
  <div class="setting-group">
    <label
      class="font-semibold mb-2 sm:mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2"
    >
      <PhRobot :size="14" class="sm:w-4 sm:h-4" />
      {{ t('aiSettings') }}
    </label>
    <div class="text-xs text-text-secondary mb-3 sm:mb-4">
      {{ t('aiSettingsDesc') }}
    </div>

    <!-- API Key -->
    <div class="setting-item mb-2 sm:mb-4">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhKey :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm">
            {{ t('aiApiKey') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('aiApiKeyDesc') }}
          </div>
        </div>
      </div>
      <input
        :value="props.settings.ai_api_key"
        type="password"
        :placeholder="t('aiApiKeyPlaceholder')"
        class="input-field w-32 sm:w-48 text-xs sm:text-sm"
        @input="
          (e) =>
            emit('update:settings', {
              ...props.settings,
              ai_api_key: (e.target as HTMLInputElement).value,
            })
        "
      />
    </div>

    <!-- Endpoint -->
    <div class="setting-item mb-2 sm:mb-4">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhLink :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm">
            {{ t('aiEndpoint') }} <span class="text-red-500">*</span>
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('aiEndpointDesc') }}
          </div>
        </div>
      </div>
      <input
        :value="props.settings.ai_endpoint"
        type="text"
        :placeholder="t('aiEndpointPlaceholder')"
        class="input-field w-32 sm:w-48 text-xs sm:text-sm"
        @input="
          (e) =>
            emit('update:settings', {
              ...props.settings,
              ai_endpoint: (e.target as HTMLInputElement).value,
            })
        "
      />
    </div>

    <!-- Model -->
    <div class="setting-item mb-2 sm:mb-4">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhBrain :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm">
            {{ t('aiModel') }} <span class="text-red-500">*</span>
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('aiModelDesc') }}
          </div>
        </div>
      </div>
      <input
        :value="props.settings.ai_model"
        type="text"
        :placeholder="t('aiModelPlaceholder')"
        class="input-field w-32 sm:w-48 text-xs sm:text-sm"
        @input="
          (e) =>
            emit('update:settings', {
              ...props.settings,
              ai_model: (e.target as HTMLInputElement).value,
            })
        "
      />
    </div>

    <!-- Custom Headers -->
    <div class="setting-item mb-2 sm:mb-4 flex-col items-stretch w-full">
      <div class="flex items-center gap-2 sm:gap-3">
        <PhSliders :size="20" class="text-text-secondary shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium text-sm">{{ t('aiCustomHeaders') }}</div>
          <div class="text-xs text-text-secondary">{{ t('aiCustomHeadersDesc') }}</div>
        </div>
      </div>

      <!-- Headers List -->
      <div class="mt-2 sm:mt-3 space-y-1.5 sm:space-y-2 w-full">
        <div
          v-for="header in customHeaders"
          :key="header.id"
          class="flex items-center gap-1.5 sm:gap-2"
        >
          <input
            v-model="header.name"
            type="text"
            :placeholder="t('aiCustomHeadersName')"
            class="input-field text-xs sm:text-sm flex-1"
            @input="saveCustomHeaders()"
          />
          <input
            v-model="header.value"
            type="text"
            :placeholder="t('aiCustomHeadersValue')"
            class="input-field text-xs sm:text-sm flex-1"
            @input="saveCustomHeaders()"
          />
          <button
            type="button"
            class="p-1.5 sm:p-2 rounded hover:bg-red-50 dark:hover:bg-red-900/20 text-text-secondary hover:text-red-500 transition-all shrink-0"
            :title="t('aiCustomHeadersRemove')"
            @click="removeCustomHeader(header.id)"
          >
            <PhTrash :size="14" class="sm:w-4 sm:h-4" />
          </button>
        </div>

        <!-- Add Header Button -->
        <button
          type="button"
          class="w-full p-1.5 sm:p-2 rounded border border-dashed border-border text-text-secondary hover:border-accent hover:text-accent hover:bg-accent/5 transition-all text-xs font-medium flex items-center justify-center gap-1.5 sm:gap-2"
          @click="addCustomHeader"
        >
          <PhPlus :size="14" class="sm:w-4 sm:h-4" />
          <span>{{ t('aiCustomHeadersAdd') }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
@reference "../../../../style.css";

.input-field {
  @apply p-1.5 sm:p-2.5 border border-border rounded-md bg-bg-secondary text-text-primary focus:border-accent focus:outline-none transition-colors;
}

.setting-item {
  @apply flex items-center sm:items-start justify-between gap-2 sm:gap-4 p-2 sm:p-3 rounded-lg bg-bg-secondary border border-border;
}

.setting-group {
  @apply mb-4 sm:mb-6;
}
</style>
