<script>
  import { Router } from '@roxi/routify';
  import { routes } from '../.routify/routes';
  import Tailwindcss from './Tailwindcss.svelte';
</script>

<style>
  :root {
    --bg-color-1: #1d3c4c;
    --bg-color-2: #3c5b6c;
    --bg-trans-1: rgba(29, 60, 76, 0.9);
    --bg-trans-2: rgba(60, 91, 108);
  }
</style>

<svelte:head>
  <link href="/fonts/fonts_B612Mono.css" rel="stylesheet" />
  <link rel="stylesheet" href="/css/all.min.css" />
</svelte:head>

<!-- <div class=""> -->
<Tailwindcss />
<Router {routes} />
<!-- </div> -->
