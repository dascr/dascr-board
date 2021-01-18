<script>
    import { onMount } from 'svelte';
    import api from '../../../utils/api';
    import { transformGameMessage } from '../../../utils/methods';
    import ws from '../../../utils/socket';
    import { goto, url } from '@roxi/routify';
    import CricketCard from './CricketCard.svelte';

    export let gameid;
    let gameData = {};
    let players = [];
    let activePlayer = {};
    let mode = '';
    let randomGhost = '';

    const update = async () => {
        const res = await api.get(`game/${gameid}/display`);
        gameData = await res.json();
        players = gameData.Player;
        activePlayer = gameData.Player[gameData.ActivePlayer];
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
    };

    onMount(async () => {
        // init websocket
        const socket = ws.init(gameid, 'Cricket Scoreboard');

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
    <p
        class="text-center border w-1/4 font-bold text-lg rounded-tl-2xl p-2 capitalize">
        Game:
        {gameData.Game}
    </p>
    <p class="text-center border w-1/4 font-bold text-lg p-2">Mode: {mode}</p>
    <p class="text-center border w-1/4 font-bold text-lg p-2">
        Random / Ghost:
        {randomGhost}
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

<div class="mt-3">
    <div class="flex flex-wrap max-w-full">
        <!-- Player cols -->
        {#each players as player, i}
            <div class="w-72 mx-2 my-2">
                <CricketCard
                    name={player.Name}
                    nickname={player.Nickname}
                    image={player.Image}
                    active={i === gameData.ActivePlayer}
                    numbers={player.Score.Numbers}
                    playerClosed={player.Score.Closed}
                    lastThree={player.LastThrows}
                    points={gameData.Podium.includes(player.UID) ? 'Place ' + (gameData.Podium.indexOf(player.UID) + 1) : 'Points: ' + player.Score.Score}
                    closedNumbers={gameData.CricketController.NumberClosed}
                    cricketNumbers={gameData.CricketController.Numbers}
                    numberRevealed={gameData.CricketController.NumberRevealed} />
            </div>
        {/each}
    </div>
</div>
