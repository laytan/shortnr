<template>
  <div class="d-flex justify-content-center">
    <div>
      <h1 class="h2">Log in to your account</h1>
      <form @submit="login" class="mb-3">
        <Email @change="creds.email = $event" class="mb-3"/>
        <Password @change="creds.password = $event"/>
        <a class="text-right d-block link-secondary mb-2" href="#">Forgot password?</a>
        <button class="btn btn-primary w-100">Log in</button>
      </form>
      <div class="d-flex justify-content-between">
        <p class="m-0">Don't have an account?</p>
        <router-link to="/register" class="link-secondary">Create account</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { reactive } from 'vue';
import Password from './forms/Password.vue';
import Email from './forms/Email.vue';

export default {
  name: 'Login',
  components: {
    Password,
    Email,
  },
  setup() {
    const creds = reactive({
      email: '',
      password: '',
    });

    const login = () => {

    };

    return { creds, login };
  },
};
</script>
