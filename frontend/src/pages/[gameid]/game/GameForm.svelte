<script>
    import { onMount } from 'svelte';
    import api from '../../../utils/api';
    import X01FormParts from './formparts/X01FormParts.svelte';
    import CricketFormParts from './formparts/CricketFormParts.svelte';
    import ATCFormParts from './formparts/ATCFormParts.svelte';
    import SplitScoreFormParts from './formparts/SplitScoreFormParts.svelte';
    import setupGame from '../../../utils/stores/gameStore';
    import { goto, url } from '@roxi/routify';

    export let gameID;
    let games = [
        { id: 'x01', text: 'X01' },
        { id: 'cricket', text: 'Cricket' },
        { id: 'atc', text: 'Around The Clock' },
        { id: 'split', text: 'Split-Score' },
    ];

    let availablePlayer = [];
    let open;
    let gameMode = X01FormParts;

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
        let playerInt = $setupGame.player.map((el) => parseInt(el));

        await api
            .post(`game/${gameID}`, {
                json: {
                    uid: gameID,
                    player: playerInt || [],
                    game: $setupGame.game || '',
                    variant: $setupGame.variant || '',
                    in: $setupGame.in || '',
                    out: $setupGame.out || '',
                    sound: $setupGame.settings.sound || false,
                    podium: $setupGame.settings.podium || false,
                    autoswitch: $setupGame.settings.autoswitch || false,
                    cricketrandom: $setupGame.cricket.random || false,
                    cricketghost: $setupGame.cricket.ghost || false,
                },
            })
            .then(async (res) => {
                if (res.status === 200) {
                    $goto($url(`/${gameID}/controller`));
                }
            })
            .catch((err) => console.error(err));
    }
</script>

<style>
    .open {
        display: block;
    }
    input:checked + svg {
        display: block;
    }
</style>

<h3 class="text-center font-bold underline text-4xl">Game Setup</h3>
<form
    id="newGame"
    enctype="multipart/form-data"
    autocomplete="off"
    on:submit|preventDefault={handleSubmit}>
    <div class="space-y-4">
        <div class="flex flex-col">
            <label for="selectPlayer" class="uppercase font-bold text-lg">Select
                Players:</label>
            <select
                bind:value={$setupGame.player}
                class="border py-2 px-3 text-gray-900"
                id="selectPlayer"
                name="player"
                multiple
                required>
                {#each availablePlayer || [] as player}
                    <option value={player.UID}>{player.Name}</option>
                {/each}
            </select>
        </div>

        <div class="flex flex-col">
            <label for="selectGame" class="uppercase font-bold text-lg">Choose
                your game type:</label>
            <select
                bind:value={$setupGame.game}
                id="selectGame"
                class="border py-2 px-3 text-gray-900"
                name="game"
                required>
                {#each games as game}
                    <option value={game.id}>{game.text}</option>
                {/each}
            </select>
        </div>

        <!-- dynamic form parts -->
        <svelte:component this={gameMode} />

        <button
            type="submit"
            class="block uppercase border-2 hover:bg-black hover:bg-opacity-30 text-lg mx-auto p-4 rounded-2xl"><i
                class="fas fa-play pr-2" />
            Start Game</button>

        <!-- Further settings -->

        <div class="flex flex-col">
            <button
                class="block uppercase border-2 hover:bg-black hover:bg-opacity-30 text-lg mx-auto p-2 rounded-2xl"
                type="button"
                on:click={() => (open = !open)}><i
                    class="fas fa-sliders-h pr-2" />Further Settings</button>
        </div>

        <div class:open class="hidden">
            <div class="flex flex-wrap">
                <div class="w-1/3">
                    <label for="sound" class="flex justify-start items-start">
                        <div
                            class="bg-white border-2 rounded border-gray-400 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2">
                            <input
                                bind:checked={$setupGame.settings.sound}
                                type="checkbox"
                                id="sound"
                                name="sound"
                                class="opacity-0 absolute" />
                            <svg
                                class="fill-current hidden w-4 h-4 text-gray-900 pointer-events-none"
                                viewBox="0 0 20 20"><path
                                    d="M0 11l2-2 5 5L18 3l2 2L7 18z" /></svg>
                        </div>
                        <div
                            class="select-none uppercase font-bold text-lg"
                            title="Play sound">
                            Sound
                        </div>
                    </label>
                </div>

                <div class="w-1/3">
                    <label for="podium" class="flex justify-start items-start">
                        <div
                            class="bg-white border-2 rounded border-gray-400 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2">
                            <input
                                bind:checked={$setupGame.settings.podium}
                                type="checkbox"
                                id="podium"
                                name="podium"
                                class="opacity-0 absolute" />
                            <svg
                                class="fill-current hidden w-4 h-4 text-gray-900 pointer-events-none"
                                viewBox="0 0 20 20"><path
                                    d="M0 11l2-2 5 5L18 3l2 2L7 18z" /></svg>
                        </div>
                        <div
                            class="select-none uppercase font-bold text-lg"
                            title="Continue game after a player has won until one player is left">
                            podium
                        </div>
                    </label>
                </div>

                <div class="w-1/3">
                    <label
                        for="autoswtich"
                        class="flex justify-start items-start">
                        <div
                            class="bg-white border-2 rounded border-gray-400 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2">
                            <input
                                bind:checked={$setupGame.settings.autoswitch}
                                type="checkbox"
                                id="autoswitch"
                                name="autoswitch"
                                class="opacity-0 absolute" />
                            <svg
                                class="fill-current hidden w-4 h-4 text-gray-900 pointer-events-none"
                                viewBox="0 0 20 20"><path
                                    d="M0 11l2-2 5 5L18 3l2 2L7 18z" /></svg>
                        </div>
                        <div
                            class="select-none uppercase font-bold text-lg"
                            title="When not using an automated recognition (cam / machine) this will switch to next player automatically when 3 throws are inserted after a delay of 2 seconds">
                            Auto Switch next player
                        </div>
                    </label>
                </div>
            </div>
        </div>
    </div>
</form>
