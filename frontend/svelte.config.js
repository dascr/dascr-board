import { resolve } from 'path';
import adapter from '@sveltejs/adapter-static';

// /** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    adapter: adapter({
      fallback: 'index.html',
    }),
    vite: {
      resolve: {
        alias: {
          $components: resolve('./src/components'),
          $stores: resolve('./src/utils/stores'),
          $utils: resolve('./src/utils'),
          $styles: resolve('./src/styles'),
        },
      },
    },
  },
  browser: {
    router: true,
  },
};

export default config;
