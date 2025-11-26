import { computed, type ComputedRef } from 'vue';
import { useAppStore } from '@/stores/app';

export interface Condition {
  id: number;
  logic?: 'and' | 'or' | null;
  negate: boolean;
  field: string;
  operator?: string | null;
  value: string;
  values: string[];
}

export interface FieldOption {
  value: string;
  labelKey: string;
  multiSelect: boolean;
  booleanField?: boolean;
}

export interface ActionOption {
  value: string;
  labelKey: string;
}

export function useRuleOptions() {
  const store = useAppStore();

  // Field options for conditions
  const fieldOptions: FieldOption[] = [
    { value: 'feed_name', labelKey: 'feedName', multiSelect: true },
    { value: 'feed_category', labelKey: 'feedCategory', multiSelect: true },
    { value: 'article_title', labelKey: 'articleTitle', multiSelect: false },
    { value: 'published_after', labelKey: 'publishedAfter', multiSelect: false },
    { value: 'published_before', labelKey: 'publishedBefore', multiSelect: false },
    { value: 'is_read', labelKey: 'readStatus', multiSelect: false, booleanField: true },
    { value: 'is_favorite', labelKey: 'favoriteStatus', multiSelect: false, booleanField: true },
    { value: 'is_hidden', labelKey: 'hiddenStatus', multiSelect: false, booleanField: true },
  ];

  // Operator options for article title
  const textOperatorOptions: Array<{ value: string; labelKey: string }> = [
    { value: 'contains', labelKey: 'contains' },
    { value: 'exact', labelKey: 'exactMatch' },
  ];

  // Boolean value options
  const booleanOptions: Array<{ value: string; labelKey: string }> = [
    { value: 'true', labelKey: 'yes' },
    { value: 'false', labelKey: 'no' },
  ];

  // Action options
  const actionOptions: ActionOption[] = [
    { value: 'favorite', labelKey: 'actionFavorite' },
    { value: 'unfavorite', labelKey: 'actionUnfavorite' },
    { value: 'hide', labelKey: 'actionHide' },
    { value: 'unhide', labelKey: 'actionUnhide' },
    { value: 'mark_read', labelKey: 'actionMarkRead' },
    { value: 'mark_unread', labelKey: 'actionMarkUnread' },
  ];

  // Feed names for multi-select
  const feedNames: ComputedRef<string[]> = computed(() => {
    return store.feeds.map((f) => f.title);
  });

  // Feed categories for multi-select
  const feedCategories: ComputedRef<string[]> = computed(() => {
    const categories = new Set<string>();
    store.feeds.forEach((f) => {
      if (f.category) {
        categories.add(f.category);
      }
    });
    return Array.from(categories);
  });

  return {
    fieldOptions,
    textOperatorOptions,
    booleanOptions,
    actionOptions,
    feedNames,
    feedCategories,
  };
}

// Helper functions for field types
export function isDateField(field: string): boolean {
  return field === 'published_after' || field === 'published_before';
}

export function isMultiSelectField(field: string): boolean {
  return field === 'feed_name' || field === 'feed_category';
}

export function isBooleanField(field: string): boolean {
  return field === 'is_read' || field === 'is_favorite' || field === 'is_hidden';
}

export function needsOperator(field: string): boolean {
  return field === 'article_title';
}
