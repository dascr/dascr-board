<!--This page need design -->
<script>
  import api from '../../utils/api';
  export let gameid;

  async function fetchGame() {
    const res = api.get(`game/${gameid}`);
    const game = await res.json();
    return game;
  }

  let promise = fetchGame();
</script>

<h1 class="text-center font-extrabold text-2xl mb-4">
  Overview of game
  {gameid}
</h1>

{#await promise}
  <p>... fetching game details</p>
{:then game}
  <p>Game: {game.Game}</p>
  <p>Variant: {game.Variant}</p>
  {#if game.In != ''}
    <p>In: {game.In}</p>
  {/if}
  {#if game.Out != ''}
    <p>Out: {game.Out}</p>
  {/if}
  <p>Throw Round: {game.ThrowRound}</p>
  <p>Game State: {game.GameState}</p>
  <p>Message: {game.Message}</p>
  <p>Settings</p>
  <ul>
    <li>Sound: {game.Settings.Sound}</li>
    <li>Podium: {game.Settings.Podium}</li>
  </ul>
  <p>Player</p>
  <ul>
    {#each game.Player as player, i}
      <li>
        {i + 1}:
        {player.Name}
        {#if player.Nickname != ''}
          -
          {player.Nickname}:
          {player.Score.Score}
        {/if}
      </li>
    {/each}
  </ul>
{:catch}
  <ul>There is no game running with the id: <i>{gameid}</i></ul>
{/await}
