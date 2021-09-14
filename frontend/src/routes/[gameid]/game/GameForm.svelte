<script>
  import { onMount } from 'svelte';
  import api from '$utils/api';
  import X01FormParts from './formparts/X01FormParts.svelte';
  import CricketFormParts from './formparts/CricketFormParts.svelte';
  import ATCFormParts from './formparts/ATCFormParts.svelte';
  import SplitScoreFormParts from './formparts/SplitScoreFormParts.svelte';
  import ShanghaiFormParts from './formparts/ShanghaiFormParts.svelte';
  import EliminationFormParts from './formparts/EliminationFormParts.svelte';
  import setupGame from '$stores/gameStore';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import DragDrop from 'svelte-dragdroplist';
  import AtcController from '../controller/ATCController.svelte';
  import AtcFormParts from './formparts/ATCFormParts.svelte';
  import AtcRules from './rules/ATCRules.svelte';

  let gameID = $page.params.gameid;
  let games = [
    { id: 'x01', text: 'X01' },
    { id: 'cricket', text: 'Cricket' },
    { id: 'atc', text: 'Around The Clock' },
    { id: 'split', text: 'Split-Score' },
    { id: 'shanghai', text: 'Shanghai' },
    { id: 'elim', text: 'Elimination' },
  ];

  let selectedPlayer = [];

  let punisherAvailable = ['x01', 'elim'];

  let availablePlayer = [];
  let open;
  let gameMode = X01FormParts;

  const updateSelected = () => {
    selectedPlayer = [];
    $setupGame.player.map((pl) => {
      console.log(pl);
      pl.Nickname
        ? selectedPlayer.push({
            UID: pl.UID,
            text: pl.Name + ' - ' + pl.Nickname,
          })
        : selectedPlayer.push({ UID: pl.UID, text: pl.Name });
    });
    console.log(selectedPlayer);
  };

  $: {
    let selectedGame = $setupGame.game;

    switch (selectedGame) {
      case 'x01':
        gameMode = X01FormParts;
        break;
      case 'cricket':
        gameMode = CricketFormParts;
        break;
      case 'atc':
        gameMode = ATCFormParts;
        break;
      case 'split':
        gameMode = SplitScoreFormParts;
        break;
      case 'shanghai':
        gameMode = ShanghaiFormParts;
        break;
      case 'elim':
        gameMode = EliminationFormParts;
        break;

      default:
        break;
    }
  }

  onMount(async () => {
    open = false;
    const res = await api.get('player');
    availablePlayer = await res.json();
  });

  async function handleSubmit() {
    let playingPlayer = [];
    let playerInt = $setupGame.player.map((el) => {
      playingPlayer.push(parseInt(el.UID));
    });

    await api
      .post(`game/${gameID}`, {
        json: {
          uid: gameID,
          player: playerInt || [],
          game: $setupGame.game || '',
          variant: $setupGame.variant || '',
          in: $setupGame.in || '',
          out: $setupGame.out || '',
          punisher: $setupGame.settings.punisher || false,
          sound: $setupGame.settings.sound || false,
          podium: $setupGame.settings.podium || false,
          autoswitch: $setupGame.settings.autoswitch || false,
          cricketrandom: $setupGame.cricket.random || false,
          cricketghost: $setupGame.cricket.ghost || false,
        },
      })
      .then(async (res) => {
        if (res.status === 200) {
          goto(`/${gameID}/controller`);
        }
      })
      .catch((err) => console.error(err));
  }
</script>

<h3 class="text-center font-bold underline text-4xl">Game Setup</h3>
<form
  id="newGame"
  enctype="multipart/form-data"
  autocomplete="off"
  on:submit|preventDefault={handleSubmit}
>
  <div class="space-y-4">
    <div class="flex flex-col">
      <label for="selectPlayer" class="uppercase font-bold text-lg"
        >Select Players:</label
      >
      <div class="flex flex-row">
        <div>
          <select
            bind:value={$setupGame.player}
            on:change={updateSelected}
            class="border py-2 px-3 text-gray-900"
            id="selectPlayer"
            name="player"
            multiple
            required
          >
            {#each availablePlayer || [] as player}
              <option value={player}>{player.Name}</option>
            {/each}
          </select>
        </div>
        <div class="py-2 px-3 text-gray-900">
          <DragDrop bind:data={selectedPlayer} removeItems={false} />
        </div>
      </div>
    </div>

    <div class="flex flex-col">
      <label for="selectGame" class="uppercase font-bold text-lg"
        >Choose your game type:</label
      >
      <select
        bind:value={$setupGame.game}
        id="selectGame"
        class="border py-2 px-3 text-gray-900"
        name="game"
        required
      >
        {#each games as game}
          <option value={game.id}>{game.text}</option>
        {/each}
      </select>
    </div>

    <!-- dynamic form parts -->
    <svelte:component this={gameMode} />

    <button
      type="submit"
      class="block uppercase border-2 hover:bg-black hover:bg-opacity-30 text-lg mx-auto p-4 rounded-2xl"
      ><i class="fas fa-play pr-2" />
      Start Game</button
    >

    <!-- Further settings -->

    <div class="flex flex-col">
      <button
        class="block uppercase border-2 hover:bg-black hover:bg-opacity-30 text-lg mx-auto p-2 rounded-2xl"
        type="button"
        on:click={() => (open = !open)}
        ><i class="fas fa-sliders-h pr-2" />Further Settings</button
      >
    </div>

    <div class:open class="hidden">
      <div class="flex flex-wrap">
        <div class="w-1/3">
          <label for="sound" class="flex justify-start items-start">
            <div
              class="bg-white border-2 rounded border-gray-400 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2"
            >
              <input
                bind:checked={$setupGame.settings.sound}
                type="checkbox"
                id="sound"
                name="sound"
                class="opacity-0 absolute"
              />
              <svg
                class="fill-current hidden w-4 h-4 text-gray-900 pointer-events-none"
                viewBox="0 0 20 20"
                ><path d="M0 11l2-2 5 5L18 3l2 2L7 18z" /></svg
              >
            </div>
            <div
              class="select-none uppercase font-bold text-lg"
              title="Play sound"
            >
              Sound
            </div>
          </label>
        </div>

        <div
          class="w-1/3"
          class:hidden={$setupGame.game === 'shanghai' ||
            $setupGame.game === 'split'}
        >
          <label for="podium" class="flex justify-start items-start">
            <div
              class="bg-white border-2 rounded border-gray-400 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2"
            >
              <input
                bind:checked={$setupGame.settings.podium}
                type="checkbox"
                id="podium"
                name="podium"
                class="opacity-0 absolute"
              />
              <svg
                class="fill-current hidden w-4 h-4 text-gray-900 pointer-events-none"
                viewBox="0 0 20 20"
                ><path d="M0 11l2-2 5 5L18 3l2 2L7 18z" /></svg
              >
            </div>
            <div
              class="select-none uppercase font-bold text-lg"
              title="Continue game after a player has won until one player is left"
            >
              podium
            </div>
          </label>
        </div>

        <div class="w-1/3">
          <label for="autoswtich" class="flex justify-start items-start">
            <div
              class="bg-white border-2 rounded border-gray-400 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2"
            >
              <input
                bind:checked={$setupGame.settings.autoswitch}
                type="checkbox"
                id="autoswitch"
                name="autoswitch"
                class="opacity-0 absolute"
              />
              <svg
                class="fill-current hidden w-4 h-4 text-gray-900 pointer-events-none"
                viewBox="0 0 20 20"
                ><path d="M0 11l2-2 5 5L18 3l2 2L7 18z" /></svg
              >
            </div>
            <div
              class="select-none uppercase font-bold text-lg"
              title="When not using an automated recognition (cam / machine) this will switch to next player automatically when 3 throws are inserted after a delay of 2 seconds"
            >
              Auto Switch next player
            </div>
          </label>
        </div>

        <div
          class="w-1/3"
          class:hidden={!punisherAvailable.includes($setupGame.game)}
        >
          <label for="punisher" class="flex justify-start items-start">
            <div
              class="bg-white border-2 rounded border-gray-400 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2"
            >
              <input
                bind:checked={$setupGame.settings.punisher}
                type="checkbox"
                id="punisher"
                name="punisher"
                class="opacity-0 absolute"
              />
              <svg
                class="fill-current hidden w-4 h-4 text-gray-900 pointer-events-none"
                viewBox="0 0 20 20"
                ><path d="M0 11l2-2 5 5L18 3l2 2L7 18z" /></svg
              >
            </div>
            <div
              class="select-none uppercase font-bold text-lg"
              title="When checked you will gain 100 when missing the numbers"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                x="0px"
                y="0px"
                width="20"
                height="20"
                viewBox="0 0 226 226"
                style=" fill:#000000;"
                ><g
                  fill="none"
                  fill-rule="nonzero"
                  stroke="none"
                  stroke-width="1"
                  stroke-linecap="butt"
                  stroke-linejoin="miter"
                  stroke-miterlimit="10"
                  stroke-dasharray=""
                  stroke-dashoffset="0"
                  font-family="none"
                  font-weight="none"
                  font-size="none"
                  text-anchor="none"
                  style="mix-blend-mode: normal"
                  ><path d="M0,226v-226h226v226z" fill="none" /><g
                    fill="#ffffff"
                    ><path
                      d="M110.01202,0.54327c-9.50721,0.20373 -17.55438,1.69772 -24.71875,2.98798c-6.96063,1.73167 -13.20823,3.32753 -19.28606,6.79087c-3.46334,2.61448 -6.79087,5.36479 -11.13702,10.59375c-4.34615,4.34615 -7.97927,8.6244 -10.59375,9.50721c-7.80949,5.22897 -12.9366,12.22356 -14.66827,20.91587c-3.46334,13.92128 -2.78425,32.05288 2.44471,53.78365c0.88282,2.61448 2.81821,6.04387 5.43269,9.50721c3.46335,4.34615 5.05919,6.96063 6.79087,8.69231c-0.88281,0 -0.71304,1.08654 -2.44471,1.08654l-2.71635,0.8149c-2.61448,0.88282 -3.53125,2.54658 -3.53125,5.16106l0.8149,3.53125l1.08654,2.44471l-1.90144,1.08654v1.62981c0.88282,4.34615 1.73167,6.89273 4.34615,9.50721c0.88282,2.61448 2.64844,4.34615 3.53125,4.34615h3.53125c0.88282,-0.88281 2.44471,-1.62981 2.44471,-1.62981c1.73167,-1.73167 3.70102,-1.69771 5.43269,-0.8149l2.44471,0.8149l3.53125,2.44471c1.73167,0 2.61448,-0.74699 4.34615,-1.62981l1.62981,-1.62981l-2.44471,-1.08654h-2.71635c-2.61448,0 -3.53125,-0.74699 -3.53125,-1.62981l5.43269,-4.34615l5.16106,-4.34615c3.46335,-1.73167 7.8774,-3.53125 12.22356,-3.53125c2.61448,0.88282 3.25962,1.79958 3.25962,3.53125c-1.73167,0.88282 -2.51262,2.54658 -1.62981,5.16106l0.8149,3.53125c0.88282,0.88282 0.8149,1.83354 0.8149,2.71635c-2.61448,-0.88281 -2.51262,0.747 -1.62981,5.97596c0.88282,2.61448 1.79958,6.04387 3.53125,9.50721c0,9.57512 0.71304,27.8086 2.44471,52.15385l1.90144,1.90144l2.44471,1.62981c2.61448,1.73167 4.54988,1.73167 5.43269,0c0.88282,-1.73167 0.8149,-12.12169 0.8149,-31.23798c0,-19.11629 -0.8149,-28.6914 -0.8149,-30.42308l-1.90144,-1.90144c1.73167,-0.88281 3.42938,-0.71304 5.16106,-2.44471c-0.88281,8.69231 0.98467,32.22265 2.71635,68.72356c0,0.88282 0.91677,0.8149 3.53125,0.8149c2.61448,0 3.53125,-0.0679 3.53125,0.8149c1.73167,-0.88281 2.37681,-11.23888 3.25962,-31.23798v-23.63221c0,-6.96063 0.06791,-12.05379 -0.8149,-14.66827h3.53125c-0.88281,11.30679 -0.8149,25.32993 -0.8149,41.83173c-0.88281,18.26743 -0.0679,27.77464 0.8149,26.89183h4.34615c2.61448,0.88282 3.46335,-8.79417 4.34615,-28.79327c0.88282,-7.80949 0.8149,-15.48317 0.8149,-24.17548c0,-7.80949 0.06791,-14.125 -0.8149,-18.47115l4.34615,-0.8149c3.46335,1.73167 5.16106,3.53125 5.16106,3.53125v0.8149l-3.53125,0.8149l-1.62981,1.08654c0,0.88282 0.06791,4.17638 -0.8149,11.13702l-0.8149,23.63221c-0.88281,21.73077 -0.0679,32.1208 0.8149,31.23798l9.50721,-1.90144c0.88282,-0.88281 1.83354,-7.70763 2.71635,-23.36058c0.88282,-13.03846 0.8149,-24.3792 0.8149,-33.95433c2.61448,-2.61448 4.24429,-6.89273 5.97596,-13.85337l0.8149,-7.0625c0,-1.73167 0.06791,-3.25962 -0.8149,-3.25962l-0.8149,-1.08654l0.8149,-1.62981l5.16106,-1.62981l5.43269,-1.90144c4.34615,-0.88281 7.70763,0.91677 10.32212,3.53125c0.88282,1.73167 1.69772,3.36148 0.8149,5.97596l-0.8149,1.90144c0,0.88282 0.06791,0.8149 -0.8149,0.8149c-0.88281,2.61448 -0.16977,4.34615 2.44471,4.34615c2.61448,0.88282 4.44802,0.91677 7.0625,-0.8149c0.88282,-1.73167 3.59916,-3.42938 7.0625,-5.16106c2.61448,-1.73167 4.27825,-4.51592 5.16106,-6.2476l0.8149,-5.97596c0,-2.61448 -0.98467,-3.46334 -2.71635,-4.34615c-0.88281,0 -0.8149,-0.8149 -0.8149,-0.8149v-1.90144c0.88282,-1.73167 1.69772,-3.42938 0.8149,-5.16106l3.53125,-16.56971c2.61448,-7.80949 3.46335,-13.78545 4.34615,-19.01442c3.46335,-15.65294 2.78426,-27.77464 -2.44471,-35.58413c-1.73167,-4.34615 -5.22896,-9.71093 -8.69231,-14.9399c-4.34615,-6.07783 -9.67698,-11.13702 -15.75481,-15.48317c-14.77013,-12.15564 -31.17007,-18.40324 -49.4375,-19.28606c-3.46334,-0.20373 -6.89273,-0.33954 -10.05048,-0.27163zM41.5601,76.60096c0.33955,-0.40745 2.17308,0.33955 5.43269,1.62981c5.22897,1.73167 9.50721,4.44802 13.85337,7.0625c2.61448,0.88282 10.42398,5.22897 21.73077,13.03846c2.61448,1.73167 7.0625,3.42938 11.40865,5.16106l11.95192,5.16106c0.88282,0.88282 0.91677,0.747 -0.8149,1.62981c-0.88281,0.88282 -3.36148,1.83354 -5.97596,2.71635l-7.0625,2.71635h-5.16106c-0.88281,0.88282 -2.71635,0.88282 -7.0625,0l-11.13702,-1.90144c-11.30679,-1.73167 -17.38462,-3.46334 -17.38462,-4.34615c-2.61448,-0.88281 -4.51592,-2.54658 -6.2476,-5.16106c-1.73167,-2.61448 -2.51262,-4.24429 -1.62981,-5.97596c-2.61448,-3.46334 -3.53125,-6.96063 -3.53125,-13.03846l0.8149,-6.2476l1.08654,-1.62981c-0.20372,-0.4414 -0.3735,-0.67909 -0.27163,-0.8149zM184.98317,77.14423c0.61118,-0.0679 1.08654,0.4414 1.08654,1.08654c0.88282,0 0.8149,0.747 0.8149,1.62981v5.43269h0.8149c-0.88281,5.22897 -1.79958,9.57512 -3.53125,13.03846c-0.88281,3.46335 -1.52794,5.90806 -3.25962,6.79087c-0.88281,1.73167 -2.61448,3.46335 -4.34615,4.34615c-0.88281,1.73167 -6.96063,3.46335 -17.38462,4.34615l-12.22356,1.90144c-3.46334,0.88282 -6.17968,0.88282 -7.0625,0c-0.88281,0.88282 -2.54658,-0.20372 -5.16106,-1.08654l-7.8774,-1.62981c-6.96063,-1.73167 -9.40535,-3.46334 -6.79087,-4.34615l13.03846,-4.34615c5.22897,-1.73167 8.52254,-4.24429 11.13702,-5.97596l11.40865,-7.0625c5.22897,-3.46334 8.59044,-5.09314 10.32212,-5.97596c3.46335,-1.73167 8.86208,-3.70102 14.9399,-5.43269c0,-1.73167 1.52794,-2.44471 3.25962,-2.44471c0.20373,-0.20372 0.61118,-0.23768 0.8149,-0.27163zM110.28365,113v15.75481h1.90144l-5.16106,7.60577l-1.08654,5.43269l-0.8149,0.8149l-2.44471,-1.90144l-1.90144,-3.25962l-1.62981,-3.53125v-2.71635l0.8149,-1.62981l1.62981,-1.90144l1.08654,-2.44471l2.44471,-6.2476v-0.8149l0.8149,-0.8149l1.08654,-0.8149l1.62981,-0.8149l0.8149,-1.90144zM114.62981,113l1.08654,1.62981l1.62981,1.08654l1.62981,1.62981v0.8149l2.71635,7.0625l3.53125,3.53125l0.8149,1.62981c0.88282,0.88282 0.88282,2.64844 0,3.53125l-1.62981,3.53125l-1.90144,3.25962l-1.62981,1.90144l-0.8149,0.8149l-1.08654,-7.8774l-5.16106,-7.8774v-11.95192z"
                    /></g
                  ></g
                ></svg
              >
            </div>
          </label>
        </div>
      </div>
    </div>
  </div>
</form>

<style>
  .open {
    display: block;
  }
  input:checked + svg {
    display: block;
  }
</style>
