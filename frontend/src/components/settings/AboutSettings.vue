<script setup lang="ts">
import { Icon } from '@iconify/vue';
</script>

<template>
  <div>
    <h2>MrRss</h2>
    <h4>Versatile, Light-weight, AI-powered RSS Reader</h4>
  </div>
  <img src="../../assets/images/appicon.png" class="logo" alt="MrRss Logo" />
  <div class="version">
    <button class="btn" title="Check update">
      <Icon icon="material-symbols:deployed-code-update-outline" />
    </button>
    <p>
      Current version: 0.1.0
    </p>
  </div>
  <div class="copyright">
    <p>© 2024 <a href="https://github.com/WCY-dt/MrRSS" rel="help" target="_blank">MrRss</a>. All rights reserved.</p>
    <p>Designed and developed by <a href="https://ch3nyang.top/about/" rel="help" target="_blank">Ch3nyang</a>.</p>
  </div>
</template>

<style lang="scss" scoped>
@import '../../styles/settings/AboutSettings.scss';
</style>