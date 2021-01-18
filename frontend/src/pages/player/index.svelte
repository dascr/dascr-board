<script>
  import api from '../../utils/api';
  import PlayerForm from './PlayerForm.svelte';
  import PlayerCard from './PlayerCard.svelte';
  import { toast } from '@zerodevx/svelte-toast';
  import { onMount } from 'svelte';
  import PreviewCard from './PreviewCard.svelte';
  import previewPlayer from '../../utils/stores/playerStore';

  let availablePlayer = [];

  onMount(async () => {
    const res = await api.get('player');
    availablePlayer = await res.json();
  });

  async function handleSubmit(player) {
    await api
      .post('player', {
        json: {
          name: player.name,
          nickname: player.nickname,
        },
      })
      .then(async (res) => {
        const data = await res.json();
        res = api
          .post(`player/${data.UID}/image`, {
            json: { b64: player.image },
          })
          .then(async () => {
            res = api.get('player');
            availablePlayer = await res.json();
            previewPlayer.resetAll();
            toast.push('Player was created, avatar was uploaded.', {
              theme: {
                '--toastBackground': '#48BB78',
                '--toastProgressBackground': '#2F855A',
              },
            });
          })
          .catch((err) => {
            console.error(err);
            toast.push('Player was created, avatar was not uploaded.', {
              theme: {
                '--toastBackground': '#ECC94B',
                '--toastProgressBackground': '#B7791F',
              },
            });
          });
      })
      .catch((err) => {
        console.error(err);
        toast.push('Player was not created', {
          theme: {
            '--toastBackground': '#F56565',
            '--toastProgressBackground': '#C53030',
          },
        });
      });
  }

  async function handleDelete(id) {
    api
      .delete(`player/${id}`)
      .then(() => {
        const newAvailablePlayer = availablePlayer.filter(
          (player) => player.UID !== id
        );
        availablePlayer = newAvailablePlayer;
        toast.push('Player was deleted', {
          theme: {
            '--toastBackground': '#48BB78',
            '--toastProgressBackground': '#2F855A',
          },
        });
      })
      .catch((err) => {
        console.log(err);
        toast.push('Player was not deleted', {
          theme: {
            '--toastBackground': '#F56565',
            '--toastProgressBackground': '#C53030',
          },
        });
      });
  }
</script>

<div class="max-w-full px-4">
  <div class="flex flex-wrap">
    <div class="w-full p-4 lg:w-1/2">
      <PlayerForm onSubmit={handleSubmit} />
    </div>

    <div class="w-full p-4 lg:w-1/2">
      <PreviewCard />
    </div>
  </div>
</div>

<div>
  <hr class="mb-2 p-4" />
</div>

<div class="max-w-full px-4">
  <h1 class="text-2xl text-center">Available Player</h1>
  <div class="flex flex-wrap">
    {#each availablePlayer || [] as player}
      <div class="w-full p-4 lg:w-1/2">
        <PlayerCard
          uid={player.UID}
          name={player.Name}
          nickname={player.Nickname}
          image={player.Image}
          onDelete={handleDelete}
          showDelete={true} />
      </div>
    {:else}
      <h2 class="text-center text-lg">
        There are no player, yet. Go ahead, create one!
      </h2>
    {/each}
  </div>
</div>
