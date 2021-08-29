<script>
    import { scoreOrPodium } from '$utils/methods';
    import Bullet from './Bullet.svelte';
    import myenv from '$utils/env';

    export let player, gameData, active;

    let apiBase = myenv.apiBase;
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
            src={apiBase + player.Image}
            alt="" />
        <div class="font-semibold text-4xl text-center">{player.Name}</div>
        {#if player.Nickname}
            <p class="text-gray-400 text-3xl text-center">{player.Nickname}</p>
        {:else}
            <!-- Placeholder for layout reasons -->
            <p class="text-transparent text-3xl text-center">-</p>
        {/if}
        <div class="font-semibold text-3xl text-center pb-2">
            <table class="w-full">
                <tr class="">
                    {#each player.LastThrows as lt}
                        <td>
                            {#if lt.Modifier === 2}
                                D{lt.Number}
                            {:else if lt.Modifier === 3}
                                T{lt.Number}
                            {:else}{lt.Number}{/if}
                        </td>
                    {/each}

                    {#each Array(3 - player.LastThrows.length) as _, __}
                        <td>-</td>
                    {/each}
                </tr>
            </table>
        </div>
        <div class="font-extrabold text-3xl mt-3 text-center">
            Points:
            {scoreOrPodium(player, gameData)}
        </div>
        <div>
            {#each player.Score.Numbers as num, i}
                {#if gameData.CricketController.NumberClosed[i]}
                    <Bullet
                        number={gameData.CricketController.Numbers[i]}
                        revealed={gameData.CricketController.NumberRevealed[i]}
                        times={3}
                        done={false}
                        closed={true} />
                {:else if !player.Score.Closed[i]}
                    <Bullet
                        number={gameData.CricketController.Numbers[i]}
                        revealed={gameData.CricketController.NumberRevealed[i]}
                        times={num}
                        done={false}
                        closed={false} />
                {:else}
                    <Bullet
                        number={gameData.CricketController.Numbers[i]}
                        revealed={gameData.CricketController.NumberRevealed[i]}
                        times={3}
                        done={true}
                        closed={false} />
                {/if}
            {/each}
        </div>
    </div>
</div>
