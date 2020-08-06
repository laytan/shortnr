<template>
  <li class="list-group-item">
    <div
    data-toggle="collapse"
    :data-target="`#link-${link.id}-collapse`"
    class="d-flex justify-content-between align-items-start">
      <div class="d-flex flex-column mr-2">
        <span>
          <span
            class="text-muted"
            style="font-size: .75rem;">
            {{ redirectBaseUrl }}/</span>{{ link.id }}
        </span>
        <a :href="link.url" rel="noopener noreferrer" target="_BLANK">
          <span class="text-muted">
            {{ link.url }}
          </span>
        </a>
      </div>
      <div class="d-flex flex-column align-items-end justify-content-end text-muted">
        <small class="mb-1">{{ timeAgo(link.created_at) }}</small>
        <Icon icon="caret-down" class="cursor-pointer"/>
      </div>
    </div>
    <div :id="`link-${link.id}-collapse`" class="collapse mt-2">
      <div class="d-flex align-items-center">
        <button
          v-if="canCopy"
          @click="copyLink"
          class="btn btn-primary mr-2">
          {{ copied ? 'copied' : 'copy' }}
        </button>
        <LoadingButton
          @click="deleteLink"
          class="btn-danger text-left mr-2"
          :style="{
            width: loading ? '6rem !important' : 'auto !important',
          }"
          text="remove"
          :loading="loading"/>
          <span v-if="deleteError" class="text-danger mr-2">
            {{ deleteError }}
          </span>
      </div>
    </div>
  </li>
</template>

<script>
import { ref, computed } from 'vue';
import { formatDistanceToNow, parseISO } from 'date-fns';
import { reqD, endpoints } from '@/api';
import { token } from '@/auth';
import LoadingButton from '@/components/forms/LoadingButton.vue';

export default {
  props: {
    link: {
      type: Object,
      required: true,
    },
  },
  components: {
    LoadingButton,
  },
  setup(props, { emit }) {
    const timeAgo = (iso) => `created ${formatDistanceToNow(parseISO(iso))} ago`;

    const redirectBaseUrl = computed(
      () => process.env.VUE_APP_REDIRECT_BASE_URL || window.location.host,
    );

    const copied = ref(false);
    const copyLink = () => {
      navigator.clipboard.writeText(`${redirectBaseUrl.value}/${props.link.id}`)
        .then(() => { copied.value = true; })
        .catch(console.error);
    };

    const loading = ref(false);
    const deleteError = ref('');
    const deleteLink = () => {
      loading.value = true;
      deleteError.value = '';
      reqD(endpoints.deleteLink(props.link.id), {}, {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      })
        .then(() => { emit('removed', props.link); })
        .catch((e) => { deleteError.value = e.message; })
        .finally(() => { loading.value = false; });
    };

    return {
      timeAgo,
      copyLink,
      canCopy: Boolean(navigator.clipboard),
      copied,
      loading,
      deleteLink,
      deleteError,
      redirectBaseUrl,
    };
  },
};
</script>
