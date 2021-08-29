import { resolve } from 'path';
import adapter from '@sveltejs/adapter-static';

// /** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    adapter: adapter({
      // default options are shown
      pages: 'build',
      assets: 'build',
      fallback: null,
    }),
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
