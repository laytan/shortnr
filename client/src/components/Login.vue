<template>
  <div class="d-flex justify-content-center">
    <div>
      <h1 class="h2">Log in to your account</h1>
      <p v-if="user.email">Hello {{ user.email }}</p>
      <Alert v-if="success.length" type="success">
        {{ success }}
      </Alert>
      <Alert v-if="error.length" type="danger">
        {{ error }}
      </Alert>
      <form @submit.prevent="doLogin" class="mb-3">
        <Email @changed="email = $event" class="mb-3"/>
        <Password @changed="password = $event"/>
        <a class="text-right d-block link-secondary mb-2" href="#">Forgot password?</a>
        <LoadingButton text="Log in" :loading="loading"/>
      </form>
      <div class="d-flex justify-content-between">
        <p class="m-0">Don't have an account?</p>
        <router-link to="/register" class="link-secondary">Create account</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import Password from './forms/Password.vue';
import Email from './forms/Email.vue';
import LoadingButton from './forms/LoadingButton.vue';
import { reqP, endpoints } from '../api';
import Alert from './Alert.vue';
import { login, user } from '../auth';

export default {
  name: 'Login',
  components: {
    Password,
    Email,
    LoadingButton,
    Alert,
  },
  setup() {
    const email = ref('');
    const password = ref('');
    const error = ref('');
    const success = ref('');
    const loading = ref(false);

    const doLogin = () => {
      error.value = '';
      success.value = '';
      loading.value = true;

      reqP(endpoints.login, {
        email: email.value,
        password: password.value,
      })
        .then((res) => {
          login(res.data.token, res.data.refreshToken);
        })
        .catch(({ message }) => { error.value = message; })
        .finally(() => { loading.value = false; });
    };

    return {
      email, password, error, success, loading, doLogin, user,
    };
  },
};
</script>
