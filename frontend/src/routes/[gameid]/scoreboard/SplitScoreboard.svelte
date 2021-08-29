<script>
    import { onMount } from 'svelte';
    import PlayerCard from '../../player/PlayerCard.svelte';
    import state from '$stores/stateStore';
    import {page} from '$app/stores';
    import {goto} from '$app/navigation';
    import { scoreOrHitorder } from '$utils/methods';

    let gameid = $page.params.gameid;

    const hitOrder = [
        '15',
        '16',
        'Any Double',
        '17',
        '18',
        'Any Triple',
        '19',
        '20',
        '25',
    ];

    onMount(async () => {
        // init websocket
        const ws = await import('$utils/socket');
        const socket = ws.init(gameid, 'ATC Scoreboard');

        await state.updateState(gameid);

        socket.addEventListener('update', async () => {
            await state.updateState(gameid);
        });

        socket.addEventListener('redirect', () => {
            goto(`/${gameid}/start`);
        });
    });
</script>

<div
    class="flex flex-row mx-auto bg-black bg-opacity-30 rounded-t-2xl overflow-hidden">
    <p class="text-center border w-1/3 font-bold text-lg rounded-tl-2xl p-2">
        Game: Split-Score
    </p>
    <p class="text-center border w-1/3 font-bold text-lg p-2">
        Variant:
        {#if $state.gameData.Variant === 'edart'}E-Dart{:else}Steel Dart{/if}
    </p>
    <p class="text-center border w-1/3 font-bold text-lg rounded-tr-2xl p-2">
        Round:
        {$state.gameData.ThrowRound}
    </p>
</div>

<div class="bg-black bg-opacity-30 rounded-b-2xl overflow-hidden">
    <p
        class="text-center border w-full font-extrabold text-4xl rounded-b-2xl p-2">
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
                    active={i === $state.gameData.ActivePlayer}>
                    <div slot="points">
                        <p class="font-extrabold text-5xl mt-5">
                            {player.Score.Score}
                        </p>
                        <p class="font-semibold text-4xl mt-5">
                            {#if $state.gameData.Variant === 'steel'}
                                {#if $state.gameData.ThrowRound >= 2}
                                    {scoreOrHitorder(player, $state.gameData, hitOrder[$state.gameData.ThrowRound - 2])}
                                {:else}Throw Start Score{/if}
                            {:else}
                                {scoreOrHitorder(player, $state.gameData, hitOrder[$state.gameData.ThrowRound - 1])}
                            {/if}
                        </p>
                        <p class="font-semibold  text-2xl mt-5 flex flex-row">
                            <img
                                src="/img/dart.png"
                                width="32"
                                height="32"
                                alt=""
                                class="mr-4" />
                            {player.TotalThrowCount}
                        </p>
                    </div>
                    <div slot="score" class="h-full">
                        <table class="w-full h-full">
                            <tr class="">
                                {#each player.LastThrows as thr, i}
                                    <td
                                        class="p-2 text-4xl font-extrabold text-center w-1/4 border-dashed border-white border-opacity-10"
                                        class:border-l={i != 0}>
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
                                        class:border-l={i != 0}>
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
