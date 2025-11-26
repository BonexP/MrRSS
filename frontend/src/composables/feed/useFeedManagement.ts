/**
 * Composable for feed management operations in settings
 */
import { useI18n } from 'vue-i18n';
import { useAppStore } from '@/stores/app';
import type { Feed } from '@/types/models';

export function useFeedManagement() {
  const { t } = useI18n();
  const store = useAppStore();

  /**
   * Import OPML file
   */
  function handleImportOPML(event: Event) {
    const target = event.target as HTMLInputElement;
    const file = target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e: ProgressEvent<FileReader>) => {
      const content = e.target?.result;
      fetch('/api/opml/import', {
        method: 'POST',
        headers: {
          'Content-Type': 'text/xml',
        },
        body: content,
      }).then(async (res) => {
        if (res.ok) {
          window.showToast(t('opmlImportedSuccess'), 'success');
          store.fetchFeeds();
        } else {
          const text = await res.text();
          window.showToast(t('importFailed', { error: text }), 'error');
        }
      });
    };
    reader.readAsText(file);
  }

  /**
   * Export OPML file
   */
  function handleExportOPML() {
    window.location.href = '/api/opml/export';
  }

  /**
   * Clean up old articles from database
   */
  async function handleCleanupDatabase() {
    const confirmed = await window.showConfirm({
      title: t('cleanDatabaseTitle'),
      message: t('cleanDatabaseMessage'),
      confirmText: t('clean'),
      cancelText: t('cancel'),
      isDanger: true,
    });
    if (!confirmed) return;

    try {
      const res = await fetch('/api/articles/cleanup', { method: 'POST' });
      if (res.ok) {
        const result = await res.json();
        window.showToast(t('databaseCleanedSuccess', { count: result.deleted }), 'success');
        store.fetchArticles();
      } else {
        window.showToast(t('errorCleaningDatabase'), 'error');
      }
    } catch (e) {
      console.error('Error cleaning database:', e);
      window.showToast(t('errorCleaningDatabase'), 'error');
    }
  }

  /**
   * Add new feed
   */
  function handleAddFeed() {
    window.dispatchEvent(new CustomEvent('show-add-feed'));
  }

  /**
   * Edit existing feed
   */
  function handleEditFeed(feed: Feed) {
    window.dispatchEvent(new CustomEvent('show-edit-feed', { detail: feed }));
  }

  /**
   * Delete a single feed
   */
  async function handleDeleteFeed(id: number) {
    const confirmed = await window.showConfirm({
      title: t('deleteFeedTitle'),
      message: t('deleteFeedMessage'),
      confirmText: t('delete'),
      cancelText: t('cancel'),
      isDanger: true,
    });
    if (!confirmed) return;

    await fetch(`/api/feeds/delete?id=${id}`, { method: 'POST' });
    store.fetchFeeds();
    window.showToast(t('feedDeletedSuccess'), 'success');
  }

  /**
   * Delete multiple feeds
   */
  async function handleBatchDelete(selectedIds: number[]) {
    const confirmed = await window.showConfirm({
      title: t('deleteMultipleFeedsTitle'),
      message: t('deleteMultipleFeedsMessage', { count: selectedIds.length }),
      confirmText: t('delete'),
      cancelText: t('cancel'),
      isDanger: true,
    });
    if (!confirmed) return;

    const promises = selectedIds.map((id: number) =>
      fetch(`/api/feeds/delete?id=${id}`, { method: 'POST' })
    );
    await Promise.all(promises);
    store.fetchFeeds();
    window.showToast(t('feedsDeletedSuccess'), 'success');
  }

  /**
   * Move multiple feeds to a new category
   */
  async function handleBatchMove(selectedIds: number[]) {
    if (!store.feeds) return;

    const newCategory = await window.showInput({
      title: t('moveFeeds'),
      message: t('enterCategoryName'),
      placeholder: t('categoryPlaceholder'),
      confirmText: t('move'),
      cancelText: t('cancel'),
    });
    if (newCategory === null) return;

    const promises = selectedIds.map((id: number) => {
      const feed = store.feeds.find((f) => f.id === id);
      if (!feed) return Promise.resolve();
      return fetch('/api/feeds/update', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          id: feed.id,
          title: feed.title,
          url: feed.url,
          category: newCategory,
        }),
      });
    });

    await Promise.all(promises);
    store.fetchFeeds();
    window.showToast(t('feedsMovedSuccess'), 'success');
  }

  return {
    handleImportOPML,
    handleExportOPML,
    handleCleanupDatabase,
    handleAddFeed,
    handleEditFeed,
    handleDeleteFeed,
    handleBatchDelete,
    handleBatchMove,
  };
}
