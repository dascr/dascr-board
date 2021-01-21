<script>
    export let gameid, gameData;

    import { endGame, rematch, nextPlayer } from '../../../utils/methods';
</script>

<div class="flex flex-row">
    <button
        class="text-center border w-1/2 font-extrabold text-2xl rounded-tl-2xl p-3 bg-black bg-opacity-30 hover:bg-opacity-50 focus:outline-none"
        on:click={() => endGame(gameid)}>End Game</button>
    <button
        class="text-center border w-1/2 font-extrabold text-2xl rounded-tr-2xl p-3 bg-black bg-opacity-30 hover:bg-opacity-50 focus:outline-none"
        class:hidden={gameData.GameState !== 'WON'}
        on:click={() => rematch(gameid)}>Rematch</button>
    <button
        class="text-center border w-1/2 font-extrabold text-2xl rounded-tr-2xl p-3 bg-black bg-opacity-30 hover:bg-opacity-50 focus:outline-none"
        class:hidden={gameData.GameState === 'WON'}
        on:click={() => nextPlayer(gameid)}>Next Player</button>
</div>
{#if $$slots.headerData}
    <slot name="headerData" />
{/if}
<div class="bg-black bg-opacity-30 rounded-b-2xl overflow-hidden">
    <p
        class="text-center border w-full font-extrabold text-2xl rounded-b-2xl p-2">
        {gameData.Message}
    </p>
</div>
