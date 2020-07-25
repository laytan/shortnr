<template>
  <div
    :class="`alert alert-${type}`"
    role="alert">
    <p class="sr-only">{{ type }} alert</p>
    <h4 v-if="heading.length">heading</h4>
    <slot></slot>
  </div>
</template>

<script>
export default {
  props: {
    type: {
      type: String,
      default: 'info',
    },
    heading: {
      type: String,
      default: '',
    },
  },
};
</script>
