import ky from 'ky';
let apiURL = 'API_URL';

const api = ky.extend({
  prefixUrl: apiURL,
  credentials: 'include',
});

export default api;
