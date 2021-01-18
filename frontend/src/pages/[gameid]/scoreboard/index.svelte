<script>
    import api from '../../../utils/api';
    import X01Scoreboard from './X01Scoreboard.svelte';
    import CricketScoreboard from './CricketScoreboard.svelte';

    export let gameid;
    let headerdiv = document.getElementsByClassName('header')[0];

    // Hide navbar in this page
    headerdiv.style.display = 'none';

    const getGame = async () => {
        let gameMode

        const res = await api.get(`game/${gameid}`);
        let game = await res.json();

        switch (game.Game) {
            case 'x01':
                gameMode = X01Scoreboard;
                break;

            case 'cricket':
                gameMode = CricketScoreboard;
                break;

            default:
                throw new Error("Controller cannot be shown")
                break;
        }
        return gameMode
    }

    let promise = getGame()
</script>

{#await promise}
    <p>...loading</p>{:then game}
    <svelte:component this={game} {gameid}/>
{:catch error}
    <p>There was an error {error.message}
{/await   }
