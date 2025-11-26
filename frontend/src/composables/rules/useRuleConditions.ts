import { ref, type Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import type { Condition } from './useRuleOptions';
import { isDateField, isMultiSelectField, isBooleanField } from './useRuleOptions';

export function useRuleConditions() {
  const { t, locale } = useI18n();
  const openDropdownIndex: Ref<number | null> = ref(null);

  function addCondition(conditions: Condition[]): void {
    conditions.push({
      id: Date.now(),
      logic: conditions.length > 0 ? 'and' : null,
      negate: false,
      field: 'article_title',
      operator: 'contains',
      value: '',
      values: [],
    });
  }

  function removeCondition(conditions: Condition[], index: number): void {
    conditions.splice(index, 1);
    if (conditions.length > 0 && index === 0) {
      conditions[0].logic = null;
    }
  }

  function onFieldChange(condition: Condition): void {
    if (isDateField(condition.field)) {
      condition.operator = null;
      condition.value = '';
      condition.values = [];
    } else if (isMultiSelectField(condition.field)) {
      condition.operator = 'contains';
      condition.value = '';
      condition.values = [];
    } else if (isBooleanField(condition.field)) {
      condition.operator = null;
      condition.value = 'true';
      condition.values = [];
    } else {
      condition.operator = 'contains';
      condition.value = '';
      condition.values = [];
    }
  }

  function toggleNegate(condition: Condition): void {
    condition.negate = !condition.negate;
  }

  function toggleDropdown(index: number): void {
    if (openDropdownIndex.value === index) {
      openDropdownIndex.value = null;
    } else {
      openDropdownIndex.value = index;
    }
  }

  function toggleMultiSelectValue(condition: Condition, val: string): void {
    const idx = condition.values.indexOf(val);
    if (idx > -1) {
      condition.values.splice(idx, 1);
    } else {
      condition.values.push(val);
    }
  }

  function getMultiSelectDisplayText(condition: Condition, labelKey: string): string {
    if (!condition.values || condition.values.length === 0) {
      return t(labelKey);
    }

    if (condition.values.length === 1) {
      return condition.values[0];
    }

    const firstItem = condition.values[0];
    const totalCount = condition.values.length;
    const remaining = totalCount - 1;

    if (locale.value === 'zh') {
      return `${firstItem} ${t('andNMore', { count: totalCount })}`;
    }
    return `${firstItem} ${t('andNMore', { count: remaining })}`;
  }

  return {
    openDropdownIndex,
    addCondition,
    removeCondition,
    onFieldChange,
    toggleNegate,
    toggleDropdown,
    toggleMultiSelectValue,
    getMultiSelectDisplayText,
  };
}
