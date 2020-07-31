<template>
  <div>
    <Alert class="w-md-25 m-3 position-absolute right-0" v-if="showLoggedInAlert">
      Welcome back {{ user.email }}!
    </Alert>
    <div class="container-fluid mt-5">
      <div class="row">
        <div class="col-md-6 pr-md-4">
          <CreateLink style="max-width: 600px;" @created="onCreatedLink"/>
        </div>
        <div class="col-md-6 pl-md-4">
          <div class="d-flex align-items-center">
            <h2>
              Your links
            </h2>
            <div class="spinner-border spinner-border-sm ml-2" role="status" v-if="loadingLinks">
              <span class="sr-only">Loading...</span>
            </div>
          </div>
          <ul class="list-group">
            <LinkListItem v-for="link in sortedLinks" :link="link" :key="link.id" />
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import { parseISO } from 'date-fns';

import { user, token } from '@/auth';
import { reqG, endpoints } from '@/api';
import router from '@/router';

import Alert from '@/components/Alert.vue';
import CreateLink from '@/components/CreateLink.vue';
import LinkListItem from '@/components/LinkListItem.vue';

export default {
  name: 'Dashboard',
  components: {
    Alert,
    CreateLink,
    LinkListItem,
  },
  setup() {
    // Show welcome message if we just logged in
    const showLoggedInAlert = ref(false);
    if (router.currentRoute.value.query['from-login'] === 'true') {
      showLoggedInAlert.value = true;
      setTimeout(() => { showLoggedInAlert.value = false; }, 3000);
    }

    const loadingLinks = ref(true);
    const links = ref([]);

    reqG(endpoints.links, {}, {
      headers: {
        Authorization: `Bearer ${token.value}`,
      },
    })
      .then((res) => { res.data.forEach(onCreatedLink); })
      .catch((e) => { console.error(e); })
      .finally(() => { loadingLinks.value = false; });

    const onCreatedLink = (link) => {
      links.value.push(link);
    };

    // Sort on time FIXME:does not work
    const sortedLinks = computed(
      () => links.value.sort(
        (a, b) => parseISO(b.created_at).getTime() - parseISO(a.created_at).getTime(),
      ),
    );

    return {
      showLoggedInAlert,
      user,
      onCreatedLink,
      sortedLinks,
      loadingLinks,
    };
  },
};
</script>
