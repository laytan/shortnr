<template>
  <li class="list-group-item d-flex justify-content-between align-items-start">
    <div class="d-flex flex-column">
      <span>
        {{ link.id }}
      </span>
      <a :href="link.url" rel="noopener noreferrer" target="_BLANK">
        <span class="text-muted">
          {{ link.url }}
        </span>
      </a>
    </div>
    <div class="d-flex flex-column align-items-end">
      <!-- <span class="badge bg-primary d-flex align-items-center">
        <Icon icon="cursor" :size="16"/>
        14 hits
      </span> -->
      <small class="text-muted">{{ timeAgo(link.created_at) }}</small>
      <button
        v-if="canCopy"
        @click="copyLink(link.id)"
        class="btn btn-primary btn-sm mt-1">
        copy
        <Icon v-if="copied" icon="check" :size="16"/>
      </button>
      <!-- <small><a role="button" class="text-danger d-flex align-items-center">
        <Icon icon="trash" :size="16"/>
        Remove
      </a></small> -->
    </div>
  </li>
</template>

<script>
import { ref } from 'vue';
import { formatDistanceToNow, parseISO } from 'date-fns';

export default {
  props: {
    link: {
      type: Object,
      required: true,
    },
  },
  setup() {
    const timeAgo = (iso) => `created ${formatDistanceToNow(parseISO(iso))} ago`;

    const copied = ref(false);
    const copyLink = (id) => {
      navigator.clipboard.writeText(`${window.location.host}/${id}`)
        .then(() => { copied.value = true; })
        .catch(console.error);
    };

    return {
      timeAgo,
      copyLink,
      canCopy: Boolean(navigator.clipboard),
      copied,
    };
  },
};
</script>
