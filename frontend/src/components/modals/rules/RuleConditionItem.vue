<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { PhProhibit, PhTrash } from '@phosphor-icons/vue';
import {
  useRuleOptions,
  type Condition,
  isDateField,
  isBooleanField,
  needsOperator,
} from '@/composables/rules/useRuleOptions';

const { t } = useI18n();

const { fieldOptions, textOperatorOptions, booleanOptions, feedNames, feedCategories } =
  useRuleOptions();

interface Props {
  condition: Condition;
  index: number;
  isDropdownOpen: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  'update:field': [value: string];
  'update:operator': [value: string];
  'update:value': [value: string];
  'update:values': [values: string[]];
  'update:negate': [];
  'toggle-dropdown': [];
  remove: [];
}>();

function handleFieldChange(event: Event): void {
  const target = event.target as HTMLSelectElement;
  emit('update:field', target.value);
}

function handleOperatorChange(event: Event): void {
  const target = event.target as HTMLSelectElement;
  emit('update:operator', target.value);
}

function handleValueChange(event: Event): void {
  const target = event.target as HTMLInputElement;
  emit('update:value', target.value);
}

function handleToggleMultiSelectValue(value: string): void {
  const currentValues = props.condition.values || [];
  const newValues = currentValues.includes(value)
    ? currentValues.filter((v) => v !== value)
    : [...currentValues, value];
  emit('update:values', newValues);
}

function getMultiSelectDisplayText(): string {
  const values = props.condition.values || [];
  if (values.length === 0) return t('selectItems');
  if (values.length === 1) return values[0];
  return t('itemsSelected', { count: values.length });
}
</script>

<template>
  <div class="condition-row bg-bg-secondary border border-border rounded-lg p-3">
    <div class="flex flex-wrap gap-2 items-end">
      <!-- NOT toggle button -->
      <div class="flex-shrink-0">
        <label class="block text-xs text-text-secondary mb-1">&nbsp;</label>
        <button
          @click="emit('update:negate')"
          :class="['not-btn', condition.negate ? 'active' : '']"
          :title="t('not')"
        >
          <PhProhibit :size="16" />
          <span class="text-xs font-medium">{{ t('not') }}</span>
        </button>
      </div>

      <!-- Field selector -->
      <div class="flex-1 min-w-[130px]">
        <label class="block text-xs text-text-secondary mb-1">{{ t('filterField') }}</label>
        <select :value="condition.field" @change="handleFieldChange" class="select-field w-full">
          <option v-for="opt in fieldOptions" :key="opt.value" :value="opt.value">
            {{ t(opt.labelKey) }}
          </option>
        </select>
      </div>

      <!-- Operator selector (only for article_title) -->
      <div v-if="needsOperator(condition.field)" class="w-28">
        <label class="block text-xs text-text-secondary mb-1">{{ t('filterOperator') }}</label>
        <select
          :value="condition.operator"
          @change="handleOperatorChange"
          class="select-field w-full"
        >
          <option v-for="opt in textOperatorOptions" :key="opt.value" :value="opt.value">
            {{ t(opt.labelKey) }}
          </option>
        </select>
      </div>

      <!-- Value input -->
      <div class="flex-1 min-w-[140px]">
        <label class="block text-xs text-text-secondary mb-1">{{ t('filterValue') }}</label>

        <!-- Date input -->
        <input
          v-if="isDateField(condition.field)"
          type="date"
          :value="condition.value"
          @input="handleValueChange"
          class="date-field w-full"
        />

        <!-- Boolean select -->
        <select
          v-else-if="isBooleanField(condition.field)"
          :value="condition.value"
          @change="handleValueChange"
          class="select-field w-full"
        >
          <option v-for="opt in booleanOptions" :key="opt.value" :value="opt.value">
            {{ t(opt.labelKey) }}
          </option>
        </select>

        <!-- Multi-select dropdown for feed name -->
        <div v-else-if="condition.field === 'feed_name'" class="dropdown-container">
          <button type="button" @click="emit('toggle-dropdown')" class="dropdown-trigger">
            <span class="dropdown-text truncate">{{ getMultiSelectDisplayText() }}</span>
            <span class="dropdown-arrow">▼</span>
          </button>
          <div v-if="isDropdownOpen" class="dropdown-menu dropdown-down">
            <div
              v-for="name in feedNames"
              :key="name"
              @click.stop="handleToggleMultiSelectValue(name)"
              :class="['dropdown-option', condition.values.includes(name) ? 'selected' : '']"
            >
              <input
                type="checkbox"
                :checked="condition.values.includes(name)"
                class="checkbox-input"
                tabindex="-1"
              />
              <span class="truncate">{{ name }}</span>
            </div>
            <div v-if="feedNames.length === 0" class="text-text-secondary text-sm p-2">
              {{ t('noArticles') }}
            </div>
          </div>
        </div>

        <!-- Multi-select dropdown for category -->
        <div v-else-if="condition.field === 'feed_category'" class="dropdown-container">
          <button type="button" @click="emit('toggle-dropdown')" class="dropdown-trigger">
            <span class="dropdown-text truncate">{{ getMultiSelectDisplayText() }}</span>
            <span class="dropdown-arrow">▼</span>
          </button>
          <div v-if="isDropdownOpen" class="dropdown-menu dropdown-down">
            <div
              v-for="cat in feedCategories"
              :key="cat"
              @click.stop="handleToggleMultiSelectValue(cat as string)"
              :class="[
                'dropdown-option',
                condition.values.includes(cat as string) ? 'selected' : '',
              ]"
            >
              <input
                type="checkbox"
                :checked="condition.values.includes(cat as string)"
                class="checkbox-input"
                tabindex="-1"
              />
              <span class="truncate">{{ cat }}</span>
            </div>
            <div v-if="feedCategories.length === 0" class="text-text-secondary text-sm p-2">
              {{ t('noArticles') }}
            </div>
          </div>
        </div>

        <!-- Regular text input -->
        <input
          v-else
          type="text"
          :value="condition.value"
          @input="handleValueChange"
          class="input-field w-full"
          :placeholder="t('filterValue')"
        />
      </div>

      <!-- Remove button -->
      <div class="flex-shrink-0">
        <label class="block text-xs text-text-secondary mb-1">&nbsp;</label>
        <button @click="emit('remove')" class="btn-danger-icon" :title="t('removeCondition')">
          <PhTrash :size="18" />
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.input-field {
  @apply p-2 border border-border rounded-md bg-bg-primary text-text-primary text-sm focus:border-accent focus:outline-none transition-colors;
}
.select-field {
  @apply p-2 border border-border rounded-md bg-bg-primary text-text-primary text-sm focus:border-accent focus:outline-none transition-colors cursor-pointer;
}
.date-field {
  @apply p-2 border border-border rounded-md bg-bg-primary text-text-primary text-sm focus:border-accent focus:outline-none transition-colors cursor-pointer;
  color-scheme: light dark;
}
.btn-danger-icon {
  @apply p-2 rounded-lg text-red-500 hover:bg-red-500/10 transition-colors cursor-pointer;
}

/* NOT button styling */
.not-btn {
  @apply flex items-center gap-1 px-2 py-2 rounded-md border transition-all cursor-pointer;
  @apply text-text-secondary bg-bg-primary border-border;
}
.not-btn:hover {
  @apply border-red-400 text-red-500;
}
.not-btn.active {
  @apply bg-red-500/10 border-red-500 text-red-500;
}

/* Dropdown multi-select styling */
.dropdown-container {
  @apply relative;
}
.dropdown-trigger {
  @apply w-full p-2 border border-border rounded-md bg-bg-primary text-text-primary text-sm;
  @apply flex items-center justify-between cursor-pointer hover:border-accent transition-colors;
}
.dropdown-text {
  @apply flex-1 text-left;
}
.dropdown-arrow {
  @apply text-text-secondary text-xs ml-2;
}
.dropdown-menu {
  @apply absolute left-0 right-0 border border-border rounded-md bg-bg-primary;
  @apply max-h-40 overflow-y-auto z-50 shadow-lg;
}
.dropdown-menu.dropdown-down {
  @apply top-full mt-1;
}
.dropdown-option {
  @apply flex items-center gap-2 px-3 py-2 cursor-pointer text-sm text-text-primary hover:bg-bg-tertiary;
}
.dropdown-option.selected {
  background-color: rgba(59, 130, 246, 0.1);
}
.checkbox-input {
  @apply w-4 h-4 accent-accent cursor-pointer;
}
</style>
