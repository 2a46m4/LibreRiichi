<script setup lang="ts">
import TopBoxElement from './top_box_element.vue'
import {AgentInfo, sort_agents} from "../game/agent_info";
import {watch} from "vue";

const props = defineProps<{
  scoreboard_values: number[]
  // Maps player index to order
  player_to_order_map: number[],
  round_wind: number,
  round_number: number,
  players: AgentInfo[]
  player_idx: number,
}>()

watch(props.players, sort_agents.bind(props.players))

function get_ordering() {
  return props.player_to_order_map[props.player_idx]
}
</script>

<template>
  <TopBoxElement :player_idx="player_idx"
                 :player_to_order_map="player_to_order_map"
                 :round_wind="round_wind"
                 :round_number="round_number"
                 :scoreboard_values="scoreboard_values"
                 :players="players"
  />
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
