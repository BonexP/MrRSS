<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { PhTextAlignLeft, PhTextT, PhPackage, PhKey, PhLink, PhRobot } from '@phosphor-icons/vue';
import type { SettingsData } from '@/types/settings';

const { t } = useI18n();

interface Props {
  settings: SettingsData;
}

defineProps<Props>();
</script>

<template>
  <div class="setting-group">
    <label
      class="font-semibold mb-2 sm:mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2"
    >
      <PhTextAlignLeft :size="14" class="sm:w-4 sm:h-4" />
      {{ t('summary') }}
    </label>
    <div class="setting-item mb-2 sm:mb-4">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhTextT :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('enableSummary') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('enableSummaryDesc') }}
          </div>
        </div>
      </div>
      <input type="checkbox" v-model="settings.summary_enabled" class="toggle" />
    </div>

    <div
      v-if="settings.summary_enabled"
      class="ml-2 sm:ml-4 space-y-2 sm:space-y-3 border-l-2 border-border pl-2 sm:pl-4"
    >
      <div class="sub-setting-item">
        <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
          <PhPackage :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
          <div class="flex-1 min-w-0">
            <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('summaryProvider') }}</div>
            <div class="text-xs text-text-secondary hidden sm:block">
              {{ t('summaryProviderDesc') }}
            </div>
          </div>
        </div>
        <select
          v-model="settings.summary_provider"
          class="input-field w-32 sm:w-48 text-xs sm:text-sm"
        >
          <option value="local">{{ t('localAlgorithm') }}</option>
          <option value="ai">{{ t('aiSummary') }}</option>
        </select>
      </div>

      <!-- AI Summary Settings -->
      <template v-if="settings.summary_provider === 'ai'">
        <div class="sub-setting-item">
          <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
            <PhKey :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
            <div class="flex-1 min-w-0">
              <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('summaryAiApiKey') }}</div>
              <div class="text-xs text-text-secondary hidden sm:block">
                {{ t('summaryAiApiKeyDesc') }}
              </div>
            </div>
          </div>
          <input
            type="password"
            v-model="settings.summary_ai_api_key"
            :placeholder="t('summaryAiApiKeyPlaceholder')"
            class="input-field w-32 sm:w-48 text-xs sm:text-sm"
          />
        </div>
        <div class="sub-setting-item">
          <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
            <PhLink :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
            <div class="flex-1 min-w-0">
              <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('summaryAiEndpoint') }}</div>
              <div class="text-xs text-text-secondary hidden sm:block">
                {{ t('summaryAiEndpointDesc') }}
              </div>
            </div>
          </div>
          <input
            type="text"
            v-model="settings.summary_ai_endpoint"
            :placeholder="t('summaryAiEndpointPlaceholder')"
            class="input-field w-32 sm:w-48 text-xs sm:text-sm"
          />
        </div>
        <div class="sub-setting-item">
          <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
            <PhRobot :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
            <div class="flex-1 min-w-0">
              <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('summaryAiModel') }}</div>
              <div class="text-xs text-text-secondary hidden sm:block">
                {{ t('summaryAiModelDesc') }}
              </div>
            </div>
          </div>
          <input
            type="text"
            v-model="settings.summary_ai_model"
            :placeholder="t('summaryAiModelPlaceholder')"
            class="input-field w-32 sm:w-48 text-xs sm:text-sm"
          />
        </div>
        <div class="sub-setting-item flex-col items-stretch gap-2">
          <div class="flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
            <PhRobot :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
            <div class="flex-1 min-w-0">
              <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('summaryAiSystemPrompt') }}</div>
              <div class="text-xs text-text-secondary hidden sm:block">
                {{ t('summaryAiSystemPromptDesc') }}
              </div>
            </div>
          </div>
          <textarea
            v-model="settings.summary_ai_system_prompt"
            :placeholder="t('summaryAiSystemPromptPlaceholder')"
            class="input-field text-xs sm:text-sm resize-y"
            rows="3"
          />
        </div>
      </template>

      <div class="sub-setting-item">
        <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
          <PhTextAlignLeft :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
          <div class="flex-1 min-w-0">
            <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('summaryLength') }}</div>
            <div class="text-xs text-text-secondary hidden sm:block">
              {{ t('summaryLengthDesc') }}
            </div>
          </div>
        </div>
        <select
          v-model="settings.summary_length"
          class="input-field w-24 sm:w-48 text-xs sm:text-sm"
        >
          <option value="short">{{ t('summaryLengthShort') }}</option>
          <option value="medium">{{ t('summaryLengthMedium') }}</option>
          <option value="long">{{ t('summaryLengthLong') }}</option>
        </select>
      </div>
    </div>
  </div>
</template>

<style scoped>
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
