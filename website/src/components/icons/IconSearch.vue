<template>
  <svg
    :class="iconClass"
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="2"
      d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
    />
  </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  class?: string;
}

const props = withDefaults(defineProps<Props>(), {
  class: 'w-7 h-7',
});

const iconClass = computed(() => props.class);
</script>
