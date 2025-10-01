<script setup lang="ts">
import TopBoxElement from './top_box_element.vue'

const props = defineProps<{
  player_number: number,
  // Maps player index to order
  player_ordering: number[],
  round_wind: number,
  round_number: number,
  starting_points: number[],
}>()

function get_ordering() {
  return props.player_ordering[props.player_number]
}
</script>

<template>
  <TopBoxElement :text="'Player' + player_number"/>
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
