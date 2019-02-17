<template>
  <div class="home">
    <Feeds />
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Feeds from '../components/Feeds.vue';

export default Vue.extend({
  name: 'home',
  components: {
    Feeds,
  },
});
</script>
