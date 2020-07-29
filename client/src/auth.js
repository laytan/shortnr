import { ref, computed } from 'vue';
import { reqP, endpoints, reqG } from './api';

// https://stackoverflow.com/questions/38552003/how-to-decode-jwt-token-in-javascript-without-using-a-library
const parseJwt = (jwtToken) => {
  const base64Url = jwtToken.split('.')[1];
  const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
  const jsonPayload = decodeURIComponent(atob(base64).split('').map((c) => `%${(`00${c.charCodeAt(0).toString(16)}`).slice(-2)}`).join(''));

  return JSON.parse(jsonPayload);
};

const getCurrentUnixTimestamp = () => Math.round((new Date()).getTime() / 1000);

const isUserDifferent = (user1, user2) => user1?.id !== user2?.id;

// try to refresh the token
export const doRefresh = () => {
  // console.log('refreshing token');
  reqP(endpoints.refresh)
    .then((res) => { login(res.data.token); })
    .catch(() => { logout(); });
};

export const token = ref(null);
export const user = ref({});

// remove token from memory
// let server remove the http-only refresh token (can't be done client-side)
export const logout = () => {
  if (token.value?.length) {
    // console.log('removing refresh cookie');
    reqG(endpoints.logout, {}, {
      method: 'DELETE',
    });
  }
  // console.log('logging out');
  token.value = null;
  user.value = {};
};

export const login = (loginToken) => {
  // console.log('logging in');
  token.value = loginToken;

  const payload = parseJwt(loginToken);
  const { exp } = payload;

  const payloadUser = {
    email: payload.email,
    id: payload.id,
    createdAt: payload.createdAt,
    updatedAt: payload.updatedAt,
  };

  if (isUserDifferent(user.value, payloadUser)) {
    user.value = payloadUser;
  }

  // exp is a unix timestamp of when the token expires
  const diff = exp - getCurrentUnixTimestamp();
  // console.log(`token expires in ${diff} seconds`);

  // Try to get a new token from the server 30 seconds before the token expires
  setTimeout(doRefresh, (diff - 30) * 1000);
};

// helper for checking login status
export const isLoggedIn = computed(() => token.value && !!token.value.length);
