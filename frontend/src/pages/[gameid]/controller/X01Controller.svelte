<script>
    import { goto, url } from '@roxi/routify';
    import { onMount } from 'svelte';
    import api from '../../../utils/api';
    import ws from '../../../utils/socket';
    import {
        transformGameMessage,
        scoreOrPodium,
    } from '../../../utils/methods';
    import TwentyFiveGrid from './inputs/TwentyFiveGrid.svelte';
    import ControllerHeader from './ControllerHeader.svelte';

    export let gameid;
    let apiBaseURL = 'API_BASE';
    let gameData = {};
    let players = [];
    let activePlayer = {};

    const update = async () => {
        const res = await api.get(`game/${gameid}/display`);
        gameData = await res.json();
        players = gameData.Player;
        activePlayer = gameData.Player[gameData.ActivePlayer];
        gameData.Message = transformGameMessage(gameData, activePlayer);
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
</style>

<!-- Header with buttons, game mode and message row -->
<ControllerHeader {gameid} {gameData}>
    <div
        slot="headerData"
        class="flex flex-row mx-auto bg-black bg-opacity-30 overflow-hidden">
        <p class="text-center border w-1/4 font-bold text-lg p-2 capitalize">
            Game:
            {gameData.Game}
        </p>
        <p class="text-center border w-1/4 font-bold text-lg p-2 capitalize">
            In:
            {gameData.In}
        </p>
        <p class="text-center border w-1/4 font-bold text-lg p-2 capitalize">
            Out:
            {gameData.Out}
        </p>
        <p class="text-center border w-1/4 font-bold text-lg p-2">
            Round:
            {gameData.ThrowRound}
        </p>
    </div>
</ControllerHeader>

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
                    {scoreOrPodium(player, gameData)}
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
                {#each Array(3 - player.LastThrows.length) as _, _}
                    <td
                        class="border-r border-l border-dashed border-opacity-10">
                        -
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
<TwentyFiveGrid {gameid} />
