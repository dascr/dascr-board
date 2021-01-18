<script>
  export let uid, name, nickname, image, onDelete, showDelete, active;
  let apiBase = 'API_BASE';
  let activeStyle =
    'border-4 border-white shadow-2xl opacity-100 transition duration-700 ease-in-out';
</script>

<div
  class="flex bg-black bg-opacity-50 rounded-lg overflow-hidden shadow-lg opacity-80 {active && activeStyle} playercard">
  <slot name="avatar">
    <div class="flex-none xs: hidden lg:block">
      <img class="" width="300" height="300" src={apiBase + image} alt="" />
    </div>
  </slot>
  <div class="m-6 flex-none w-60 overflow-hidden">
    <p class="font-semibold text-4xl">{name}</p>
    {#if nickname}
      <p class="text-gray-400 text-3xl">{nickname}</p>
    {:else}<br />{/if}
    {#if $$slots.points}
      <slot name="points" />
    {/if}

    {#if showDelete}
      <button
        class="block uppercase border-2 hover:bg-black hover:bg-opacity-30 text-lg  p-4 mt-4 rounded-2xl"
        on:click={onDelete(uid)}><i class="fas fa-trash" />
        Delete</button>
    {/if}
  </div>
  {#if $$slots.score}
    <div class="m-6 w-full">
      <slot name="score" />
    </div>
  {/if}
</div>
