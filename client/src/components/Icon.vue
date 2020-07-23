<template>
  <svg :class="`bi ${spin ? 'spin' : ''}`" :width="size" :height="size" fill="currentColor">
    <use :xlink:href="`bootstrap-icons.svg#${icon}`" />
  </svg>
</template>

<script>
export default {
  props: {
    size: {
      type: Number,
      default: 24,
    },
    icon: {
      type: String,
      required: true,
    },
    spin: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss" scoped>
.spin {
  animation: spin 1s infinite;
}

@keyframes spin {
  0% {
    transform: rotateZ(0deg);
  }
  100% {
    transform: rotateZ(360deg);
  }
}
</style>
