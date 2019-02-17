<template>
    <div>
        <div v-if="error" class="notification is-danger">
            {{error}}
        </div>

        <form @submit.prevent="submit">
            <div class="field">
                <label class="label">Username</label>
                <div class="control">
                    <input
                        class="input"
                        type="text"
                        v-model="username"
                        @input="clearError"
                        placeholder="username"
                    />
                </div>
            </div>
            <div class="field">
                <label class="label">Password</label>
                <div class="control">
                    <input
                        class="input"
                        type="password"
                        v-model="password"
                        @input="clearError"
                        placeholder="password"
                    />
                </div>
            </div>

            <div class="field">
                <div class="control">
                    <button class="button">Log In</button>
                </div>
            </div>
        </form>
    </div>
</template>

<script lang="ts">
import flatten from 'lodash/flatten';
import Vue from 'vue';

import router from '../router';
import * as newsblur from '../lib/newsblur';

export default Vue.extend({
    data() {
        return {
            username: '',
            password: '',
            error: null as string | null,
        };
    },

    methods: {
        clearError() {
            this.error = null;
        },

        async submit() {
            const response = await newsblur.logIn(this.username, this.password);
            if (!response.authenticated) {
                this.error = flatten(Object.values(response.errors)).join(', ');
                return;
            }

            this.$store.commit('setLoggedIn', true);
            router.push({name: 'home'});
        },
    },
});
</script>
