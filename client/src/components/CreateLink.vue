<template>
  <form style="max-width: 600px;" @submit.prevent="onSubmit">
    <h2>
      Create a new link
    </h2>
    <Alert v-if="error.length" type="danger">
      {{ error }}
    </Alert>
    <label for="id" class="sr-only">Custom ID (optional)</label>
    <input
    v-model="id"
    id="id"
    type="text"
    placeholder="Custom ID (optional)"
    class="form-control mb-2"
    autocomplete="off">
    <label for="url" class="sr-only">Link URL (required)</label>
    <input
    v-model="url"
    id="url"
    type="url"
    placeholder="URL (required)"
    class="form-control mb-2"
    autocomplete="url">
    <LoadingButton text="Create" :loading="loading" />
  </form>
</template>

<script>
import { ref } from 'vue';
import { endpoints, reqP } from '@/api';
import { token } from '@/auth';
import LoadingButton from '@/components/forms/LoadingButton.vue';
import Alert from '@/components/Alert.vue';

export default {
  components: {
    LoadingButton,
    Alert,
  },
  setup(props, { emit }) {
    const id = ref('');
    const url = ref('');
    const loading = ref(false);
    const error = ref('');
    const onSubmit = () => {
      error.value = '';
      if (url.value.length < 1) {
        error.value = 'please fill in a URL to shorten';
        return;
      }

      const payload = {
        url: url.value,
      };

      if (id.value.length > 0) {
        payload.id = id.value;
      }

      loading.value = true;
      reqP(endpoints.shorten, payload, {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      })
        .then((res) => { emit('created', res.data); })
        .catch((e) => { error.value = e.message; })
        .finally(() => { loading.value = false; });
    };

    return {
      onSubmit,
      url,
      loading,
      error,
      id,
    };
  },
};
</script>
