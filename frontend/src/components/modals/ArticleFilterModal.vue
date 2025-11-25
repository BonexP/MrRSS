<script setup>
import { ref, computed, onMounted } from 'vue';
import { store } from '../../store.js';
import { PhPlus, PhTrash, PhFunnel } from "@phosphor-icons/vue";

const props = defineProps({
    show: { type: Boolean, default: false },
    currentFilters: { type: Array, default: () => [] }
});

const emit = defineEmits(['close', 'apply']);

// Filter conditions
const conditions = ref([]);

// Field options
const fieldOptions = [
    { value: 'feed_name', labelKey: 'feedName' },
    { value: 'feed_category', labelKey: 'feedCategory' },
    { value: 'article_title', labelKey: 'articleTitle' },
    { value: 'published_after', labelKey: 'publishedAfter' },
    { value: 'published_before', labelKey: 'publishedBefore' }
];

// Operator options for text fields
const textOperatorOptions = [
    { value: 'contains', labelKey: 'contains' },
    { value: 'exact', labelKey: 'exactMatch' }
];

// Logic options
const logicOptions = [
    { value: 'and', labelKey: 'and' },
    { value: 'or', labelKey: 'or' },
    { value: 'not', labelKey: 'not' }
];

// Feed names for autocomplete
const feedNames = computed(() => {
    return store.feeds.map(f => f.title);
});

// Feed categories for autocomplete
const feedCategories = computed(() => {
    const categories = new Set();
    store.feeds.forEach(f => {
        if (f.category) {
            categories.add(f.category);
        }
    });
    return Array.from(categories);
});

onMounted(() => {
    // Load existing filters if provided
    if (props.currentFilters && props.currentFilters.length > 0) {
        conditions.value = JSON.parse(JSON.stringify(props.currentFilters));
    }
});

function addCondition() {
    conditions.value.push({
        id: Date.now(),
        logic: conditions.value.length > 0 ? 'and' : null,
        field: 'article_title',
        operator: 'contains',
        value: ''
    });
}

function removeCondition(index) {
    conditions.value.splice(index, 1);
    // Reset first condition's logic to null
    if (conditions.value.length > 0 && index === 0) {
        conditions.value[0].logic = null;
    }
}

function isDateField(field) {
    return field === 'published_after' || field === 'published_before';
}

function onFieldChange(index) {
    const condition = conditions.value[index];
    if (isDateField(condition.field)) {
        condition.operator = null;
        condition.value = '';
    } else {
        condition.operator = 'contains';
    }
}

function clearFilters() {
    conditions.value = [];
}

function applyFilters() {
    // Validate conditions - only include conditions with a value
    const validConditions = conditions.value.filter(c => c.value);
    
    emit('apply', validConditions);
    emit('close');
}

function close() {
    emit('close');
}
</script>

<template>
    <div v-if="show" class="fixed inset-0 z-[60] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
        <div class="bg-bg-primary w-full max-w-2xl max-h-[90vh] flex flex-col rounded-2xl shadow-2xl border border-border overflow-hidden animate-fade-in">
            <!-- Header -->
            <div class="p-4 sm:p-5 border-b border-border flex justify-between items-center shrink-0">
                <h3 class="text-lg font-semibold m-0 flex items-center gap-2">
                    <PhFunnel :size="20" />
                    {{ store.i18n.t('filterArticles') }}
                </h3>
                <span @click="close" class="text-2xl cursor-pointer text-text-secondary hover:text-text-primary">&times;</span>
            </div>
            
            <!-- Content -->
            <div class="flex-1 overflow-y-auto p-4 sm:p-6">
                <!-- Empty state -->
                <div v-if="conditions.length === 0" class="text-center text-text-secondary py-8">
                    <PhFunnel :size="48" class="mx-auto mb-3 opacity-50" />
                    <p>{{ store.i18n.t('noFiltersApplied') }}</p>
                </div>
                
                <!-- Condition list -->
                <div v-else class="space-y-3">
                    <div v-for="(condition, index) in conditions" :key="condition.id" class="condition-row bg-bg-secondary border border-border rounded-lg p-3 sm:p-4">
                        <!-- Logic connector (AND/OR/NOT) for non-first conditions -->
                        <div v-if="index > 0" class="mb-3">
                            <select v-model="condition.logic" class="select-field w-24">
                                <option v-for="opt in logicOptions" :key="opt.value" :value="opt.value">
                                    {{ store.i18n.t(opt.labelKey) }}
                                </option>
                            </select>
                        </div>
                        
                        <div class="flex flex-wrap gap-2 items-end">
                            <!-- Field selector -->
                            <div class="flex-1 min-w-[140px]">
                                <label class="block text-xs text-text-secondary mb-1">{{ store.i18n.t('filterField') }}</label>
                                <select v-model="condition.field" @change="onFieldChange(index)" class="select-field w-full">
                                    <option v-for="opt in fieldOptions" :key="opt.value" :value="opt.value">
                                        {{ store.i18n.t(opt.labelKey) }}
                                    </option>
                                </select>
                            </div>
                            
                            <!-- Operator selector (only for text fields) -->
                            <div v-if="!isDateField(condition.field)" class="w-28">
                                <label class="block text-xs text-text-secondary mb-1">{{ store.i18n.t('filterOperator') }}</label>
                                <select v-model="condition.operator" class="select-field w-full">
                                    <option v-for="opt in textOperatorOptions" :key="opt.value" :value="opt.value">
                                        {{ store.i18n.t(opt.labelKey) }}
                                    </option>
                                </select>
                            </div>
                            
                            <!-- Value input -->
                            <div class="flex-1 min-w-[140px]">
                                <label class="block text-xs text-text-secondary mb-1">{{ store.i18n.t('filterValue') }}</label>
                                
                                <!-- Date input for date fields -->
                                <input v-if="isDateField(condition.field)" 
                                       type="date" 
                                       v-model="condition.value" 
                                       class="date-field w-full">
                                
                                <!-- Dropdown select for feed name -->
                                <select v-else-if="condition.field === 'feed_name'" 
                                        v-model="condition.value" 
                                        class="select-field w-full">
                                    <option value="">{{ store.i18n.t('feedName') }}</option>
                                    <option v-for="name in feedNames" :key="name" :value="name">{{ name }}</option>
                                </select>
                                
                                <!-- Dropdown select for category -->
                                <select v-else-if="condition.field === 'feed_category'" 
                                        v-model="condition.value" 
                                        class="select-field w-full">
                                    <option value="">{{ store.i18n.t('feedCategory') }}</option>
                                    <option v-for="cat in feedCategories" :key="cat" :value="cat">{{ cat }}</option>
                                </select>
                                
                                <!-- Regular text input -->
                                <input v-else 
                                       type="text" 
                                       v-model="condition.value" 
                                       class="input-field w-full"
                                       :placeholder="store.i18n.t('filterValue')">
                            </div>
                            
                            <!-- Remove button -->
                            <div class="flex-shrink-0">
                                <button @click="removeCondition(index)" class="btn-danger-icon" :title="store.i18n.t('removeCondition')">
                                    <PhTrash :size="18" />
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
                
                <!-- Add condition button -->
                <button @click="addCondition" class="btn-secondary w-full mt-4 flex items-center justify-center gap-2">
                    <PhPlus :size="18" />
                    {{ store.i18n.t('addCondition') }}
                </button>
            </div>
            
            <!-- Footer -->
            <div class="p-4 sm:p-5 border-t border-border bg-bg-secondary flex justify-between gap-3 shrink-0">
                <button @click="clearFilters" class="btn-secondary" :disabled="conditions.length === 0">
                    {{ store.i18n.t('clearFilters') }}
                </button>
                <button @click="applyFilters" class="btn-primary">
                    {{ store.i18n.t('applyFilters') }}
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
.btn-primary {
    @apply bg-accent text-white border-none px-5 py-2.5 rounded-lg cursor-pointer font-semibold hover:bg-accent-hover transition-colors disabled:opacity-50 disabled:cursor-not-allowed;
}
.btn-secondary {
    @apply bg-bg-tertiary text-text-primary border border-border px-4 py-2.5 rounded-lg cursor-pointer font-medium hover:bg-bg-secondary transition-colors disabled:opacity-50 disabled:cursor-not-allowed;
}
.btn-danger-icon {
    @apply p-2 rounded-lg text-red-500 hover:bg-red-500/10 transition-colors cursor-pointer;
}
.animate-fade-in {
    animation: modalFadeIn 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
@keyframes modalFadeIn {
    from { transform: translateY(-20px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
}
</style>
