<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { PhChartLine, PhArrowCounterClockwise } from '@phosphor-icons/vue';
import type { SettingsData } from '@/types/settings';

const { t } = useI18n();

interface Props {
  settings: SettingsData;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  'update:settings': [settings: SettingsData];
}>();

// AI usage tracking
const aiUsage = ref<{
  usage: number;
  limit: number;
  limit_reached: boolean;
}>({
  usage: 0,
  limit: 0,
  limit_reached: false,
});

async function fetchAIUsage() {
  try {
    const response = await fetch('/api/ai-usage');
    if (response.ok) {
      aiUsage.value = await response.json();
    }
  } catch (e) {
    console.error('Failed to fetch AI usage:', e);
  }
}

async function resetAIUsage() {
  if (!window.confirm(t('aiUsageResetConfirm'))) {
    return;
  }
  try {
    const response = await fetch('/api/ai-usage/reset', { method: 'POST' });
    if (response.ok) {
      await fetchAIUsage();
      // Reset the local settings value as well
      emit('update:settings', {
        ...props.settings,
        ai_usage_tokens: '0',
      });
      window.showToast(t('aiUsageResetSuccess'), 'success');
    }
  } catch (e) {
    console.error('Failed to reset AI usage:', e);
    window.showToast(t('aiUsageResetError'), 'error');
  }
}

// Calculate usage percentage
function getUsagePercentage(): number {
  if (aiUsage.value.limit === 0) return 0;
  return Math.min(100, (aiUsage.value.usage / aiUsage.value.limit) * 100);
}

onMounted(() => {
  fetchAIUsage();
});
</script>

<template>
  <div class="setting-group">
    <!-- AI Usage Display -->
    <div class="setting-group mb-2 sm:mb-4">
      <div class="flex items-center justify-between mb-2 sm:mb-3">
        <label
          class="font-semibold text-text-secondary uppercase text-xs tracking-wider flex items-center gap-2"
        >
          <PhChartLine :size="14" class="sm:w-4 sm:h-4" />
          {{ t('aiUsage') }}
        </label>
      </div>

      <!-- Usage Status Display (Similar to Network Settings) -->
      <div
        class="flex flex-col sm:flex-row sm:items-stretch sm:justify-between gap-3 sm:gap-4 p-2 sm:p-3 rounded-lg bg-bg-secondary border border-border"
      >
        <!-- Tokens Used Box -->
        <div class="flex items-center">
          <div
            class="flex flex-col gap-2 p-3 rounded-lg bg-bg-primary border border-border w-full sm:min-w-[120px]"
          >
            <span class="text-sm text-text-secondary text-left">{{ t('aiUsageTokens') }}</span>
            <div class="flex items-baseline gap-1">
              <span class="text-xl sm:text-2xl font-bold text-text-primary"
                >{{ aiUsage.usage.toLocaleString() }} /
                {{ aiUsage.limit > 0 ? aiUsage.limit.toLocaleString() : '∞' }}</span
              >
              <span class="text-sm text-text-secondary">{{ t('tokens') }}</span>
            </div>
          </div>
        </div>

        <div class="flex flex-col sm:justify-between flex-1 gap-2 sm:gap-0">
          <div class="flex justify-center sm:justify-end">
            <button type="button" class="btn-secondary" @click="resetAIUsage">
              <PhArrowCounterClockwise :size="16" />
              {{ t('aiUsageReset') }}
            </button>
          </div>

          <!-- Progress bar (only shown if limit is set) -->
          <div
            v-if="aiUsage.limit > 0"
            class="flex flex-row items-center justify-center sm:justify-end gap-2"
          >
            <div class="flex items-center justify-between text-xs text-text-secondary">
              <span>{{ t('progress') }}</span>
              <span class="text-accent">{{ getUsagePercentage().toFixed(2) }}%</span>
            </div>
            <div class="relative h-2 bg-bg-tertiary rounded-full overflow-hidden">
              <div
                class="absolute top-0 left-0 h-full transition-all duration-300 rounded-full"
                :class="aiUsage.limit_reached ? 'bg-red-500' : 'bg-accent'"
                :style="{ width: getUsagePercentage() + '%' }"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Set AI Usage Limit -->
    <div class="setting-item mb-2 sm:mb-4">
      <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
        <PhChartLine :size="20" class="text-text-secondary mt-0.5 shrink-0 sm:w-6 sm:h-6" />
        <div class="flex-1 min-w-0">
          <div class="font-medium mb-0 sm:mb-1 text-sm">{{ t('setUsageLimit') }}</div>
          <div class="text-xs text-text-secondary hidden sm:block">
            {{ t('setUsageLimitDesc') }}
          </div>
        </div>
      </div>
      <input
        :value="props.settings.ai_usage_limit"
        type="number"
        min="0"
        :placeholder="t('aiUsageLimitPlaceholder')"
        class="input-field w-32 sm:w-48 text-xs sm:text-sm"
        @input="
          (e) =>
            emit('update:settings', {
              ...props.settings,
              ai_usage_limit: (e.target as HTMLInputElement).value,
            })
        "
      />
    </div>
  </div>
</template>

<style scoped>
@reference "../../../../style.css";

.btn-secondary {
  @apply bg-bg-tertiary border border-border text-text-primary px-3 sm:px-4 py-1.5 sm:py-2 rounded-md cursor-pointer flex items-center gap-1.5 sm:gap-2 font-medium hover:bg-bg-secondary transition-colors;
}

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
