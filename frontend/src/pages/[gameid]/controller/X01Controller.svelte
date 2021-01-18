<script>
    import { goto, url } from '@roxi/routify';
    import { onMount } from 'svelte';
    import api from '../../../utils/api';
    import ws from '../../../utils/socket';
    import { transformGameMessage } from '../../../utils/methods';

    export let gameid;
    let apiBaseURL = 'API_BASE';
    let gameData = {};
    let players = [];
    let activePlayer = {};
    let double = false;
    let triple = false;
    let autoswitch = false;

    const update = async () => {
        const res = await api.get(`game/${gameid}/display`);
        gameData = await res.json();
        players = gameData.Player;
        activePlayer = gameData.Player[gameData.ActivePlayer];
        gameData.Message = transformGameMessage(gameData, activePlayer);
        autoswitch = gameData.Settings.AutoSwitch;

        // Check for autoswitch and do the switch
        if (gameData.Settings.AutoSwitch) {
            if (
                gameData.GameState === 'NEXTPLAYER' ||
                gameData.GameState === 'NEXTPLAYERWON' ||
                gameData.GameState.includes('BUST')
            ) {
                setTimeout(() => {
                    nextPlayer();
                }, 2000);
            }
        }
    };

    const endGame = () => {
        if (confirm('Really end game?')) {
            api.delete(`game/${gameid}`);
        }
    };

    const rematch = () => {
        api.post(`game/${gameid}/rematch`);
    };

    const nextPlayer = () => {
        api.post(`game/${gameid}/nextPlayer`);
    };

    const undo = () => {
        api.post(`game/${gameid}/undo`);
    };

    const toggleDouble = () => {
        if (double) {
            double = false;
            triple = false;
        } else if (!double) {
            double = true;
            triple = false;
        }
    };

    const toggleTriple = () => {
        if (triple) {
            double = false;
            triple = false;
        } else if (!triple) {
            double = false;
            triple = true;
        }
    };

    const insertThrow = (number) => {
        let modifier = 1;
        if (double) {
            modifier = 2;
        }
        if (triple) {
            modifier = 3;
        }
        api.post(`game/${gameid}/throw/${number}/${modifier}`);

        double = false;
        triple = false;
        navigator.vibrate(200);
    };

    onMount(async () => {
        // init websocket
        const socket = ws.init(gameid, 'X01 Controller');

        update();

        socket.addEventListener('redirect', () => {
            $goto($url(`/${gameid}/game`));
        });

        socket.addEventListener('update', () => {
            update();
        });
    });
</script>

<style>
    .active {
        background: rgba(0, 0, 0, 0.5);
        border: 2px solid white;
    }
    button:disabled {
        background: rgba(126, 126, 126, 0.5);
        cursor: not-allowed;
    }
</style>

<!-- Header with buttons, game mode and message row -->
<div class="flex flex-row">
    <button
        class="text-center border w-1/2 font-extrabold text-2xl rounded-tl-2xl p-3 bg-black bg-opacity-30 hover:bg-opacity-50 focus:outline-none"
        on:click={endGame}>End Game</button>
    <button
        class="text-center border w-1/2 font-extrabold text-2xl rounded-tr-2xl p-3 bg-black bg-opacity-30 hover:bg-opacity-50 focus:outline-none"
        class:hidden={gameData.GameState !== 'WON'}
        on:click={rematch}>Rematch</button>
    <button
        class="text-center border w-1/2 font-extrabold text-2xl rounded-tr-2xl p-3 bg-black bg-opacity-30 hover:bg-opacity-50 focus:outline-none"
        class:hidden={gameData.GameState === 'WON'}
        disabled={autoswitch}
        on:click={nextPlayer}>Next Player</button>
</div>
<div class="flex flex-row mx-auto bg-black bg-opacity-30 overflow-hidden">
    <p class="text-center border w-1/4 font-bold text-lg p-2">
        Game:
        {gameData.Game}
    </p>
    <p class="text-center border w-1/4 font-bold text-lg p-2">
        In:
        {gameData.In}
    </p>
    <p class="text-center border w-1/4 font-bold text-lg p-2">
        Out:
        {gameData.Out}
    </p>
    <p class="text-center border w-1/4 font-bold text-lg p-2">
        Round:
        {gameData.ThrowRound}
    </p>
</div>
<div class="bg-black bg-opacity-30 rounded-b-2xl overflow-hidden">
    <p
        class="text-center border w-full font-extrabold text-2xl rounded-b-2xl p-2">
        {gameData.Message}
    </p>
</div>

<!-- Table with player details -->
<table class="w-full mt-3">
    <thead class="mx-3 border-b border-dashed border-opacity-10">
        <tr class="text-center">
            <td class="" />
            <td class="border-r border-dashed font-bold px-3 text-left">
                Player
            </td>
            <td class="border-l border-r border-dashed font-bold px-3">
                Score
            </td>
            <td class="border-l border-r border-dashed font-bold px-3">
                Average
            </td>
            <td
                class="border-l border-r border-dashed font-bold px-3"
                colspan="3">
                Throws
            </td>
            <td class="border-l border-r border-dashed font-bold px-3">
                Total
            </td>
            <td class="border-l border-dashed font-bold px-3">
                <div class="flex justify-center">
                    <img
                        src="/img/dart.png"
                        width="16"
                        height="16"
                        alt=""
                        class="" />
                </div>
            </td>
        </tr>
    </thead>
    <tbody>
        {#each players as player, i}
            <tr
                class="text-center flex-none"
                class:active={i === gameData.ActivePlayer}>
                <td>
                    <img
                        src={`${apiBaseURL}${player.Image}`}
                        class="rounded-full"
                        alt=""
                        height="32"
                        width="32" />
                </td>
                <td
                    class="px-3 text-left border-r border-dashed border-opacity-10">
                    {player.Name}
                    {player.Nickname && '- ' + player.Nickname}
                </td>
                <td
                    class="px-3 border-r border-l border-dashed border-opacity-10">
                    {gameData.Podium.includes(player.UID) ? 'Place ' + (gameData.Podium.indexOf(player.UID) + 1) : player.Score.Score}
                </td>
                <td
                    class="px-3 border-r border-l border-dashed border-opacity-10">
                    {player.Average}
                </td>

                {#each player.LastThrows as thr}
                    <td
                        class="border-r border-l border-dashed border-opacity-10">
                        {#if thr.Modifier === 2}
                            D
                            {thr.Number}
                        {:else if thr.Modifier === 3}
                            T
                            {thr.Number}
                        {:else}{thr.Number}{/if}
                    </td>
                {/each}
                <td
                    class="px-3 border-l border-r border-dashed border.opacity-10">
                    {player.ThrowSum}
                </td>
                <td class="px-3 border-l border-dashed border.opacity-10">
                    {player.TotalThrowCount}
                </td>
            </tr>
        {/each}
    </tbody>
</table>

<!-- Input controls -->
<div class="mt-12">
    <div class="grid grid-cols-8 grid-rows-3 gap-4">
        {#each Array(7) as _, i}
            <button
                class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
                on:click={() => {
                    insertThrow(i + 1);
                }}>{i + 1}</button>
        {/each}
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none row-span-3"
            on:click={nextPlayer}>Next Player</button>

        {#each Array(13) as _, i}
            <button
                class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
                on:click={() => {
                    insertThrow(i + 8);
                }}>{i + 8}</button>
        {/each}

        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            disabled={triple}
            on:click={() => {
                insertThrow(25);
            }}>25</button>
    </div>
    <div class="grid grid-cols-4 grid-rows-1 gap-4 mt-4">
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            on:click={() => {
                insertThrow(0);
            }}>0</button>
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            class:active={double}
            on:click={toggleDouble}>Double</button>
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            class:active={triple}
            on:click={toggleTriple}>Triple</button>
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            on:click={undo}>Undo</button>
    </div>
    <!-- <div class="grid grid-cols-1 grid-rows-1 gap-4 mt-4">
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            on:click={nextPlayer}>Next Player</button>
    </div> -->
</div>
