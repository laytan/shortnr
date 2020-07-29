import axios from 'axios';

const APIURL = process.env.VUE_APP_API_URL;

// keep endpoints in one place so we dont have random strings hardcoded everywhere
export const endpoints = {
  register: '/api/v1/users/signup',
  login: '/api/v1/users/login',
  refresh: '/api/v1/users/refresh',
  logout: '/api/v1/users/logout',
};

const handleRequest = (request) => axios(request)
  .then((response) => response.data.res)
  .catch((error) => {
    if (error.response) {
      // If rate limited, try again after delay
      if (error.response.status === 429) {
        // eslint-disable-next-line no-console
        console.warn('Request rate-limited');
        return new Promise(
          (resolve, reject) => setTimeout(
            () => handleRequest(request).then(resolve).catch(reject),
            4000,
          ),
        );
        // Error from server
      } if (error.response.data.err?.msg) {
        throw new Error(error.response.data.err.msg);
        // No response from server
      } else {
        throw new Error('something went wrong, please try again');
      }
      // Error constructing request / error with config
    } else {
      throw new Error('something went wrong, please check your connection and try again');
    }
  });

export const reqP = (endpoint, body = {}, opts = {}) => handleRequest({
  method: 'POST',
  url: `${APIURL}${endpoint}`,
  data: body,
  ...opts,
});

export const reqG = (endpoint, body = {}, opts = {}) => handleRequest({
  method: 'GET',
  url: `${APIURL}${endpoint}`,
  params: body,
  ...opts,
});
