<script setup lang="ts">
import { computed, type ComputedRef } from 'vue';
import { useI18n } from 'vue-i18n';
import { PhTrash } from '@phosphor-icons/vue';

interface ActionOption {
  value: string;
  labelKey: string;
}

interface Props {
  action: string;
  index: number;
  selectedActions: string[];
  allActionOptions: ActionOption[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  update: [value: string];
  remove: [];
}>();

const { t } = useI18n();

// Get available actions (exclude already selected ones, except current)
const availableActions: ComputedRef<ActionOption[]> = computed(() => {
  const selectedSet = new Set(props.selectedActions);
  return props.allActionOptions.filter(
    (opt) => !selectedSet.has(opt.value) || opt.value === props.action
  );
});

function handleUpdate(event: Event): void {
  const value = (event.target as HTMLSelectElement).value;
  emit('update', value);
}
</script>

<template>
  <div class="action-row">
    <span class="text-xs text-text-secondary">{{ index + 1 }}.</span>
    <select :value="action" class="select-field flex-1" @change="handleUpdate">
      <option v-for="opt in availableActions" :key="opt.value" :value="opt.value">
        {{ t(opt.labelKey) }}
      </option>
    </select>
    <button class="btn-danger-icon" :title="t('removeAction')" @click="emit('remove')">
      <PhTrash :size="16" />
    </button>
  </div>
</template>

<style scoped>
@reference "../../../style.css";

.action-row {
  @apply flex items-center gap-2 p-2 bg-bg-secondary border border-border rounded-lg;
}

.select-field {
  @apply p-2 border border-border rounded-md bg-bg-primary text-text-primary text-sm focus:border-accent focus:outline-none transition-colors cursor-pointer;
}

.btn-danger-icon {
  @apply p-2 rounded-lg text-red-500 hover:bg-red-500/10 transition-colors cursor-pointer;
}
</style>
