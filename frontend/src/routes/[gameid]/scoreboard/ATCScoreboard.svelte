<script>
  import { onMount } from 'svelte';
  import PlayerCard from '../../player/PlayerCard.svelte';
  import { scoreOrCurrentNumber, playSound } from '$utils/methods';
  import state from '$stores/stateStore';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';

  let gameid = $page.params.gameid;

  onMount(async () => {
    // init websocket
    const im = await import('$utils/socket');
    const ws = im.default;
    const socket = ws.init(gameid, 'ATC Scoreboard');

    await state.updateState(gameid).then(async () => {
      if ($state.gameData.Settings.Sound) {
        await playSound($state.soundToPlay);
      }
    });

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
  });
</script>

<div
  class="flex flex-row mx-auto bg-black bg-opacity-30 rounded-t-2xl overflow-hidden"
>
  <p class="text-center border w-1/3 font-bold text-lg rounded-tl-2xl p-2">
    Game: Around The Clock
  </p>
  <p class="text-center border w-1/3 font-bold text-lg p-2">
    Variant:
    {$state.gameData.Variant}
  </p>
  <p class="text-center border w-1/3 font-bold text-lg rounded-tr-2xl p-2">
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

<div class="max-w-full space-y-2">
  <div class="flex flex-wrap">
    {#each $state.players as player, i}
      <div class="w-full p-2 2xl:w-1/2">
        <PlayerCard
          uid={player.UID}
          name={player.Name}
          nickname={player.Nickname}
          image={player.Image}
          showDelete={false}
          onDelete={() => {}}
          active={i === $state.gameData.ActivePlayer}
        >
          <div slot="points">
            <p class="font-extrabold text-5xl mt-5">
              {scoreOrCurrentNumber(player, $state.gameData)}
            </p>
            <p class="font-semibold  text-2xl mt-5 flex flex-row">
              <img
                src="/img/dart.png"
                width="32"
                height="32"
                alt=""
                class="mr-4"
              />
              {player.TotalThrowCount}
            </p>
          </div>
          <div slot="score" class="h-full">
            <table class="w-full h-full">
              <tr class="">
                {#each player.LastThrows as thr, i}
                  <td
                    class="p-2 text-4xl font-extrabold text-center w-1/4 border-dashed border-white border-opacity-10"
                    class:border-l={i != 0}
                  >
                    {#if thr.Modifier === 2}
                      D{thr.Number}
                    {:else if thr.Modifier === 3}
                      T{thr.Number}
                    {:else}{thr.Number}{/if}
                  </td>
                {/each}
                {#each Array(3 - player.LastThrows.length) as _, __}
                  <td
                    class="p-2 text-4xl font-extrabold text-center w-1/4 border-dashed border-white border-opacity-10"
                    class:border-l={i != 0}
                  >
                    -
                  </td>
                {/each}
              </tr>
            </table>
          </div>
        </PlayerCard>
      </div>
    {/each}
  </div>
</div>
