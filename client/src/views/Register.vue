<template>
  <div class="container-fluid mt-5 d-flex justify-content-center">
    <div :style="{ maxWidth: '500px' }">
      <h1 class="h2">Create a free account!</h1>
      <Alert v-if="success.length" type="success">
        {{ success }}
      </Alert>
      <Alert v-if="error.length" type="danger">
        {{ error }}
      </Alert>
      <form @submit.prevent="doRegister" class="mb-3">
        <Email @changed="email = $event" class="mb-3"/>
        <Password auto-complete="new-password" @changed="password = $event" class="mb-3"/>
        <LoadingButton text="Register" :loading="loading"/>
      </form>
      <div class="d-flex justify-content-between">
        <p class="m-0">
          Already have an account?
        </p>
        <router-link to="/" class="link-secondary">Log in</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import { reqP } from '../api';
import Email from '../components/forms/Email.vue';
import Password from '../components/forms/Password.vue';
import LoadingButton from '../components/forms/LoadingButton.vue';
import Alert from '../components/Alert.vue';

export default {
  name: 'Register',
  components: {
    Email,
    Password,
    LoadingButton,
    Alert,
  },
  setup() {
    const email = ref('');
    const password = ref('');

    const loading = ref(false);

    const success = ref('');
    const error = ref('');

    const doRegister = () => {
      // Reset alerts
      success.value = '';
      error.value = '';

      // Validate required fields
      if (!email.value.length || !password.value.length) {
        error.value = 'Please fill in all fields';
        return;
      }

      // Validate password length
      if (password.value.length < 6) {
        error.value = 'Password must be at least 6 characters';
        return;
      }

      // Start showing loading
      loading.value = true;

      // Request signup and handle success and errors setting loading=false at the end
      reqP('signup', {
        email: email.value,
        password: password.value,
      })
        .then(({ msg }) => { success.value = msg; })
        .catch(({ msg }) => { error.value = msg; })
        .finally(() => { loading.value = false; });
    };

    return {
      email, password, doRegister, error, success, loading,
    };
  },
};
</script>
