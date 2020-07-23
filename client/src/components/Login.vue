<template>
  <div class="d-flex justify-content-center">
    <div>
      <form class="mb-3">
        <h2 class="mb-3">Log in to your account</h2>
        <Email class="mb-3"/>
        <Password/>
        <a class="text-right d-block text-secondary mb-2" href="#">Forgot password?</a>
        <button class="btn btn-primary w-100">Log in</button>
      </form>
      <div class="d-flex justify-content-between">
        <p class="m-0">Don't have an account?</p>
        <a href="#" class="text-secondary">Create account</a>
      </div>
    </div>
  </div>
</template>

<script>
import Password from './forms/Password.vue';
import Email from './forms/Email.vue';

export default {
  name: 'Login',
  components: {
    Password,
    Email,
  },
};
</script>
