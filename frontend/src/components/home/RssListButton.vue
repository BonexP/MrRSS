<script lang="ts" setup>
import { Icon } from '@iconify/vue';
import { useRouter } from 'vue-router';

const props = defineProps({
  isRefreshing: Boolean
})

const router = useRouter();

const goToSettings = () => {
  router.push('/settings');
}
</script>

<template>
  <nav>
    <div class="util">
      <button class="btn" title="Filter">
        <Icon icon="material-symbols:filter-list" />
      </button>
    </div>
    <div class="function">
      <button class="btn" @click="$emit('handle-click-refresh')" title="Refresh">
        <Icon v-if="!isRefreshing" icon="material-symbols:refresh" />
        <Icon v-else icon="material-symbols:error" disabled="true" />
      </button>
      <button class="btn" title="Translate">
        <Icon icon="material-symbols:g-translate" />
      </button>
      <button class="btn" @click="$emit('delete-history-content')" title="Delete history">
        <Icon icon="material-symbols:auto-delete-outline" />
      </button>
      <button class="btn" @click="goToSettings" title="Settings">
        <Icon icon="material-symbols:settings" />
      </button>
    </div>
  </nav>
</template>

<style lang="scss" scoped>
@import '../../styles/home/RssListButton.scss';
</style>