<template>
  <div class="mt-5 d-flex justify-content-center">
    <div>
      <h1 class="h2">Create a free account!</h1>
      <form @submit.prevent="doRegister" class="mb-3">
        <Email @changed="email = $event" class="mb-3"/>
        <Password auto-complete="new-password" @changed="password = $event" class="mb-3"/>
        <LoadingButton text="Register" :loading="loading"/>
      </form>
      <div class="d-flex justify-content-between">
        <p class="m-0">
          Already have an account?
        </p>
        <router-link to="/" class="text-secondary">Log in</router-link>
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

export default {
  name: 'Register',
  components: {
    Email,
    Password,
    LoadingButton,
  },
  setup() {
    const email = ref('');
    const password = ref('');

    const loading = ref(false);

    const success = ref('');
    const error = ref('');

    const doRegister = () => {
      success.value = '';
      error.value = '';

      loading.value = true;
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
