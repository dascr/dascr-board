const plugin = require('tailwindcss/plugin');

module.exports = {
  mode: 'jit',
  future: {
    removeDeprecatedGapUtilities: true,
    purgeLayersByDefault: true,
  },
  purge: {
    content: ['./src/**/*.svelte', './public/*.html'],
    css: ['./public/**/*.css'],
  },
  darkMode: false, // or 'media' or 'class'
  theme: {
    gradientColorStops: (theme) => ({
      ...theme('colors'),
      primary: '#1d3c4c',
      secondary: '#3c5b6c',
    }),
    extend: {},
  },
  variants: {
    extend: {},
  },
  plugins: [],
};
