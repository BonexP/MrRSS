<script lang="ts" setup>
const props = defineProps({
  selectedFeed: Object
})
</script>

<template>
  <article v-if="selectedFeed">
    <h1>{{ selectedFeed.Title }}</h1>
    <div class="RssInfo">
      <div class="RssTitle">
        <img v-if="selectedFeed.FeedImage" :src="selectedFeed.FeedImage"
          :alt="`Image of feed: ${selectedFeed.FeedTitle}`" class="RssImage" />
        <p>{{ selectedFeed.FeedTitle }}</p>
      </div>
      <time>{{ selectedFeed.Time }}</time>
    </div>
    <div v-html="selectedFeed.Content" class="Content"></div>
  </article>
  <article v-else>
  </article>
</template>

<style lang="scss" scoped>
@import '../../styles/home/RssContent.scss';
</style>