import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';
import Feed from './components/Feed.vue';
import Story from './components/Story.vue';

import * as newsblur from './lib/newsblur';
import store from './store';

Vue.use(Router);

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
      children: [
        {
          name: 'feed',
          path: 'feed/:feed',
          component: Feed,
        },

        {
          name: 'story',
          path: 'feed/:feed/page/:page/story/:story',
          component: Story,
        },
      ],
    },
    {
      path: '/login',
      name: 'log-in',
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import(/* webpackChunkName: "about" */ './views/LogIn.vue'),
    },
  ],
});

router.beforeEach(async (to, from, next) => {
  if (!store.state.loggedIn && to.name !== 'log-in') {
    // Attempt to validate whether we're logged in
    try {
      await newsblur.fetchFeeds();
      store.commit('setLoggedIn', true);
      return next();
    } catch (ex) {
      return next({name: 'log-in'});
    }
  }
  next();
});

export default router;
