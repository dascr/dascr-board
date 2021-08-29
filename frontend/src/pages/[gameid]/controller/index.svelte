<script>
    import api from '../../../utils/api';
    import X01Controller from './X01Controller.svelte';
    import CricketController from './CricketController.svelte';
    import AtcController from './ATCController.svelte';
    import SplitController from './SplitController.svelte';
    import ShanghaiController from './ShanghaiController.svelte';

    export let gameid;
    let headerdiv = document.getElementsByClassName('header')[0];

    // Hide navbar in this page
    headerdiv.style.display = 'none';

    const getGame = async () => {
        let gameMode;

        const res = await api.get(`game/${gameid}`);
        let game = await res.json();

        switch (game.Game) {
            case 'x01':
                gameMode = X01Controller;
                break;

            case 'cricket':
                gameMode = CricketController;
                break;

            case 'atc':
                gameMode = AtcController;
                break;

            case 'split':
                gameMode = SplitController;
                break;

            case 'shanghai':
                gameMode = ShanghaiController
                break;

            default:
                throw new Error('Controller cannot be shown');
        }
        return gameMode;
    };

    let promise = getGame();
</script>

{#await promise}
    <p>...loading</p>
{:then game}
    <svelte:component this={game} {gameid} />
{:catch error}
    <p>There was an error {error.message}</p>
{/await}
