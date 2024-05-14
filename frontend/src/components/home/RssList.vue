<script lang="ts" setup>
const props = defineProps({
  feedContent: {
    type: Object,
    default: () => ({ feedList: [] })
  }
})
</script>

<template>
  <ul>
    <li v-for="feed in feedContent.feedList" :key="feed.Title" @click="$emit('feed-clicked', feed)">
      <button class="ListItem" title="Read more">
        <img v-if="feed.Image" :src="feed.Image" :alt="`Image of post: ${feed.Title}`" class="Image" />
        <div class="Info">
          <div class="RssInfo">
            <div class="RssTitle">
              <img v-if="feed.FeedImage" :src="feed.FeedImage" :alt="`Image of feed: ${feed.FeedTitle}`"
                class="RssImage" />
              <p>{{ feed.FeedTitle }}</p>
            </div>
            <time>{{ feed.TimeSince }}</time>
          </div>
          <p class="Title" :style="{ color: feed.Readed ? 'gray' : 'black' }">{{ feed.Title }}</p>
        </div>
      </button>
    </li>
  </ul>
</template>

<style lang="scss" scoped>
@import '../../styles/home/RssList.scss'
</style>
