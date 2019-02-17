<template>
  <div v-if="loading">Loading...</div>
  <div v-else>
    <h1 class="title">{{story.story_title}}</h1>
    <div v-html="story.story_content"/>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';

import * as newsblur from '../lib/newsblur';

export default Vue.extend({
  data() {
    return {
      loading: false,
      story: null,
    };
  },

  async beforeMount() {
    this.loading = true;
    // TODO: Handle if the story has moved off of the suppled page
    const feed = await newsblur.fetchStories(
      this.$route.params.feed,
      parseInt(this.$route.params.page, 10),
    );
    this.story = feed.stories.find(
      (story: any) => story.id === this.$route.params.story,
    );
    this.loading = false;
  },
});
</script>
