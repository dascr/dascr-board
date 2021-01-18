<script>
    import { onMount } from 'svelte';

    import setupGame from '../../../../utils/stores/gameStore';

    onMount(() => {
        $setupGame.variant = 'cut';
    });
</script>

<style>
    input:checked + svg {
        display: block;
    }
    input:disabled + svg {
        display: block;
        background: rgba(0, 0, 0, 0.5);
        cursor: not-allowed;
    }
</style>

<div class="flex flex-col">
    <label class="uppercase font-bold text-lg" for="selectVariant"
        >Select Variant:</label
    >
    <select
        bind:value="{$setupGame.variant}"
        class="border py-2 px-3 text-gray-900"
        id="selectVariant"
        name="variant"
    >
        <option value="cut">Cut-Throat</option>
        <option value="normal">Normal</option>
        <option value="no">No Score</option>
    </select>
</div>

<div class="flex flex-wrap">
    <div class="w-1/2">
        <label for="random" class="flex justify-start items-start">
            <div
                class="bg-white border-2 rounded border-gray-400 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2"
            >
                <input
                    bind:checked="{$setupGame.cricket.random}"
                    disabled="{$setupGame.cricket.ghost}"
                    type="checkbox"
                    id="random"
                    name="random"
                    class="opacity-0 absolute"
                />
                <svg
                    class="fill-current hidden w-4 h-4 text-gray-900 pointer-events-none"
                    viewBox="0 0 20 20"
                >
                    <path d="M0 11l2-2 5 5L18 3l2 2L7 18z" />
                </svg>
            </div>
            <div
                class="select-none uppercase font-bold text-lg"
                title="Pick 7 random numbers"
            >
                Random numbers
            </div>
        </label>
    </div>
    <div class="w-1/2">
        <label for="ghost" class="flex justify-start items-start">
            <div
                class="bg-white border-2 rounded border-gray-400 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2"
            >
                <input
                    bind:checked="{$setupGame.cricket.ghost}"
                    type="checkbox"
                    id="ghost"
                    name="ghost"
                    class="opacity-0 absolute"
                />
                <svg
                    class="fill-current hidden w-4 h-4 text-gray-900 pointer-events-none"
                    viewBox="0 0 20 20"
                >
                    <path d="M0 11l2-2 5 5L18 3l2 2L7 18z" />
                </svg>
            </div>
            <div
                class="select-none uppercase font-bold text-lg"
                title="Reveal random number only after first hit"
            >
                ghost
            </div>
        </label>
    </div>
</div>
