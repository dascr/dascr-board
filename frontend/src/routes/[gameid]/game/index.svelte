<script>
    import GameForm from './GameForm.svelte';
    import X01Rules from './rules/X01Rules.svelte';
    import HighscoreRules from './rules/HighscoreRules.svelte';
    import CricketRules from './rules/CricketRules.svelte';
    import ATCRules from './rules/ATCRules.svelte';
    import SplitScoreRules from './rules/SplitScoreRules.svelte';
    import ShanghaiRules from './rules/ShanghaiRules.svelte';
    import setupGame from '$stores/gameStore';
    import {onMount} from "svelte";
    import {page} from "$app/stores";
    import {goto} from "$app/navigation";

    let gameid = $page.params.gameid;
    let gameMode = X01Rules;

    $: {
        let selectedGame = $setupGame.game;

        switch (selectedGame) {
            case 'x01':
                gameMode = X01Rules;
                break;
            case 'cricket':
                gameMode = CricketRules;
                break;
            case 'atc':
                gameMode = ATCRules;
                break;
            case 'split':
                gameMode = SplitScoreRules;
                break;
            case 'shanghai':
                gameMode = ShanghaiRules;
                break;
            case 'high':
                gameMode = HighscoreRules;
                break

            default:
                break;
        }
    }

    onMount(async () => {
        // Hide navbar in this page
        let headerdiv = document.getElementsByClassName('header')[0];
        headerdiv.style.display = 'block';

        // init websocket
        const im = await import('$utils/socket');
        const ws = im.default;
        const socket = ws.init(gameid, 'Game Setup Page');

        socket.addEventListener("redirect", () => {
            goto(`/${gameid}/controller`);
        })
    })
</script>

<div class="px-4 max-w-full">
    <div class="flex flex-wrap">
        <div class="w-full p-4 lg:w-1/2 space-y-4">
            <GameForm gameID={gameid}/>
        </div>
        <div class="w-full p-4 lg:w-1/2 space-y-4">
            <h3 class="text-center font-bold underline text-4xl">Game Rules</h3>
            <svelte:component this={gameMode}/>
        </div>
    </div>
</div>
