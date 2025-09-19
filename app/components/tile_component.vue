<script setup lang="ts">
import Tile from "./tile_component.vue";

  const props = defineProps<{
    tile_paths: URL[]
  }>()

</script>

<template>
  <img v-for="tile_path in props.tile_paths" :src=tile_path.toString() alt="A tile"/>
</template>

<style scoped>
div {
  position: relative;
  display: inline;
}
img {
  position: absolute;
}
</style>