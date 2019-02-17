<template>
    <div>
        <div v-if="this.loading">Loading feeds...</div>
        <Folder v-if="hasFeeds" :name="null" :feeds="feeds.feeds" :contents="feeds.folders" />
        <div v-if="!loading && !hasFeeds">No feeds to show</div>
    </div>
</template>

<script lang="ts">
import Vue from 'vue';

import Folder from './Folder.vue';
import * as newsblur from '../lib/newsblur';

export default Vue.extend({
   components: {Folder},
   data() {
       return {
           loading: false,
           feeds: {feeds: []},
       };
   },

   computed: {
       hasFeeds(): boolean {
           return Object.values(this.feeds.feeds).length > 0;
       },
   },

   async beforeMount() {
       this.loading = true;
       this.feeds = await newsblur.fetchFeeds();
       this.loading = false;
   },
});
</script>
