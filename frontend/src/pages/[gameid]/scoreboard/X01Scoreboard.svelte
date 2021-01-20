<script>
    import { url, goto } from '@roxi/routify';
    import api from '../../../utils/api';
    import { onMount } from 'svelte';
    import ws from '../../../utils/socket';
    import PlayerCard from '../../player/PlayerCard.svelte';
    import { transformGameMessage } from '../../../utils/methods';

    export let gameid;
    let gameData = {};
    let players = [];
    let activePlayer = {};

    async function update() {
        const res = await api.get(`game/${gameid}/display`);
        gameData = await res.json();
        players = gameData.Player;
        activePlayer = gameData.Player[gameData.ActivePlayer];
        gameData.Message = transformGameMessage(gameData, activePlayer);
    }

    onMount(async () => {
        // init websocket
        const socket = ws.init(gameid, 'X01 Scoreboard');

        update();

        socket.addEventListener('update', () => {
            update();
        });

        socket.addEventListener('redirect', () => {
            $goto($url(`/${gameid}/start`));
        });
    });
</script>

<div
    class="flex flex-row mx-auto bg-black bg-opacity-30 rounded-t-2xl overflow-hidden">
    <p class="text-center border w-1/4 font-bold text-lg rounded-tl-2xl p-2">
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
    <p class="text-center border w-1/4 font-bold text-lg rounded-tr-2xl p-2">
        Round:
        {gameData.ThrowRound}
    </p>
</div>
<div class="bg-black bg-opacity-30 rounded-b-2xl overflow-hidden">
    <p
        class="text-center border w-full font-extrabold text-4xl rounded-b-2xl p-2">
        {gameData.Message}
    </p>
</div>

<div class="max-w-full space-y-2">
    <div class="flex flex-wrap">
        {#each players as player, i}
            <div class="w-full p-2 2xl:w-1/2">
                <PlayerCard
                    uid={player.UID}
                    name={player.Name}
                    nickname={player.Nickname}
                    image={player.Image}
                    showDelete={false}
                    onDelete={() => {}}
                    active={i === gameData.ActivePlayer}>
                    <div slot="points">
                        <p class="font-extrabold text-5xl mt-5">
                            {gameData.Podium.includes(player.UID) ? 'Place ' + (gameData.Podium.indexOf(player.UID) + 1) : player.Score.Score}
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
                                {#each player.LastThrows as thr}
                                    <td
                                        class="p-2 text-4xl font-extrabold text-center w-1/4 border-r border-dashed border-white border-opacity-10">
                                        {#if thr.Modifier === 2}
                                            D{thr.Number}
                                        {:else if thr.Modifier === 3}
                                            T{thr.Number}
                                        {:else}{thr.Number}{/if}
                                    </td>
                                {/each}
                                {#each Array(3 - player.LastThrows.length) as _, _}
                                    <td
                                        class="p-2 text-4xl font-extrabold text-center w-1/4 border-r border-dashed border-white border-opacity-10">
                                        -
                                    </td>
                                {/each}
                                <td
                                    rowspan="2"
                                    class="pl-4 p-2 text-4xl font-extrabold text-left w-1/4">
                                    <div class="flex flex-row">
                                        <div class="mr-2">Ø</div>
                                        <div>{player.Average}</div>
                                    </div>
                                </td>
                            </tr>
                            <tr class="">
                                <td
                                    colspan="3"
                                    class="text-center p-2 text-4xl font-extrabold w-3/4 border-t border-white border-dashed border-opacity-10">
                                    {player.ThrowSum}
                                </td>
                            </tr>
                        </table>
                    </div>
                </PlayerCard>
            </div>
        {/each}
    </div>
</div>
