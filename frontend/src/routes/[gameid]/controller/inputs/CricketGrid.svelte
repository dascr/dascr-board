<script>
    import {
        nextPlayer,
        insertThrow,
        undo,
        miss,
    } from '$utils/methods';
    import TwentyFiveGrid from './TwentyFiveGrid.svelte';
    import state from '$stores/stateStore';

    export let revealed, numbers, allRevealed, gameid;

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

<div class="mt-12">
    <div class="grid grid-cols-8 grid-rows-1 gap-4">
        {#each Array(7) as _, i}
            {#if !revealed[i]}
                <button
                    class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
                    disabled
                    on:click={() => {
                        sendThrow(0);
                    }}>?</button>
            {:else}
                <button
                    class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
                    disabled={numbers[i] === 25 && $state.triple}
                    on:click={() => {
                        sendThrow(numbers[i]);
                    }}>{numbers[i]}</button>
            {/if}
        {/each}
        <button
            class="text-2xl font-extrabold p-7 bg-black bg-opacity-30 hover:bg-opacity-50 border text-center rounded focus:outline-none"
            on:click={() => nextPlayer(gameid)}>Next Player</button>
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

    <!-- for ghost mode, hide when all revealed -->
    {#if !allRevealed}
        <TwentyFiveGrid {gameid} />
    {/if}
</div>
