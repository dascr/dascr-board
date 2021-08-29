<script>
    export let times, done, closed, number, revealed;
</script>

<style>
    .done {
        color: white;
        text-shadow: 2px 2px var(--bg-color-1);
    }
    .closed {
        opacity: 0.1;
    }
</style>

<div class="grid grid-cols-4 grid-rows-1 content-center" class:closed>
    <div class="text-gray-400 text-2xl font-extrabold" class:done>
        {#if !revealed}?:{:else}{number}:{/if}
    </div>
    {#each Array(times) as _}
        <svg height="32" width="112">
            <ellipse
                cx="32"
                cy="16"
                rx="16"
                ry="8"
                style="fill:white;stroke:white;stroke-width:2" />
        </svg>
    {/each}
    {#each Array(3 - times) as _}
        <svg height="32" width="112">
            <ellipse
                cx="32"
                cy="16"
                rx="16"
                ry="8"
                style="fill:black;stroke:white;stroke-width:2"
                opacity="0.5" />
        </svg>
    {/each}
</div>
