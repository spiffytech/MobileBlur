<template>
  <div>
    <div v-if="name" class="folderName">{{name}}</div>
    <div v-for="item in filteredContents" :key="isFolder(item) ? folderName(item) : item">
      <!-- Some items don't have corresponding feeds or folders. No idea what they are. -->
      <div v-if="!isFolder(item) && feeds[item]">
        <router-link :to="{name: 'feed', params: {feed: item}}">{{feeds[item].feed_title}}</router-link>
        {{feeds[item].ng}} {{feeds[item].nt}} {{feeds[item].ps}}
      </div>

      <Folder v-else :name="folderName(item)" :feeds="feeds" :contents="folderContents(item)"/>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';

import * as newsblur from '../lib/newsblur';

export default Vue.extend({
  props: ['contents', 'feeds', 'name'],
  name: 'Folder',

  computed: {
    filteredContents(): any {
      return this.contents.filter(
        (item: any) =>
          this.isFolder(item) ||
          (this.feeds[item] &&
            newsblur.canShowFeed(this.$store.state.filter, this.feeds[item])),
      );
    },
  },
  methods: {
    /**
     * Returns whether or not the given feed-id-or-folder is a folder
     */
    isFolder(item: number | { [key: string]: any }): boolean {
      return typeof item !== 'number';
    },

    /**
     * Folders are represented as an object with a single key pointing to an
     * array of children. Return that single key.
     */
    folderName(folder: { [key: string]: any }) {
      return Object.keys(folder)[0];
    },

    /**
     * Folders are represented as an object with a single key pointing to an
     * array of children. Return that single key's value..
     */
    folderContents(folder: { [key: string]: any }) {
      return Object.values(folder)[0];
    },
  },
});
</script>

<style scoped>
.folderName {
  font-weight: bold;
}
</style>
