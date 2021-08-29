<script>
    export let gameid;
    import {
        insertThrow,
        miss,
        nextPlayer,
        undo,
    } from '$utils/methods';
    import state from '$stores/stateStore';

    const sendThrow = (number) => {
        insertThrow(gameid, number, $state.double, $state.triple);
        $state.double = false;
        $state.triple = false;
    };
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

<!-- Input controls -->
<div class="mt-12">
    <div class="grid grid-cols-8 grid-rows-3 gap-4">
        {#each Array(7) as _, i}
            <button
                class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
                on:click={() => {
                    sendThrow(i + 1);
                }}>{i + 1}</button>
        {/each}
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none row-span-3"
            on:click={() => nextPlayer(gameid)}>Next Player</button>

        {#each Array(13) as _, i}
            <button
                class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
                on:click={() => {
                    sendThrow(i + 8);
                }}>{i + 8}</button>
        {/each}

        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            disabled={$state.triple}
            on:click={() => {
                sendThrow(25);
            }}>25</button>
    </div>
    <div class="grid grid-cols-4 grid-rows-1 gap-4 mt-4">
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            on:click={() => miss(gameid)}>0</button>
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            class:active={$state.double}
            on:click={state.toggleDouble}>Double</button>
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            class:active={$state.triple}
            on:click={state.toggleTriple}>Triple</button>
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            on:click={() => undo(gameid)}>Undo</button>
    </div>
</div>
