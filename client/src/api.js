const APIURL = process.env.VUE_APP_API_URL;

export const req = (endpoint, opts = {}) => fetch(`${APIURL}${endpoint}`, opts)
  .then((res) => {
    if (!res.ok) {
      throw res;
    }
    return res.json().then((data) => data.res);
  })
  .catch((res) => {
    if (!res.json) {
      return {
        msg: 'Could not connect to server, check your connection',
      };
    }
    return res.json().then((data) => data.err);
  });

export const reqP = (endpoint, body) => req(endpoint, {
  method: 'post',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(body),
});
