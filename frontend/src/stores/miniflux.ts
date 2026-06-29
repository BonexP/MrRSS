import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { MinifluxFeed, MinifluxCategory, MinifluxEntry, MinifluxEntriesResponse, MinifluxCounters } from '@/types/models';

export interface MinifluxArticleFilter {
  feed_id?: number;
  category_id?: number;
  status?: string;
  starred?: boolean;
  limit?: number;
  offset?: number;
  search?: string;
}

export const useMinifluxStore = defineStore('miniflux', () => {
  const feeds = ref<MinifluxFeed[]>([]);
  const categories = ref<MinifluxCategory[]>([]);
  const counters = ref<MinifluxCounters>({ reads: {}, unreads: {} });
  const entries = ref<MinifluxEntry[]>([]);
  const totalEntries = ref(0);
  const isLoading = ref(false);
  const isEnabled = ref(false);

  async function fetchFeeds(): Promise<MinifluxFeed[]> {
    try {
      const res = await fetch('/api/miniflux/feeds');
      if (!res.ok) {
        if (res.status === 400) {
          console.log('[Miniflux Store] fetchFeeds: Miniflux not enabled (400)');
          isEnabled.value = false;
        } else {
          console.error('[Miniflux Store] fetchFeeds failed:', res.status, res.statusText);
        }
        return [];
      }
      const data: MinifluxFeed[] = await res.json();
      feeds.value = data;
      isEnabled.value = true;
      console.log('[Miniflux Store] fetchFeeds:', data.length, 'feeds');
      return data;
    } catch (e) {
      console.error('[Miniflux Store] fetchFeeds error:', e);
      return [];
    }
  }

  async function fetchCategories(): Promise<MinifluxCategory[]> {
    try {
      const res = await fetch('/api/miniflux/categories');
      if (!res.ok) {
        console.error('[Miniflux Store] fetchCategories failed:', res.status, res.statusText);
        return [];
      }
      const data: MinifluxCategory[] = await res.json();
      categories.value = data;
      console.log('[Miniflux Store] fetchCategories:', data.length, 'categories');
      return data;
    } catch (e) {
      console.error('[Miniflux Store] fetchCategories error:', e);
      return [];
    }
  }

  async function fetchCounters(): Promise<MinifluxCounters | null> {
    try {
      const res = await fetch('/api/miniflux/counters');
      if (!res.ok) {
        console.error('[Miniflux Store] fetchCounters failed:', res.status, res.statusText);
        return null;
      }
      const data: MinifluxCounters = await res.json();
      counters.value = data;
      console.log('[Miniflux Store] fetchCounters: unreads=%d, reads=%d',
        Object.keys(data.unreads || {}).length,
        Object.keys(data.reads || {}).length);
      return data;
    } catch (e) {
      console.error('[Miniflux Store] fetchCounters error:', e);
      return null;
    }
  }

  async function fetchEntries(filter: MinifluxArticleFilter = {}): Promise<MinifluxEntriesResponse | null> {
    isLoading.value = true;
    try {
      const params = new URLSearchParams();
      if (filter.feed_id) params.append('feed_id', String(filter.feed_id));
      if (filter.category_id) params.append('category_id', String(filter.category_id));
      if (filter.status) params.append('status', filter.status);
      if (filter.starred) params.append('starred', 'true');
      if (filter.limit) params.append('limit', String(filter.limit));
      if (filter.offset) params.append('offset', String(filter.offset));
      if (filter.search) params.append('search', filter.search);

      const qs = params.toString();
      const url = `/api/miniflux/articles${qs ? '?' + qs : ''}`;
      console.log('[Miniflux Store] fetchEntries:', url);
      const res = await fetch(url);
      if (!res.ok) {
        console.error('[Miniflux Store] fetchEntries failed:', res.status, res.statusText);
        return null;
      }
      const data: MinifluxEntriesResponse = await res.json();
      entries.value = data.entries;
      totalEntries.value = data.total;
      console.log('[Miniflux Store] fetchEntries: %d entries (total=%d)', data.entries.length, data.total);
      return data;
    } catch (e) {
      console.error('[Miniflux Store] fetchEntries error:', e);
      return null;
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchEntry(id: number): Promise<MinifluxEntry | null> {
    try {
      const res = await fetch(`/api/miniflux/articles/detail?id=${id}`);
      if (!res.ok) {
        console.error('[Miniflux Store] fetchEntry failed (id=%d):', id, res.status, res.statusText);
        return null;
      }
      const data = await res.json();
      console.log('[Miniflux Store] fetchEntry (id=%d): %q', id, data.title);
      return data;
    } catch (e) {
      console.error('[Miniflux Store] fetchEntry error (id=%d):', id, e);
      return null;
    }
  }

  async function updateStatus(entryIDs: number[], status: string, starred?: boolean): Promise<boolean> {
    try {
      const res = await fetch('/api/miniflux/articles/status', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ entry_ids: entryIDs, status, starred }),
      });
      if (!res.ok) {
        console.error('[Miniflux Store] updateStatus failed:', res.status, res.statusText);
      } else {
        console.log('[Miniflux Store] updateStatus: %d entries -> status=%s, starred=%s',
          entryIDs.length, status, starred);
      }
      return res.ok;
    } catch (e) {
      console.error('[Miniflux Store] updateStatus error:', e);
      return false;
    }
  }

  async function markFeedAsRead(feedID: number): Promise<boolean> {
    try {
      const res = await fetch(`/api/miniflux/feeds/mark-read?id=${feedID}`, { method: 'POST' });
      if (!res.ok) {
        console.error('[Miniflux Store] markFeedAsRead failed (feed=%d):', feedID, res.status, res.statusText);
      } else {
        console.log('[Miniflux Store] markFeedAsRead (feed=%d)', feedID);
      }
      return res.ok;
    } catch (e) {
      console.error('[Miniflux Store] markFeedAsRead error (feed=%d):', feedID, e);
      return false;
    }
  }

  async function markCategoryAsRead(categoryID: number): Promise<boolean> {
    try {
      const res = await fetch(`/api/miniflux/categories/mark-read?id=${categoryID}`, { method: 'POST' });
      if (!res.ok) {
        console.error('[Miniflux Store] markCategoryAsRead failed (category=%d):', categoryID, res.status, res.statusText);
      } else {
        console.log('[Miniflux Store] markCategoryAsRead (category=%d)', categoryID);
      }
      return res.ok;
    } catch (e) {
      console.error('[Miniflux Store] markCategoryAsRead error (category=%d):', categoryID, e);
      return false;
    }
  }

  async function markAllAsRead(): Promise<boolean> {
    try {
      const res = await fetch('/api/miniflux/mark-all-read', { method: 'POST' });
      if (!res.ok) {
        console.error('[Miniflux Store] markAllAsRead failed:', res.status, res.statusText);
      } else {
        console.log('[Miniflux Store] markAllAsRead: success');
      }
      return res.ok;
    } catch (e) {
      console.error('[Miniflux Store] markAllAsRead error:', e);
      return false;
    }
  }

  async function toggleBookmark(entryID: number): Promise<boolean> {
    try {
      const res = await fetch('/api/miniflux/articles/status', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ entry_ids: [entryID], starred: true }),
      });
      if (!res.ok) {
        console.error('[Miniflux Store] toggleBookmark failed (entry=%d):', entryID, res.status, res.statusText);
      } else {
        console.log('[Miniflux Store] toggleBookmark (entry=%d)', entryID);
      }
      return res.ok;
    } catch (e) {
      console.error('[Miniflux Store] toggleBookmark error (entry=%d):', entryID, e);
      return false;
    }
  }

  function reset() {
    feeds.value = [];
    categories.value = [];
    counters.value = { reads: {}, unreads: {} };
    entries.value = [];
    totalEntries.value = 0;
    isLoading.value = false;
  }

  return {
    feeds,
    categories,
    counters,
    entries,
    totalEntries,
    isLoading,
    isEnabled,
    fetchFeeds,
    fetchCategories,
    fetchCounters,
    fetchEntries,
    fetchEntry,
    updateStatus,
    markFeedAsRead,
    markCategoryAsRead,
    markAllAsRead,
    toggleBookmark,
    reset,
  };
});
