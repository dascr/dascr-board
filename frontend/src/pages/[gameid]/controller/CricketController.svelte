<script>
    import api from '../../../utils/api';
    import { transformGameMessage } from '../../../utils/methods';
    import { onMount } from 'svelte';
    import ws from '../../../utils/socket';
    import { goto, url } from '@roxi/routify';

    const checker = (arr) => arr.every((v) => v === true);

    export let gameid;
    let apiBaseURL = 'API_BASE';
    let gameData = {};
    let players = [];
    let numbers = [];
    let revealed = [];
    let allRevealed = false;
    let activePlayer = {};
    let mode = '';
    let randomGhost = '';
    let double = false;
    let triple = false;
    let autoswitch = false;

    const update = async () => {
        const res = await api.get(`game/${gameid}/display`);
        gameData = await res.json();
        players = gameData.Player;
        activePlayer = gameData.Player[gameData.ActivePlayer];
        numbers = gameData.CricketController.Numbers;
        revealed = gameData.CricketController.NumberRevealed;
        autoswitch = gameData.Settings.AutoSwitch;
        gameData.Message = transformGameMessage(gameData, activePlayer);

        switch (gameData.Variant) {
            case 'cut':
                mode = 'Cut Throat';
                break;
            case 'normal':
                mode = 'Normal';
                break;
            case 'no':
                mode = 'No Score';
                break;
        }

        if (gameData.CricketController.Ghost) {
            randomGhost = 'Yes / Yes';
        } else if (gameData.CricketController.Random) {
            randomGhost = 'Yes / No';
        } else {
            randomGhost = 'No / No';
        }

        // Check if all are revealed
        allRevealed = checker(revealed);

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
        const socket = ws.init(gameid, 'Cricket Controller');

        await update();

        socket.addEventListener('update', () => {
            update();
        });

        socket.addEventListener('redirect', () => {
            $goto($url(`/${gameid}/game`));
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
    <p class="text-center border w-1/4 font-bold text-lg p-2">Mode: {mode}</p>
    <p class="text-center border w-1/4 font-bold text-lg p-2">
        Random / Ghost:
        {randomGhost}
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
                Points
            </td>
            {#each numbers as num, i}
                <td class="border-l border-r border-dashed font-bold px-3">
                    {#if revealed[i]}{num}{:else}?{/if}
                </td>
            {/each}
            <td
                class="border-l border-r border-dashed font-bold px-3"
                colspan="3">
                Throws
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
                {#each player.Score.Numbers as num}
                    <td
                        class="px-3 border-r border-l border-dashed border-opacity-10">
                        {num}
                    </td>
                {/each}
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
                <td class="px-3 border-l border-dashed border.opacity-10">
                    {player.TotalThrowCount}
                </td>
            </tr>
        {/each}
    </tbody>
</table>

<!-- Input controls -->
<div class="mt-12">
    <div class="grid grid-cols-8 grid-rows-1 gap-4">
        {#each Array(7) as _, i}
            {#if !revealed[i]}
                <button
                    class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
                    disabled
                    on:click={() => {
                        insertThrow(0);
                    }}>?</button>
            {:else}
                <button
                    class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
                    disabled={numbers[i] === 25 && triple}
                    on:click={() => {
                        insertThrow(numbers[i]);
                    }}>{numbers[i]}</button>
            {/if}
        {/each}
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            on:click={nextPlayer}>Next Player</button>
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

    <!-- for ghost mode, hide when all revealed -->
    {#if !allRevealed}
        <div class="mt-12">
            <div class="grid grid-cols-7 grid-rows-3 gap-4">
                {#each Array(20) as _, i}
                    <button
                        class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
                        on:click={() => {
                            insertThrow(i + 1);
                        }}>{i + 1}</button>
                {/each}

                <button
                    class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
                    disabled={triple}
                    on:click={() => {
                        insertThrow(25);
                    }}>25</button>
            </div>
        </div>
    {/if}
</div>
