<script setup lang="ts">
import { useRouter } from 'vue-router';
import { computed } from 'vue';
import { Icon } from '@iconify/vue';

const router = useRouter();

const goBack = () => {
  router.push('/');
}

const pageTitle = computed(() => {
  switch (router.currentRoute.value.path) {
    case '/settings/rss':
      return 'RSS Settings';
    case '/settings/preference':
      return 'Preferences';
    case '/settings/about':
      return 'About';
    default:
      return 'Settings';
  }
});
</script>

<template>
  <div class="settings">
    <div class="nav">
      <div class="title">
        <button @click="goBack" class="btn done" title="Back to home">
          <Icon icon="material-symbols:arrow-back" />
        </button>
        <h1>{{ pageTitle }}</h1>
      </div>
      <div class="router">
        <router-link to="/settings/rss" class="btn" title="RSS settings">
          <Icon icon="material-symbols:rss-feed" />
        </router-link>
        <router-link to="/settings/preference" class="btn" title="Preference settings">
          <Icon icon="material-symbols:inbox-customize" />
        </router-link>
        <router-link to="/settings/about" class="btn" title="About MrRSS">
          <Icon icon="material-symbols:page-info" />
        </router-link>
      </div>
    </div>
    <router-view />
  </div>
</template>

<style lang="scss" scoped>
@import '../../styles/settings/Settings.scss';
</style>