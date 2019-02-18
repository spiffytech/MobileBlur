<template>
  <div id="app">
    <div id="nav">
      <router-link to="/">Home</router-link> |
      <router-link to="/about">About</router-link>
      <button @click="logout">Log Out</button>
    </div>
    <router-view/>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';

import * as newsblur from './lib/newsblur';

export default Vue.extend({
  methods: {
    async logout() {
      await newsblur.logOut();
    }
  }
});
</script>


<style></style>
