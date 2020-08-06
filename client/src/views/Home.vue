<template>
  <div class="grid">
    <div class="bg-secondary background"></div>
    <div class="m-3 bg-primary text-light content d-flex align-items-center">
      <div class="container">
        <h1 class="d-inline-block text-secondary font-weight-bold">Shortnr</h1>
        <p class="lead">
          Shortnr is a URL shortener and click tracker.
          It is written in VueJS with the composition API on the front-end and Go on the back-end.
          Complete with JWT Auth with refresh tokens.
          Shortnr is also downloadable as a PWA app for super easy access.
          And it is all FREE!
        </p>
        <div v-if="isLoggedIn">
          Hello {{ user.email }},
          <router-link to="/dashboard" class="link-secondary">go to dashboard</router-link>
        </div>
        <div v-else>
          <router-link to="/register" class="btn btn-secondary mr-2">Register for free</router-link>
          <router-link to="/login" class="btn btn-link text-secondary font-weight-bold">
            Log in
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { isLoggedIn, user } from '@/auth';

export default {
  name: 'Home',
  setup() {
    return {
      isLoggedIn,
      user,
    };
  },
};
</script>

<style lang="scss" scoped>
.grid {
  display: grid;

  .background {
    grid-row: 2 / 10;
    grid-column: 1;
  }

  .content {
    grid-row: 1 / 10;
    grid-column: 1;
  }
}

@media screen and(min-width: 800px) {
  .lead {
    width: 60%;
  }
}
</style>
