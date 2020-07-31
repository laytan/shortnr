<template>
  <form style="max-width: 600px;" @submit.prevent="onSubmit">
    <h2>
      Create a new link
    </h2>
    <Alert v-if="error.length" type="danger">
      {{ error }}
    </Alert>
    <label for="url" class="sr-only">Link URL</label>
    <input
    v-model="url"
    id="url"
    type="url"
    placeholder="URL"
    class="form-control mb-2"
    autocomplete="url">
    <!-- <label for="custom-id" class="sr-only">Custom ID (Optional)</label>
    <input
    placeholder="Custom ID (Optional)"
    type="text"
    class="form-control mb-2"
    autocomplete="off"
    id="custom-id"> -->
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
    const url = ref('');
    const loading = ref(false);
    const error = ref('');
    const onSubmit = () => {
      error.value = '';
      if (url.value.length < 1) {
        error.value = 'Please fill in all fields';
        return;
      }

      loading.value = true;
      reqP(endpoints.shorten, { url: url.value }, {
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
    };
  },
};
</script>
