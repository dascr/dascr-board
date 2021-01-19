<script>
    import { url } from '@roxi/routify';

    import { onMount } from 'svelte';
    import api from '../utils/api';
    let baseURL = window.location.protocol + '//' + window.location.host;
    let allGames = [];
    onMount(async () => {
        const res = await api.get('game');
        allGames = await res.json();
    });
</script>

<!-- Intro row -->
<div class="flex flex-col space-y-4">
    <div>
        <h2 class="text-center font-bold text-2xl">Introduction</h2>
        <p class="text-center">
            Hello and welcome to your personal darts game sponsored and brought
            to you by
            <i>DaSCR</i>
            - Your friendly darts tracking and scoring system.
        </p>
    </div>
    <div>
        <h3 class="text-center font-bold text-xl mb-2">
            Setup your personal game
        </h3>
        <p class="text-center">
            To start a game you choose a unique game id within the url. So there
            can be more running games concurrently.
        </p>
        <p class="text-center">
            To start a game with the ID "dascr" for example you navigate your
            scoreboard machine to:
        </p>
        <p class="text-center">
            <a
                href={$url(`/dascr/start`)}
                class="font-bold">{baseURL}/dascr/start</a>
        </p>
        <p class="text-center">
            You will then be presented with the QR Code to scan with your
            smartphone or tablet to setup your personal game space.
        </p>
    </div>
    <div>
        <h3 class="text-center font-bold text-xl">Player Management</h3>
        <p class="text-center">
            To manage the available players you can use the menu above (top
            left) or navigate your browser to:
        </p>
        <p class="text-center">
            <a href={$url(`/player`)} class="font-bold">{baseURL}/player</a>
        </p>
    </div>
    <!-- hr -->
    <div>
        <hr class="pl-5 pr-5" />
    </div>
    <!-- Running games -->
    <h2 class="text-2xl font-bold text-center">Running games</h2>
    <table class="table-fixed w-full text-center" id="players-table">
        <thead>
            <tr>
                <th class="w-1/4">Game ID</th>
                <th class="w-1/4">Game</th>
                <th class="w-1/4">Variant</th>
                <th class="w-1/4"># of Players</th>
            </tr>
        </thead>
        <tbody />
        {#each allGames || [] as game}
            <tr id={game.uid}>
                <td>
                    <a
                        href={$url(`/${game.uid}/scoreboard`)}
                        class="font-bold">{game.uid}</a>
                </td>
                <td>{game.game}</td>
                <td>{game.variant}</td>
                <td>{game.player.length}</td>
            </tr>
        {/each}
    </table>
</div>
