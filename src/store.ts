import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    loggedIn: false,
    filter: 'nt',
  },
  mutations: {
    setLoggedIn(state, isLoggedIn) {
      state.loggedIn = isLoggedIn;
    },
  },
  actions: {

  },
});
