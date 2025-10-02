<script setup lang="ts">

const props = defineProps<{
  scoreboard_values: number[],
  player_to_order_map: number[],
  round_wind: number,
  round_number: number,
  player_names: number[],
}>()

</script>

<template>
  <div class="fixed top-2 right-2 outline p-2 bg-white rounded shadow-md text-black z-10">
    <h1 class="text-xl">Scores</h1>
    <p>You: {{}}</p>
    <p></p>
  </div>
</template>

<style scoped></style>
