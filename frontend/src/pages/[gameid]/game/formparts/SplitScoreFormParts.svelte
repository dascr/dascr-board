<script>
  import setupGame from '../../../../utils/stores/gameStore';
</script>

<div class="flex flex-col">
  <label class="uppercase font-bold text-lg" for="selectVariant">Select Variant:</label>
  <select
    bind:value={$setupGame.variant}
    class="border py-2 px-3 text-gray-900"
    id="selectVariant"
    name="variant">
    <option value="steel">Steel Dart</option>
    <option value="edart">E-Dart</option>
  </select>
</div>
