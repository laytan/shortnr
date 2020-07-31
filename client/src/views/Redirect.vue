<template>
  <div class="d-flex justify-content-center align-items-center">
    <div class="spinner-border" role="status">
      <span class="sr-only">Loading...</span>
    </div>
  </div>
</template>

<script>
import { reqG, endpoints } from '@/api';
import router from '@/router';

export default {
  name: 'Redirect',
  setup() {
    const { id } = router.currentRoute.value.params;
    reqG(endpoints.link(id))
      .then((res) => { window.location.replace(res.data.url); })
      .catch(() => { router.push({ name: 'NotFound' }); });
  },
};
</script>
