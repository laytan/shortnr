<template>
  <div>
    <label for="password" class="form-label sr-only">Password</label>
    <div class="input-group mb-1">
      <span class="input-group-text bg-primary text-white">
        <Icon icon="lock"/>
      </span>
      <input
        :type="passwordShown ? 'text' : 'password'"
        placeholder="Password"
        class="form-control"
        id="password">
        <button
          type="button"
          class="btn btn-primary btn-sm"
          @click="passwordShown = !passwordShown">
          <Icon :icon="passwordShown ? 'eye-slash' : 'eye'" :size="16"/>
        </button>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';

export default {
  setup() {
    const passwordShown = ref(false);
    return { passwordShown };
  },
};
</script>
