<script>
    import {
        insertThrowNumberMod,
        nextPlayer,
        miss,
        undo,
    } from '$utils/methods';

    export let gameid, activePlayer;
</script>

<div class="mt-12">
    <div class="grid grid-cols-3 grid-rows-2 gap-4">
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            on:click={() => {
                insertThrowNumberMod(gameid, activePlayer.Score.CurrentNumber, 1);
            }}>Single</button>
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            on:click={() => {
                insertThrowNumberMod(gameid, activePlayer.Score.CurrentNumber, 2);
            }}>Double</button>
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none row-span-2"
            on:click={() => nextPlayer(gameid)}>Next Player</button>
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            on:click={() => miss(gameid)}>Miss</button>
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            on:click={() => undo(gameid)}>Undo</button>
    </div>
</div>
