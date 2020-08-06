<template>
  <div class="view-nav-container vh-100">
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
              <router-link class="nav-link" to="/login">Log in</router-link>
            </li>
            <li class="nav-item" v-if="!isLoggedIn">
              <router-link class="nav-link" to="/register">Register</router-link>
            </li>
            <li class="nav-item" v-if="isLoggedIn">
              <router-link class="nav-link" to="/dashboard">Dashboard</router-link>
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
    <router-view class="view" />
    <Alert
      v-if="updateExists"
      type="secondary"
      heading="Update available"
      class="position-fixed w-75 left-0 right-0 bottom-0 mx-auto my-2">
      <span class="mr-1">
        There is an update available
      </span>
      <a role="button" href="#" @click="refreshApp">
        Update
      </a>
    </Alert>
  </div>
</template>

<script>
import useUpdate from '@/composition/update';
import Alert from '@/components/Alert.vue';
import { isLoggedIn, logout } from './auth';

export default {
  components: {
    Alert,
  },
  setup() {
    const { refreshApp, updateExists } = useUpdate();

    return {
      isLoggedIn,
      logout,
      updateExists,
      refreshApp,
    };
  },
};
</script>

<style lang="scss" scoped>
.view-nav-container {
  display: flex;
  flex-direction: column;

  .view {
    flex-grow: 1;
  }
}
</style>
