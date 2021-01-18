<script>
  import Navbar from './Navbar.svelte'
  import Sidebar from './Sidebar.svelte'
  let open = false
</script>

<div class="header">
<Sidebar bind:open/>
<Navbar bind:sidebar={open} />
</div>
