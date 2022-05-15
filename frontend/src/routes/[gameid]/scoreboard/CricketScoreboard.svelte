<script>
  import { onMount } from 'svelte';
  import CricketCard from './CricketCard.svelte';
  import { setCricketModeHeader, playSound } from '$utils/methods';
  import state from '$stores/stateStore';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';

  let gameid = $page.params.gameid;

  let mode = '';
  let randomGhost = '';

  onMount(async () => {
    // init websocket
    const im = await import('$utils/socket');
    const ws = im.default;
    const socket = ws.init(gameid, 'Cricket Scoreboard');

    socket.addEventListener('update', async () => {
      await state.updateState(gameid).then(async () => {
        if ($state.gameData.Settings.Sound) {
          await playSound($state.soundToPlay);
        }
      });
    });

    socket.addEventListener('redirect', () => {
      goto(`/${gameid}/start`);
    });

    await state.updateState(gameid).then(async () => {
      if ($state.gameData.Settings.Sound) {
        await playSound($state.soundToPlay);
      }
    });

    const res = setCricketModeHeader($state.gameData);
    mode = res[0];
    randomGhost = res[1];

  });
</script>

<div
  class="flex flex-row mx-auto bg-black bg-opacity-30 rounded-t-2xl overflow-hidden"
>
  <p
    class="text-center border w-1/4 font-bold text-lg rounded-tl-2xl p-2 capitalize"
  >
    Game:
    {$state.gameData.Game}
  </p>
  <p class="text-center border w-1/4 font-bold text-lg p-2">Mode: {mode}</p>
  <p class="text-center border w-1/4 font-bold text-lg p-2">
    Random / Ghost:
    {randomGhost}
  </p>
  <p class="text-center border w-1/4 font-bold text-lg rounded-tr-2xl p-2">
    Round:
    {$state.gameData.ThrowRound}
  </p>
</div>

<div class="bg-black bg-opacity-30 rounded-b-2xl overflow-hidden">
  <p
    class="text-center border w-full font-extrabold text-4xl rounded-b-2xl p-2"
  >
    {$state.message}
  </p>
</div>

<div class="mt-3">
  <div class="flex flex-wrap max-w-full">
    <!-- Player cols -->
    {#each $state.players as player, i}
      <div class="w-72 mx-2 my-2">
        <CricketCard
          {player}
          gameData={$state.gameData}
          active={i === $state.gameData.ActivePlayer}
        />
      </div>
    {/each}
  </div>
</div>
