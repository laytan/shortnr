<template>
  <li class="list-group-item">
    <div @click="toggleCollapse" class="d-flex justify-content-between align-items-start">
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
    <div ref="collapseEl" class="collapse mt-2">
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
          <span class="badge rounded-pill bg-primary">
            Clicks: {{ clicks }}
          </span>
          <span v-if="error" class="text-danger mr-2">
            {{ error }}
          </span>
      </div>
      <div ref="chartEl"></div>
    </div>
  </li>
</template>

<script>
import { ref, watch } from 'vue';
import {
  formatDistanceToNow,
  parseISO,
  isSameDay,
  format,
} from 'date-fns';
import { Collapse } from 'bootstrap/dist/js/bootstrap.bundle';
import { Chart } from 'frappe-charts/dist/frappe-charts.esm';
import 'frappe-charts/dist/frappe-charts.min.css';

import { reqD, endpoints, reqG } from '@/api';
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

    const redirectBaseUrl = process.env.VUE_APP_REDIRECT_BASE_URL || window.location.host;

    const error = ref('');

    /** Handle link copying */
    const copied = ref(false);
    const copyLink = () => {
      error.value = '';
      navigator.clipboard.writeText(`${redirectBaseUrl}/${props.link.id}`)
        .then(() => { copied.value = true; })
        .catch((e) => { error.value = e.message; });
    };

    /** Handle link deleting */
    const loading = ref(false);
    const deleteLink = () => {
      loading.value = true;
      error.value = '';
      reqD(endpoints.deleteLink(props.link.id), {}, {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      })
        .then(() => { emit('removed', props.link); })
        .catch((e) => { error.value = e.message; })
        .finally(() => { loading.value = false; });
    };

    /** Amount of clicks on the link */
    const clicks = ref('loading...');
    const clickElements = ref(null);
    const loadClicks = () => {
      if (clicks.value !== 'loading...') {
        return;
      }

      error.value = '';
      reqG(endpoints.clicks(props.link.id), {}, {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      })
        .then((res) => { clicks.value = res.data.length; clickElements.value = res.data; })
        .catch((e) => { error.value = e.message; });
    };

    /** Show and hide the collapse */
    let collapseShown = false;
    const collapseEl = ref(null);
    const toggleCollapse = () => {
      if (collapseShown) {
        collapseShown = false;
        new Collapse(collapseEl.value).hide();
      } else {
        collapseShown = true;
        new Collapse(collapseEl.value).show();

        // Trigger loading amount of clicks for link when opening collapse
        loadClicks();
      }
    };

    /** Handle chart logic / displaying */
    const chartEl = ref(null);
    watch(clickElements, (value) => {
      if (value.length === 0) {
        return;
      }

      // Change our array with every single click into an array of objects where the same day
      // clicks count up
      const clicksPerDay = value.reduce((aggr, curr) => {
        const date = parseISO(curr.created_at);
        const added = aggr.some((v, i) => {
          if (isSameDay(date, v.date)) {
            // eslint-disable-next-line no-param-reassign
            aggr[i].clicks += 1;
            return true;
          }
          return false;
        });

        if (!added) {
          aggr.push({ date, clicks: 1 });
        }

        return aggr;
      }, []);

      // eslint-disable-next-line no-new
      new Chart(chartEl.value, {
        title: 'Clicks',
        data: {
          labels: clicksPerDay.map((day) => format(day.date, 'd MMM yy')),
          datasets: [
            {
              values: clicksPerDay.map((day) => day.clicks),
            },
          ],
        },
        type: 'line',
        height: 250,
        colors: ['#168967'],
      });
    });

    return {
      timeAgo,
      error,
      copyLink,
      canCopy: Boolean(navigator.clipboard),
      copied,
      loading,
      deleteLink,
      redirectBaseUrl,
      toggleCollapse,
      collapseEl,
      clicks,
      chartEl,
    };
  },
};
</script>
