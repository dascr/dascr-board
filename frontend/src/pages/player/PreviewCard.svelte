<script>
  import previewPlayer from '../../utils/stores/playerStore';
  import Cropper from 'svelte-easy-crop';
  import PlayerCard from './PlayerCard.svelte';

  let apiBase = 'API_BASE';

  // Cropper options
  const options = {
    crop: { x: 0, y: 0 },
    zoom: 1,
    aspect: 1 / 1,
    minZoom: 0,
    maxZoom: 3,
    cropShape: 'rect',
    cropSize: { width: 300, height: 300 },
    showGrid: true,
    zoomSpeed: 0.1,
    crossOrigin: true,
    restrictPosition: false,
  };

  const handleCropComplete = (e) => {
    let cropFactor = { ...e.detail.pixels };
    previewPlayer.setCropFactor(cropFactor);
  };

  $: name = $previewPlayer.name || 'Max';
  $: nickname = $previewPlayer.nickname || 'The Demo Player';
  $: image = $previewPlayer.image || `${apiBase}images/na.png`;
</script>

<style>
  .cropperWrapper {
    position: relative;
    height: 300px;
    width: 300px;
    background-color: white;
  }
</style>

<PlayerCard uid="0" {name} {nickname} showDelete={false} onDelete={() => {}}>
  <div slot="avatar">
    <div class="cropperWrapper">
      {#if image != `${apiBase}images/na.png`}
        <!--Croppy thingy-->
        <Cropper {image} {...options} on:cropcomplete={handleCropComplete} />
      {:else}
        <!--Placeholder-->
        <img width="300" height="300" src={image} alt="previewPlaceholder" />
      {/if}
    </div>
  </div>
</PlayerCard>
