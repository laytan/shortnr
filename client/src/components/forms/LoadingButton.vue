<template>
  <button class="btn btn-primary w-100 position-relative">
    {{ text }}
    <span class="position-absolute">
      <Icon
        :spin="true"
        v-if="loading"
        icon="arrow-clockwise"
        :size="size"
        />
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
    size: {
      type: Number,
      required: false,
    },
  },
};
</script>

<style lang="scss" scoped>
span {
  right: .5rem;
  bottom: .3rem;
}
</style>
