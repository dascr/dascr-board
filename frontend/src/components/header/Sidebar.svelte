<script>
  export let open = false;
</script>

<style>
  aside {
    left: -100%;
    transition: left;
  }

  .open {
    left: 0%;
    width: auto;
  }
</style>

<aside
  class="absolute w-full h-full bg-gradient-to-r from-primary to-secondary border-r-2 shadow-lg"
  class:open>
  <nav class="p-12 text-xl">
    <a
      class="block"
      href={'/'}
      on:click={() => {
        open = false;
      }}><i class="fas fa-home w-10" />Start</a>
    <a
      class="block"
      href={'/player'}
      on:click={() => {
        open = false;
      }}><i class="fas fa-user w-10" />Player</a>
    <a
      class="block"
      href={'/credits'}
      on:click={() => {
        open = false;
      }}><i class="fas fa-info-circle w-10" />Credits</a>
  </nav>
</aside>
