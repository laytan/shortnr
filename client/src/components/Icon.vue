<template>
  <svg class="bi" :width="size" :height="size" fill="currentColor">
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
  },
};
</script>
