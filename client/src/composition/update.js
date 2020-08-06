import { ref } from 'vue';

export default function useUpdate() {
  const refreshing = ref(false);
  const registration = ref(null);

  // Used to check if an update is needed in other components
  const updateExists = ref(false);

  // Event from service worker that an update is available
  const updateAvailable = (e) => {
    registration.value = e.detail;
    updateExists.value = true;
  };
  document.addEventListener('swUpdated', updateAvailable, { once: true });

  // Prevent multiple refreshes
  navigator.serviceWorker.addEventListener('controllerchange', () => {
    if (refreshing.value === true) {
      return;
    }
    refreshing.value = true;

    window.location.reload();
  });

  // Invoked on accepting the update
  const refreshApp = () => {
    updateExists.value = false;

    if (!registration.value || !registration.value.waiting) {
      return;
    }

    registration.value.waiting.postMessage({ type: 'SKIP_WAITING' });
  };

  return {
    updateExists,
    refreshApp,
  };
}
