<script>
  import { onMount } from 'svelte';

  import previewPlayer from '$stores/playerStore';

  export let onSubmit;

  let files;

  $: {
    files && previewPlayer.setImage(files[0]);
  }
  $: cropFactor = $previewPlayer.cropfactor || {};

  onMount(() => {
    // Eventlistener für image change
    var source = document.getElementById('playerImage');
    source.addEventListener('change', () => {
      // Write image to source img tag
      if (source.files && source.files[0]) {
        var fr = new FileReader();
        fr.onload = () => {
          document.getElementById('source').src = fr.result;
        };
        fr.readAsDataURL(files[0]);
      }
    });
  });

  const transformImage = () => {
    var canvas = document.getElementById('canvas');
    var context = canvas.getContext('2d');
    var img = document.getElementById('source');
    context.drawImage(
      img,
      cropFactor.x,
      cropFactor.y,
      cropFactor.width,
      cropFactor.height,
      0,
      0,
      300,
      300
    );

    return canvas.toDataURL('image/png');
  };

  const handleSubmit = (e) => {
    const { name, nickname, image } = e.target;

    let resultImage;

    // Process crop image
    if (image.files && image.files[0]) {
      // Crop
      resultImage = transformImage();
    } else {
      resultImage = null;
    }

    let player = {
      name: name.value,
      nickname: nickname.value,
      image: resultImage,
    };
    let frm = document.getElementById('newPlayer');
    frm.reset();
    onSubmit(player);
  };
</script>

<!-- Form col -->
<h1 class="text-2xl text-center">Add Player</h1>
<form
  id="newPlayer"
  enctype="multipart/form-data"
  autocomplete="off"
  on:submit|preventDefault={handleSubmit}>
  <div class="space-y-4">
    <div class="flex flex-col">
      <label for="playerName" class="uppercase font-bold text-lg">Name
        <small class="lowercase">*required</small></label>
      <input
        bind:value={$previewPlayer.name}
        type="text"
        class="border py-2 px-3 text-gray-900"
        id="playerName"
        name="name"
        placeholder="Enter player name"
        required />
    </div>
    <div class="flex flex-col">
      <label
        for="playerNickname"
        class="uppercase font-bold text-lg">Nickname</label>
      <input
        bind:value={$previewPlayer.nickname}
        type="text"
        class="border py-2 px-3 text-gray-900"
        id="playerNickname"
        name="nickname"
        placeholder="Enter player nickname"
        maxlength="12" />
    </div>
    <div class="flex flex-col">
      <label for="playerImage" class="text-center"><span
          class="block uppercase border-2  p-4 rounded-2xl font-bold text-lg hover:bg-black hover:bg-opacity-30"><i
            class="fas fa-upload" />
          upload avatar
          <input
            bind:files
            type="file"
            class="hidden"
            id="playerImage"
            accept="image/*"
            name="image" /></span></label>
    </div>
    <button
      type="submit"
      class="block uppercase border-2 hover:bg-black hover:bg-opacity-30 text-lg mx-auto p-4 rounded-2xl"><i
        class="fas fa-plus" />
      Add Player</button>
  </div>
</form>

<img class="hidden" id="source" src="" alt="" />
<canvas class="hidden" id="canvas" width="300" height="300" />
