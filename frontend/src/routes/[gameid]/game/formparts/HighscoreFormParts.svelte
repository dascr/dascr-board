<script>
  import { onMount } from 'svelte';
  import setupGame from '$stores/gameStore';

  onMount(() => {
    $setupGame.variant = '301';
  });
</script>

<style>
    input:checked + svg {
        display: block;
    }
    input:disabled + svg {
        display: block;
        background: rgba(0, 0, 0, 0.5);
        cursor: not-allowed;
    }
</style>

<div class="flex flex-col">
  <label for="selectVariant" class="uppercase font-bold text-lg"
    >Which Highscore do you want to play:</label
  >
  <select
    bind:value={$setupGame.variant}
    class="border py-2 px-3 text-gray-900"
    id="selectVariant"
    name="variant"
  >
    <option value="301" selected="selected">301</option>
    <option value="501">501</option>
    <option value="701">701</option>
    <option value="901">901</option>
  </select>
</div>
<div class="flex flex-col">
  <label for="selectInCondition" class="uppercase font-bold text-lg"
    >Condition to get into the game:</label
  >
  <select
    bind:value={$setupGame.in}
    class="border py-2 px-3 text-gray-900"
    id="selectInCondition"
    name="inCondition"
  >
    <option value="straight">Straight</option>
    <option value="double">Double</option>
    <option value="master">Master</option>
  </select>
</div>
<div class="flex flex-col">
  <label for="selectOutCondition" class="uppercase font-bold text-lg"
    >Condition to end the game:</label
  >
  <select
    bind:value={$setupGame.out}
    class="border py-2 px-3 text-gray-900"
    id="selectOutCondition"
    name="outCondition"
  >
    <option value="straight">Straight</option>
    <option value="double" selected="selected">Double</option>
    <option value="master">Master</option>
  </select>
</div>
<div class="flex flex-col">
  <label for="elimination" class="flex justify-start items-start">
    <div
      class="bg-white border-2 rounded border-gray-400 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2"
    >
      <input
        bind:checked={$setupGame.elimination}
        type="checkbox"
        id="elimination"
        name="eliminiation"
        class="opacity-0 absolute"
      />
      <svg
        class="fill-current hidden w-4 h-4 text-gray-900 pointer-events-none"
        viewBox="0 0 20 20"
      >
        <path d="M0 11l2-2 5 5L18 3l2 2L7 18z" />
      </svg>
    </div>
    <div
      class="select-none uppercase font-bold text-lg"
      title="Eliminate player if you match score - sending them to 0"
    >
      Elimination
    </div>
  </label>
</div>
