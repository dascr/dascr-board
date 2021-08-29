<script>
    import {
        setCricketModeHeader,
        scoreOrPodium,
    } from '$utils/methods';
    import { onMount } from 'svelte';
    import ws from '$utils/socket';
    import ControllerHeader from './ControllerHeader.svelte';
    import CricketGrid from './inputs/CricketGrid.svelte';
    import state from '$stores/stateStore';
    import {page} from '$app/stores';
    import {goto} from '$app/navigation';
    import myenv from '$utils/env';

    let apiBaseURL = myenv.apiBase;
    let gameid = $page.params.gameid;

    let mode = '';
    let randomGhost = '';

    onMount(async () => {
        // init websocket
        const socket = ws.init(gameid, 'Cricket Controller');

        await state.updateState(gameid);

        const res = setCricketModeHeader($state.gameData);
        mode = res[0];
        randomGhost = res[1];

        // allRevealed = checker(revealed);

        socket.addEventListener('update', async () => {
            await state.updateState(gameid);
        });

        socket.addEventListener('redirect', () => {
            goto(`/${gameid}/game`);
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
<ControllerHeader {gameid} gameData={$state.gameData}>
    <div
        slot="headerData"
        class="flex flex-row mx-auto bg-black bg-opacity-30 overflow-hidden">
        <p class="text-center border w-1/4 font-bold text-lg p-2 capitalize">
            Game:
            {$state.gameData.Game}
        </p>
        <p class="text-center border w-1/4 font-bold text-lg p-2">
            Mode:
            {mode}
        </p>
        <p class="text-center border w-1/4 font-bold text-lg p-2">
            Random / Ghost:
            {randomGhost}
        </p>
        <p class="text-center border w-1/4 font-bold text-lg p-2">
            Round:
            {$state.gameData.ThrowRound}
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
                Points
            </td>
            {#each $state.numbers as num, i}
                <td class="border-l border-r border-dashed font-bold px-3">
                    {#if $state.revealed[i]}{num}{:else}?{/if}
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
        {#each $state.players as player, i}
            <tr
                class="text-center flex-none"
                class:active={i === $state.gameData.ActivePlayer}>
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
                    {scoreOrPodium(player, $state.gameData)}
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
                {#each Array(3 - player.LastThrows.length) as _, __}
                    <td
                        class="border-r border-l border-dashed border-opacity-10">
                        -
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
<CricketGrid
    {gameid}
    revealed={$state.revealed}
    numbers={$state.numbers}
    allRevealed={$state.allRevealed} />
