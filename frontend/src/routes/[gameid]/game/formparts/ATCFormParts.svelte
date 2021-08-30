<script>
    import { onMount } from 'svelte';

    import setupGame from '$stores/gameStore';

    onMount(() => {
        $setupGame.variant = 'normal';
    });
</script>

<div class="flex flex-col">
    <label class="uppercase font-bold text-lg" for="selectVariant">Select
        Variant:</label>
    <select
        bind:value={$setupGame.variant}
        class="border py-2 px-3 text-gray-900"
        id="selectVariant"
        name="variant">
        <option value="normal">Normal</option>
        <option value="fast">Fast</option>
    </select>
</div>
