<script>
    import GameForm from './GameForm.svelte';
    import X01Rules from './rules/X01Rules.svelte';
    import CricketRules from './rules/CricketRules.svelte';
    import ATCRules from './rules/ATCRules.svelte';
    import SplitScoreRules from './rules/SplitScoreRules.svelte';
    import ws from '../../../utils/socket';
    import setupGame from '../../../utils/stores/gameStore';
    import {onMount} from "svelte";
    import {goto, url} from '@roxi/routify'

    export let gameid;
    let gameMode = X01Rules;
    let headerdiv = document.getElementsByClassName('header')[0];

    // Hide navbar in this page
    headerdiv.style.display = 'block';

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

            default:
                break;
        }
    }

    onMount(() => {
        // init websocket
        const socket = ws.init(gameid, 'Game Setup Page');

        socket.addEventListener("redirect", () => {
            $goto($url(`/${gameid}/controller`))
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
