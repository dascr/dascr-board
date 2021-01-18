<script>
    import Bullet from './Bullet.svelte';
    export let name,
        nickname,
        image,
        active,
        numbers,
        lastThree,
        points,
        closedNumbers,
        playerClosed,
        cricketNumbers,
        numberRevealed;
    let apiBase = 'API_BASE';
    let activeStyle =
        'border-4 border-white shadow-2xl opacity-100 transition duration-700 ease-in-out';
</script>

<div
    class="flex bg-black bg-opacity-50 rounded-lg overflow-hidden shadow-md opacity-80 px-2 py-2 {active && activeStyle}">
    <div class="mb-2 space-y-4">
        <img
            class="rounded-full hidden lg:block mx-auto"
            width="200"
            height="200"
            src={apiBase + image}
            alt="" />
        <div class="font-semibold text-4xl text-center">{name}</div>
        {#if nickname}
            <p class="text-gray-400 text-3xl text-center">{nickname}</p>
        {:else}
            <!-- Placeholder for layout reasons -->
            <p class="text-transparent text-3xl text-center">-</p>
        {/if}
        <div class="font-semibold text-3xl text-center pb-2">
            <table class="w-full">
                <tr class="">
                    {#each lastThree as lt}
                        <td>
                            {#if lt.Modifier === 2}
                                D{lt.Number}
                            {:else if lt.Modifier === 3}
                                T{lt.Number}
                            {:else}{lt.Number}{/if}
                        </td>
                    {/each}
                </tr>
            </table>
        </div>
        <div class="font-extrabold text-3xl mt-3 text-center">{points}</div>
        <div>
            {#each numbers as num, i}
                {#if closedNumbers[i]}
                    <Bullet
                        number={cricketNumbers[i]}
                        revealed={numberRevealed[i]}
                        times={3}
                        done={false}
                        closed={true} />
                {:else if !playerClosed[i]}
                    <Bullet
                        number={cricketNumbers[i]}
                        revealed={numberRevealed[i]}
                        times={num}
                        done={false}
                        closed={false} />
                {:else}
                    <Bullet
                        number={cricketNumbers[i]}
                        revealed={numberRevealed[i]}
                        times={3}
                        done={true}
                        closed={false} />
                {/if}
            {/each}
        </div>
    </div>
</div>
