const APIURL = process.env.VUE_APP_API_URL;

// keep endpoints in one place so we dont have random strings hardcoded everywhere
export const endpoints = {
  register: '/api/v1/users/signup',
  login: '/api/v1/users/login',
  refresh: '/api/v1/users/refresh',
  logout: '/api/v1/users/logout',
};

export const req = (endpoint, opts = {}) => fetch(`${APIURL}${endpoint}`, {
  credentials: 'include',
  ...opts,
})
  .then((res) => res.json())
  .catch(() => { throw new Error('Network error, check your connection'); })
  .then((data) => {
    if (data.err?.msg?.length) {
      throw new Error(data.err.msg);
    }
    return data.res;
  });

export const reqP = (endpoint, body) => req(endpoint, {
  method: 'post',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(body),
});
