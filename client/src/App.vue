<template>
  <div>
    <nav class="navbar navbar-expand-md navbar-dark bg-secondary">
      <div class="container-fluid">
        <router-link class="navbar-brand text-primary" to="/">Shortnr</router-link>
        <button
          class="navbar-toggler"
          data-toggle="collapse"
          data-target="#primaryNav"
          aria-controls="primaryNav"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="primaryNav">
          <ul class="navbar-nav ml-auto mb-2 mb-md-0">
            <li class="nav-item" v-if="!isLoggedIn">
              <router-link class="nav-link" to="/">Log in</router-link>
            </li>
            <li class="nav-item" v-if="!isLoggedIn">
              <router-link class="nav-link" to="/register">Register</router-link>
            </li>
            <li class="nav-item" v-if="isLoggedIn">
              <button
              class="nav-link bg-transparent border-0"
              style="font-weight: 600;"
              @click="logout">
                Log out
              </button>
            </li>
          </ul>
        </div>
      </div>
    </nav>
    <router-view/>
  </div>
</template>

<script>
import { doRefresh, isLoggedIn, logout } from './auth';

export default {
  setup() {
    doRefresh();

    return {
      isLoggedIn,
      logout,
    };
  },
};
</script>
