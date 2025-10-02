<script setup lang="ts">

import {AgentInfo, sort_agents} from "../game/agent_info";
import {watch} from "vue";

const props = defineProps<{
  scoreboard_values: number[],
  player_to_order_map: number[],
  round_wind: number,
  round_number: number,
  players: AgentInfo[],
  player_idx: number
}>()

watch(props.players, sort_agents.bind(props.players))

function get_value(player_idx: number) {
  return props.scoreboard_values[player_idx]
}

</script>

<template>
  <div class="fixed top-2 right-2 outline p-2 bg-white rounded shadow-md text-black z-10">
    <h1 class="text-xl">Scores</h1>
    <p>You: {{get_value(props.player_idx)}}</p>
    <p></p>
  </div>
</template>

<style scoped></style>
