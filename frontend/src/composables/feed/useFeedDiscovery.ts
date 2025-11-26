/**
 * Composable for feed discovery operations
 */
import { ref, onUnmounted, type Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAppStore } from '@/stores/app';
import type { DiscoveredFeed, ProgressCounts, ProgressState } from '@/types/discovery';

interface StartResult {
  status: string;
  message?: string;
  total?: number;
}

export function useFeedDiscovery() {
  const { t } = useI18n();
  const store = useAppStore();

  const isDiscovering = ref(false);
  const discoveredFeeds: Ref<DiscoveredFeed[]> = ref([]);
  const selectedFeeds: Ref<Set<number>> = ref(new Set());
  const errorMessage = ref('');
  const progressMessage = ref('');
  const progressDetail = ref('');
  const progressCounts: Ref<ProgressCounts> = ref({ current: 0, total: 0, found: 0 });
  const isSubscribing = ref(false);
  let pollInterval: ReturnType<typeof setInterval> | null = null;

  /**
   * Extract hostname from URL
   */
  function getHostname(url: string): string {
    try {
      return new URL(url).hostname;
    } catch {
      return url;
    }
  }

  /**
   * Start polling for discovery progress
   */
  function startPolling(statusEndpoint: string) {
    pollInterval = setInterval(async () => {
      try {
        const statusResponse = await fetch(statusEndpoint);
        if (!statusResponse.ok) {
          throw new Error(`Status check failed: ${statusResponse.status}`);
        }

        const state: ProgressState = await statusResponse.json();

        // Update progress
        if (state.progress) {
          progressMessage.value = state.progress.message || '';
          progressDetail.value = state.progress.detail || '';

          if (state.progress.current !== undefined && state.progress.total !== undefined) {
            progressCounts.value.current = state.progress.current;
            progressCounts.value.total = state.progress.total;
          }
          if (state.progress.found_count !== undefined) {
            progressCounts.value.found = state.progress.found_count;
          }
        }

        // Check if complete
        if (state.is_complete) {
          if (pollInterval) {
            clearInterval(pollInterval);
            pollInterval = null;
          }
          isDiscovering.value = false;

          if (state.error) {
            errorMessage.value = state.error;
            window.showToast(state.error, 'error');
          } else if (state.feeds) {
            discoveredFeeds.value = state.feeds;
            // Auto-select all discovered feeds
            selectedFeeds.value = new Set(state.feeds.map((_, idx) => idx));
            if (state.feeds.length > 0) {
              window.showToast(t('discoveredFeeds', { count: state.feeds.length }), 'success');
            } else {
              window.showToast(t('noFeedsDiscovered'), 'info');
            }
          }
        }
      } catch (pollError) {
        console.error('Error polling status:', pollError);
        if (pollInterval) {
          clearInterval(pollInterval);
          pollInterval = null;
        }
        isDiscovering.value = false;
        errorMessage.value = t('errorPollingStatus');
      }
    }, 500); // Poll every 500ms
  }

  /**
   * Start single feed discovery
   */
  async function startSingleFeedDiscovery(feedId: number) {
    isDiscovering.value = true;
    errorMessage.value = '';
    discoveredFeeds.value = [];
    selectedFeeds.value.clear();
    progressMessage.value = t('fetchingHomepage');
    progressDetail.value = '';
    progressCounts.value = { current: 0, total: 0, found: 0 };

    if (pollInterval) {
      clearInterval(pollInterval);
      pollInterval = null;
    }

    try {
      if (!feedId) {
        throw new Error('Invalid feed ID');
      }

      // Clear previous state
      await fetch('/api/feeds/discover/clear', { method: 'POST' });

      // Start discovery
      const startResponse = await fetch('/api/feeds/discover/start', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ feed_id: feedId }),
      });

      if (!startResponse.ok) {
        throw new Error(`Failed to start discovery: ${startResponse.status}`);
      }

      const startResult: StartResult = await startResponse.json();

      if (startResult.status === 'started') {
        startPolling('/api/feeds/discover/status');
      } else {
        throw new Error(startResult.message || 'Failed to start discovery');
      }
    } catch (error) {
      console.error('Error starting discovery:', error);
      isDiscovering.value = false;
      errorMessage.value = error instanceof Error ? error.message : String(error);
      window.showToast(t('errorDiscoveringFeeds'), 'error');
    }
  }

  /**
   * Start batch discovery (all feeds)
   */
  async function startBatchDiscovery() {
    isDiscovering.value = true;
    errorMessage.value = '';
    discoveredFeeds.value = [];
    selectedFeeds.value.clear();
    progressMessage.value = t('preparingDiscovery');
    progressDetail.value = '';
    progressCounts.value = { current: 0, total: 0, found: 0 };

    if (pollInterval) {
      clearInterval(pollInterval);
      pollInterval = null;
    }

    try {
      // Clear previous state
      await fetch('/api/feeds/discover-all/clear', { method: 'POST' });

      // Start batch discovery
      const startResponse = await fetch('/api/feeds/discover-all/start', {
        method: 'POST',
      });

      if (!startResponse.ok) {
        throw new Error(`Failed to start batch discovery: ${startResponse.status}`);
      }

      const startResult: StartResult = await startResponse.json();

      if (startResult.status === 'started') {
        progressCounts.value.total = startResult.total || 0;
        startPolling('/api/feeds/discover-all/status');
      } else {
        throw new Error(startResult.message || 'Failed to start batch discovery');
      }
    } catch (error) {
      console.error('Error starting batch discovery:', error);
      isDiscovering.value = false;
      errorMessage.value = error instanceof Error ? error.message : String(error);
      window.showToast(t('errorDiscoveringFeeds'), 'error');
    }
  }

  /**
   * Toggle feed selection
   */
  function toggleFeedSelection(index: number) {
    if (selectedFeeds.value.has(index)) {
      selectedFeeds.value.delete(index);
    } else {
      selectedFeeds.value.add(index);
    }
  }

  /**
   * Toggle all feeds selection
   */
  function toggleAllFeeds() {
    if (selectedFeeds.value.size === discoveredFeeds.value.length) {
      selectedFeeds.value.clear();
    } else {
      selectedFeeds.value = new Set(discoveredFeeds.value.map((_, idx) => idx));
    }
  }

  /**
   * Subscribe to selected feeds
   */
  async function subscribeToFeeds() {
    if (selectedFeeds.value.size === 0) {
      window.showToast(t('pleaseSelectFeeds'), 'warning');
      return;
    }

    isSubscribing.value = true;
    const feedsToSubscribe = Array.from(selectedFeeds.value)
      .map((idx) => discoveredFeeds.value[idx])
      .filter(Boolean);

    let successCount = 0;
    let failCount = 0;

    for (const feed of feedsToSubscribe) {
      try {
        const response = await fetch('/api/feeds/add', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            title: feed.name,
            url: feed.rss_feed,
            category: '',
          }),
        });

        if (response.ok) {
          successCount++;
        } else {
          failCount++;
        }
      } catch (error) {
        console.error('Error subscribing to feed:', error);
        failCount++;
      }
    }

    isSubscribing.value = false;

    if (successCount > 0) {
      await store.fetchFeeds();
      window.showToast(t('feedsSubscribedSuccess', { count: successCount }), 'success');
    }

    if (failCount > 0) {
      window.showToast(t('someFeedsFailedToSubscribe', { count: failCount }), 'error');
    }

    return { successCount, failCount };
  }

  /**
   * Cleanup on unmount
   */
  onUnmounted(() => {
    if (pollInterval) {
      clearInterval(pollInterval);
      pollInterval = null;
    }
  });

  return {
    isDiscovering,
    discoveredFeeds,
    selectedFeeds,
    errorMessage,
    progressMessage,
    progressDetail,
    progressCounts,
    isSubscribing,
    getHostname,
    startSingleFeedDiscovery,
    startBatchDiscovery,
    toggleFeedSelection,
    toggleAllFeeds,
    subscribeToFeeds,
  };
}
