import ky from 'ky';
let apiURL = import.meta.env.VITE_API_URL;

const api = ky.extend({
  prefixUrl: apiURL,
  credentials: 'include',
});

export default api;
