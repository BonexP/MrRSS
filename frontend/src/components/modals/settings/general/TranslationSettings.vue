<script setup lang="ts">
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { PhGlobe, PhArticle, PhPackage, PhKey, PhLink, PhRobot } from '@phosphor-icons/vue';
import type { SettingsData } from '@/types/settings';

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
</script>

<template>
  <div class="setting-group">
    <label
      class="font-semibold mb-2 sm:mb-3 text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2"
    >
      <PhGlobe :size="14" class="sm:w-4 sm:h-4" />
      {{ t('translation') }}
    </label>
    <div class="setting-item mb-2 sm:mb-4">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhArticle :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
            {{ t('enableTranslation') }}
          </div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('enableTranslationDesc') }}
          </div>
        </div>
      </div>
      <input v-model="localSettings.translation_enabled" type="checkbox" class="toggle" />
    </div>

    <div
      v-if="settings.translation_enabled"
      class="ml-2 sm:ml-4 space-y-2 sm:space-y-3 border-l-2 border-border pl-2 sm:pl-4"
    >
      <div class="sub-setting-item">
        <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
          <PhPackage :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
          <div class="flex-1 min-w-0">
            <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('translationProvider') }}</div>
            <div class="text-xs text-text-secondary hidden sm:block">
              {{ t('translationProviderDesc') || 'Choose the translation service to use' }}
            </div>
          </div>
        </div>
        <select
          v-model="localSettings.translation_provider"
          class="input-field w-32 sm:w-48 text-xs sm:text-sm"
        >
          <option value="google">{{ t('googleTranslate') }}</option>
          <option value="deepl">{{ t('deeplApi') }}</option>
          <option value="baidu">{{ t('baiduTranslate') }}</option>
          <option value="ai">{{ t('aiTranslation') }}</option>
        </select>
      </div>

      <!-- Google Translate Endpoint -->
      <div v-if="settings.translation_provider === 'google'" class="sub-setting-item">
        <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
          <PhLink :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
          <div class="flex-1 min-w-0">
            <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('googleTranslateEndpoint') }}</div>
            <div class="text-xs text-text-secondary hidden sm:block">
              {{ t('googleTranslateEndpointDesc') }}
            </div>
          </div>
        </div>
        <select
          v-model="localSettings.google_translate_endpoint"
          class="input-field w-32 sm:w-48 text-xs sm:text-sm"
        >
          <option value="translate.googleapis.com">
            {{ t('googleTranslateEndpointDefault') }}
          </option>
          <option value="clients5.google.com">{{ t('googleTranslateEndpointAlternate') }}</option>
        </select>
      </div>

      <!-- DeepL API Key -->
      <div v-if="settings.translation_provider === 'deepl'" class="sub-setting-item">
        <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
          <PhKey :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
          <div class="flex-1 min-w-0">
            <div class="font-medium mb-0 sm:mb-1 text-sm">
              {{ t('deeplApiKey') }} <span class="text-red-500">*</span>
            </div>
            <div class="text-xs text-text-secondary hidden sm:block">
              {{ t('deeplApiKeyDesc') || 'Enter your DeepL API key' }}
            </div>
          </div>
        </div>
        <input
          v-model="localSettings.deepl_api_key"
          type="password"
          :placeholder="t('deeplApiKeyPlaceholder')"
          :class="[
            'input-field w-32 sm:w-48 text-xs sm:text-sm',
            settings.translation_enabled &&
            settings.translation_provider === 'deepl' &&
            !settings.deepl_api_key?.trim()
              ? 'border-red-500'
              : '',
          ]"
        />
      </div>

      <!-- Baidu Translate Settings -->
      <template v-if="settings.translation_provider === 'baidu'">
        <div class="sub-setting-item">
          <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
            <PhKey :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
            <div class="flex-1 min-w-0">
              <div class="font-medium mb-0 sm:mb-1 text-sm">
                {{ t('baiduAppId') }} <span class="text-red-500">*</span>
              </div>
              <div class="text-xs text-text-secondary hidden sm:block">
                {{ t('baiduAppIdDesc') }}
              </div>
            </div>
          </div>
          <input
            v-model="localSettings.baidu_app_id"
            type="text"
            :placeholder="t('baiduAppIdPlaceholder')"
            :class="[
              'input-field w-32 sm:w-48 text-xs sm:text-sm',
              settings.translation_enabled &&
              settings.translation_provider === 'baidu' &&
              !settings.baidu_app_id?.trim()
                ? 'border-red-500'
                : '',
            ]"
          />
        </div>
        <div class="sub-setting-item">
          <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
            <PhKey :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
            <div class="flex-1 min-w-0">
              <div class="font-medium mb-0 sm:mb-1 text-sm">
                {{ t('baiduSecretKey') }} <span class="text-red-500">*</span>
              </div>
              <div class="text-xs text-text-secondary hidden sm:block">
                {{ t('baiduSecretKeyDesc') }}
              </div>
            </div>
          </div>
          <input
            v-model="localSettings.baidu_secret_key"
            type="password"
            :placeholder="t('baiduSecretKeyPlaceholder')"
            :class="[
              'input-field w-32 sm:w-48 text-xs sm:text-sm',
              settings.translation_enabled &&
              settings.translation_provider === 'baidu' &&
              !settings.baidu_secret_key?.trim()
                ? 'border-red-500'
                : '',
            ]"
          />
        </div>
      </template>

      <!-- AI Translation Settings -->
      <template v-if="settings.translation_provider === 'ai'">
        <div class="sub-setting-item">
          <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
            <PhKey :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
            <div class="flex-1 min-w-0">
              <div class="font-medium mb-0 sm:mb-1 text-sm">
                {{ t('aiApiKey') }} <span class="text-red-500">*</span>
              </div>
              <div class="text-xs text-text-secondary hidden sm:block">
                {{ t('aiApiKeyDesc') }}
              </div>
            </div>
          </div>
          <input
            v-model="localSettings.ai_api_key"
            type="password"
            :placeholder="t('aiApiKeyPlaceholder')"
            :class="[
              'input-field w-32 sm:w-48 text-xs sm:text-sm',
              settings.translation_enabled &&
              settings.translation_provider === 'ai' &&
              !settings.ai_api_key?.trim()
                ? 'border-red-500'
                : '',
            ]"
          />
        </div>
        <div class="sub-setting-item">
          <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
            <PhLink :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
            <div class="flex-1 min-w-0">
              <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('aiEndpoint') }}</div>
              <div class="text-xs text-text-secondary hidden sm:block">
                {{ t('aiEndpointDesc') }}
              </div>
            </div>
          </div>
          <input
            v-model="localSettings.ai_endpoint"
            type="text"
            :placeholder="t('aiEndpointPlaceholder')"
            class="input-field w-32 sm:w-48 text-xs sm:text-sm"
          />
        </div>
        <div class="sub-setting-item">
          <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
            <PhRobot :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
            <div class="flex-1 min-w-0">
              <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('aiModel') }}</div>
              <div class="text-xs text-text-secondary hidden sm:block">
                {{ t('aiModelDesc') }}
              </div>
            </div>
          </div>
          <input
            v-model="localSettings.ai_model"
            type="text"
            :placeholder="t('aiModelPlaceholder')"
            class="input-field w-32 sm:w-48 text-xs sm:text-sm"
          />
        </div>
        <div class="sub-setting-item flex-col items-stretch gap-2">
          <div class="flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
            <PhRobot :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
            <div class="flex-1 min-w-0">
              <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('aiSystemPrompt') }}</div>
              <div class="text-xs text-text-secondary hidden sm:block">
                {{ t('aiSystemPromptDesc') }}
              </div>
            </div>
          </div>
          <textarea
            v-model="localSettings.ai_system_prompt"
            class="input-field w-full text-xs sm:text-sm resize-none"
            rows="3"
          />
        </div>
      </template>

      <div class="sub-setting-item">
        <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
          <PhGlobe :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
          <div class="flex-1 min-w-0">
            <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('targetLanguage') }}</div>
            <div class="text-xs text-text-secondary hidden sm:block">
              {{ t('targetLanguageDesc') || 'Language to translate article titles to' }}
            </div>
          </div>
        </div>
        <select
          v-model="localSettings.target_language"
          class="input-field w-24 sm:w-48 text-xs sm:text-sm"
        >
          <option value="en">{{ t('english') }}</option>
          <option value="es">{{ t('spanish') }}</option>
          <option value="fr">{{ t('french') }}</option>
          <option value="de">{{ t('german') }}</option>
          <option value="zh">{{ t('chinese') }}</option>
          <option value="ja">{{ t('japanese') }}</option>
        </select>
      </div>
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
