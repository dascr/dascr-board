<script>
  import { SvelteToast } from '@zerodevx/svelte-toast';
  import Footer from '$components/footer/Footer.svelte';
  import Header from '$components/header/Header.svelte';
  import '$styles/tailwind-output.css';
</script>

<svelte:head>
  <link href="/fonts/fonts_B612Mono.css" rel="stylesheet" />
  <link rel="stylesheet" href="/css/all.min.css" />
</svelte:head>

<Header />
<main class="flex-grow p-3">
  <slot />
  <SvelteToast />
</main>
<Footer />

<style>
  main {
    flex: 1;
  }

  :root {
    --bg-color-1: #1d3c4c;
    --bg-color-2: #3c5b6c;
    --bg-trans-1: rgba(29, 60, 76, 0.9);
    --bg-trans-2: rgba(60, 91, 108);
    min-height: 100%;
    background: linear-gradient(90deg, var(--bg-trans-1), var(--bg-trans-2)), url('/img/dart-board.png');
    background-position: center;
    background-size: cover;
    font-family: B612Mono;
  }
</style>
