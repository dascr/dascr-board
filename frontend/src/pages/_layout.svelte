<script>
  import { SvelteToast } from '@zerodevx/svelte-toast';
  import Footer from '../components/footer/Footer.svelte';
  import Header from '../components/header/Header.svelte';
</script>

<style>
  main {
    flex: 1;
  }
</style>

<Header />
<main class="flex-grow p-3">
  <slot />
  <SvelteToast />
</main>
<Footer />
