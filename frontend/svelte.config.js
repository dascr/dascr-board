import { resolve } from 'path';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    // hydrate the <div id="svelte"> element in src/app.html
    target: '#svelte',
    ssr: false,
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
    router: true,
  },
};

export default config;
