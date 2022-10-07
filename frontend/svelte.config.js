import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    adapter: adapter({
      fallback: 'index.html',
    }),
    alias: {
      $components: './src/components',
      $stores: './src/utils/stores',
      $utils: './src/utils',
      $styles: './src/styles',
    },
  },
};

export default config;
