<script setup lang="ts">
  const props = defineProps<{
    tile_path: URL
  }>()
</script>

<template>
  <img :src=props.tile_path.toString() alt="A tile"/>
</template>

<style scoped>

</style>