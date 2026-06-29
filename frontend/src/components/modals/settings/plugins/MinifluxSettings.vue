<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { PhLink, PhKey, PhTestTube, PhArrowClockwise, PhCloudCheck } from '@phosphor-icons/vue';
import type { SettingsData } from '@/types/settings';
import { useAppStore } from '@/stores/app';
import { NestedSettingsContainer, SubSettingItem, InputControl } from '@/components/settings';

const { t } = useI18n();
const appStore = useAppStore();

interface Props {
  settings: SettingsData;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  'update:settings': [settings: SettingsData];
  'settings-changed': [];
}>();

const isSyncing = ref(false);
const syncStatus = ref<{
  last_sync_time: string | null;
}>({
  last_sync_time: null,
});

let statusPollInterval: ReturnType<typeof setInterval> | null = null;

async function fetchSyncStatus() {
  try {
    const response = await fetch('/api/miniflux/status');
    if (response.ok) {
      const data = await response.json();
      syncStatus.value = data;
    } else {
      console.error('[Miniflux Settings] fetchSyncStatus failed:', response.status, response.statusText);
    }
  } catch (e) {
    console.error('[Miniflux Settings] fetchSyncStatus error:', e);
  }
}

function startStatusPolling() {
  fetchSyncStatus();
  statusPollInterval = setInterval(fetchSyncStatus, 5000);
}

function stopStatusPolling() {
  if (statusPollInterval) {
    clearInterval(statusPollInterval);
    statusPollInterval = null;
  }
}

onMounted(() => {
  if (props.settings.miniflux_enabled) {
    startStatusPolling();
  }
});

onUnmounted(() => {
  stopStatusPolling();
});

function updateSetting(key: keyof SettingsData, value: any) {
  emit('update:settings', {
    ...props.settings,
    [key]: value,
  });
}

async function testConnection() {
  try {
    const response = await fetch('/api/miniflux/test-connection', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    });
    const data = await response.json();
    if (data.success) {
      console.log('[Miniflux Settings] testConnection: success');
      window.showToast(t('setting.miniflux.connectionSuccess'), 'success');
    } else {
      console.error('[Miniflux Settings] testConnection: failed -', data.error);
      window.showToast(data.error || t('setting.miniflux.connectionFailed'), 'error');
    }
  } catch (e) {
    console.error('[Miniflux Settings] testConnection error:', e);
    window.showToast(t('setting.miniflux.connectionFailed'), 'error');
  }
}

async function syncNow() {
  isSyncing.value = true;
  try {
    const response = await fetch('/api/miniflux/sync', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    });
    if (response.ok) {
      console.log('[Miniflux Settings] syncNow: sync started');
      window.showToast(t('setting.miniflux.syncStarted'), 'success');
    } else {
      const errMsg = `sync failed (${response.status})`;
      console.error('[Miniflux Settings] syncNow:', errMsg);
      throw new Error(t('setting.miniflux.syncFailed'));
    }
  } catch (error) {
    console.error('[Miniflux Settings] syncNow error:', error);
    window.showToast(
      error instanceof Error ? error.message : t('setting.miniflux.syncFailed'),
      'error'
    );
  } finally {
    isSyncing.value = false;
  }
}

async function handleMinifluxToggle(event: Event) {
  const target = event.target as HTMLInputElement;
  const newEnabled = target.checked;

  if (!newEnabled && props.settings.miniflux_enabled) {
    const confirmed = await window.showConfirm({
      title: t('setting.miniflux.enabled'),
      message: t('setting.miniflux.disableConfirm'),
      isDanger: true,
    });
    if (!confirmed) {
      target.checked = true;
      return;
    }
  }

  console.log('[Miniflux Settings] handleMinifluxToggle:', newEnabled ? 'enabled' : 'disabled');
  updateSetting('miniflux_enabled', newEnabled);
}

watch(
  () => props.settings.miniflux_enabled,
  async (newVal: boolean, oldVal: boolean) => {
    if (oldVal && !newVal) {
      console.log('[Miniflux Settings] Miniflux disabled, stopping polling');
      appStore.stopMinifluxStatusPolling();
      stopStatusPolling();
      setTimeout(async () => {
        await appStore.fetchFeeds();
        await appStore.fetchArticles();
        await appStore.fetchUnreadCounts();
      }, 1000);
    } else if (!oldVal && newVal) {
      console.log('[Miniflux Settings] Miniflux enabled, starting polling');
      await appStore.fetchFeeds();
      await appStore.startMinifluxStatusPolling();
      startStatusPolling();
      emit('settings-changed');
    }
  }
);

function formatSyncTime(timeStr: string | null): string {
  if (!timeStr) return t('setting.miniflux.never');
  const date = new Date(timeStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 60000);

  if (diffMins < 1) return t('setting.miniflux.justNow');
  if (diffMins < 60) return t('setting.miniflux.minsAgo', { n: diffMins });
  const diffHours = Math.floor(diffMins / 60);
  if (diffHours < 24) return t('setting.miniflux.hoursAgo', { n: diffHours });
  const diffDays = Math.floor(diffHours / 24);
  return t('setting.miniflux.daysAgo', { n: diffDays });
}
</script>

<template>
  <div class="setting-item">
    <div class="flex-1 flex items-center sm:items-start gap-2 sm:gap-3 min-w-0">
      <img
        src="/assets/plugin_icons/miniflux.svg"
        alt="Miniflux"
        class="w-5 h-5 sm:w-6 sm:h-6 mt-0.5 shrink-0"
      />
      <div class="flex-1 min-w-0">
        <div class="font-medium mb-0 sm:mb-1 text-sm sm:text-base">
          {{ t('setting.miniflux.enabled') }}
        </div>
        <div class="text-xs text-text-secondary hidden sm:block">
          {{ t('setting.miniflux.enabledDesc') }}
        </div>
      </div>
    </div>
    <input
      type="checkbox"
      :checked="props.settings.miniflux_enabled"
      class="toggle"
      @change="handleMinifluxToggle"
    />
  </div>
  <NestedSettingsContainer v-if="props.settings.miniflux_enabled">
    <SubSettingItem
      :icon="PhLink"
      :title="t('setting.miniflux.serverUrl')"
      :description="t('setting.miniflux.serverUrlDesc')"
      required
    >
      <InputControl
        type="url"
        :model-value="props.settings.miniflux_server_url"
        :placeholder="t('setting.miniflux.serverUrlPlaceholder')"
        width="md"
        @update:model-value="updateSetting('miniflux_server_url', $event)"
      />
    </SubSettingItem>

    <SubSettingItem
      :icon="PhKey"
      :title="t('setting.miniflux.apiKey')"
      :description="t('setting.miniflux.apiKeyDesc')"
      required
    >
      <InputControl
        type="password"
        :model-value="props.settings.miniflux_api_key"
        :placeholder="t('setting.miniflux.apiKeyPlaceholder')"
        width="md"
        @update:model-value="updateSetting('miniflux_api_key', $event)"
      />
    </SubSettingItem>

    <SubSettingItem
      :icon="PhTestTube"
      :title="t('setting.miniflux.testConnection')"
      :description="t('setting.miniflux.testConnectionDesc')"
    >
      <button class="btn-secondary" @click="testConnection">
        {{ t('setting.miniflux.test') }}
      </button>
    </SubSettingItem>

    <SubSettingItem
      :icon="PhCloudCheck"
      :title="t('setting.miniflux.syncNow')"
      :description="t('setting.miniflux.syncNowDesc')"
    >
      <template #description>
        <div>
          {{ t('setting.miniflux.syncNowDesc') }}
          <div class="text-xs text-text-secondary mt-1">
            {{ t('setting.miniflux.lastSync') }}:
            <span class="theme-number">{{ formatSyncTime(syncStatus.last_sync_time) }}</span>
          </div>
        </div>
      </template>
      <button class="btn-secondary" :disabled="isSyncing" @click="syncNow">
        <PhArrowClockwise :size="16" class="sm:w-5 sm:h-5" :class="{ 'animate-spin': isSyncing }" />
        {{ isSyncing ? t('setting.miniflux.syncing') : t('setting.miniflux.sync') }}
      </button>
    </SubSettingItem>
  </NestedSettingsContainer>
</template>

<style scoped>
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

.btn-secondary {
  @apply bg-bg-tertiary border border-border text-text-primary px-3 sm:px-4 py-1.5 sm:py-2 rounded-md cursor-pointer flex items-center gap-1.5 sm:gap-2 font-medium hover:bg-bg-secondary transition-colors;
}
.btn-secondary:disabled {
  @apply cursor-not-allowed opacity-50;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.animate-spin {
  animation: spin 1s linear infinite;
}

.theme-number {
  @apply text-accent font-semibold;
}
</style>
