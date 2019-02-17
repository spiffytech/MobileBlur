<template>
    <div v-if="loading">Loading...</div>
    <div v-else>
        <h1 class="title">{{feed.feed_title}}</h1>
        <div
            v-for="story in this.stories"
            :key="story.id"
            :class="{read: story.read_status === 1}"
        >
            <router-link :to="{name: 'story', params: {feed: feed.id, story: story.id, page: story.page}}">
                <header>{{story.story_title}}</header>
                <p>{{story.long_parsed_date}}</p>
            </router-link>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from 'vue';

import * as newsblur from '../lib/newsblur';

export default Vue.extend({
   data() {
       return {
           lastPage: 0,
           loading: false,
           feed: null,
           stories: [] as any[],
       };
   },
   async beforeMount() {
       this.loading = true;
       await this.loadStories();
       this.loading = false;
   },

   methods: {
       async loadStories() {
        const currentPage = this.lastPage + 1;
        const response =
            await newsblur.fetchStories(this.$route.params.feed, currentPage);
        if (!this.feed) this.feed = response;
        const storiesWithPageNumber =
            response.stories.map((story: any) => ({...story, page: currentPage}));
        this.stories = this.stories.concat(storiesWithPageNumber);
        this.lastPage = currentPage;
       },
   },
});
</script>
