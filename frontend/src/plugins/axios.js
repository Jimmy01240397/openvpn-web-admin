import axios from 'axios';
import router from '@/router';

const maxRetryTime = 3;
const retryDelay = 1000;

const option = {
  baseURL:         window.location.origin,
  headers:         { 'Content-Type': 'application/json' },
  withCredentials: true,
};

const instance = axios.create(option);

// set up the default timeout threshold
instance.defaults.timeout = 10000;

instance.interceptors.request.use(request => {
  return request;
});

instance.interceptors.response.use(response => {
  return response;
}, async error => {
  const { config } = error;
  // reject if there is no config field
  if (!config) {
    return Promise.reject(error);
  }
  if (!error.response) {
    error.response = { data: 'Network Error' };
  }
  if (error.response.status === 401) {
    router.push('/');
  }
  // reject if the error is not caused at server side
  if (error.response.status !== 500) {
    return Promise.reject(error);
  }

  config.headers['x-retry-count'] = config.headers['x-retry-count'] || 0;
  if (config.headers['x-retry-count'] >= maxRetryTime) {
    return Promise.reject(error);
  }

  config.headers['x-retry-count'] += 1;

  // create a new promise to handle exponential backoff
  const backoffDelay = (1 / 2) * Math.pow(2, config.headers['x-retry-count']) * retryDelay + Math.floor(Math.random() * 150 + 50);
  await new Promise((resolve) => {
    setTimeout(resolve, backoffDelay);
  });
  return instance(config);
})

export default instance;
