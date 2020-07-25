<template>
  <button class="btn btn-primary w-100 position-relative">
    {{ text }}
    <span class="position-absolute" v-if="loading">
      <div class="spinner-border spinner-border-sm" role="status">
        <span class="sr-only">Loading...</span>
      </div>
    </span>
  </button>
</template>

<script>
export default {
  props: {
    text: {
      type: String,
      required: true,
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss" scoped>
span.position-absolute {
  right: .5rem;
  bottom: .5rem;
}
</style>
