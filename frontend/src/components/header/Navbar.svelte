<script>
  import Hamburger from './Hamburger.svelte';

  export let sidebar = false;
</script>

<header
  class="flex justify-between bg-gradient-to-r from-primary to-secondary p-2 items-center text-white border-b-2 pb-3 mb-4">
  <nav class="flex">
    <Hamburger bind:open={sidebar} />
    <a href={'/'}>DaSCR - score and more</a>
  </nav>
</header>
