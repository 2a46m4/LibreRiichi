<script setup lang="ts">

import { AgentInfo } from "../game/agent_info";
import { Tile as TileValue } from "../messaging/tile";

const props = defineProps<{
  scoreboard_values: number[],
  round_wind: number,
  round_number: number,
  players: AgentInfo[],
  player_idx: number,
  debug_info?: string,
}>()

function round_wind() {
  switch (props.round_wind) {
    case TileValue.EastTile:
      return "East"
    case TileValue.SouthTile:
      return "South"
    case TileValue.WestTile:
      return "West"
    case TileValue.NorthTile:
      return "North"
    default:
      return "Unknown"
  }
}

</script>

<template>
  <div class="fixed top-2 right-2 outline p-2 bg-white rounded shadow-md text-black z-10">
    <h1 class="text-xl">Scores</h1>
    <p>{{scoreboard_values
      .map(n => n.toString())
      .join(", ")}}</p>
    <h1 class="text-xl">Round wind:</h1>
    <p>{{ round_wind() }}</p>
    <h1 class="text-xl">Player index:</h1>
    <p>{{ player_idx }}</p>
    <div v-if="debug_info !== undefined">
      <h1 class="text-xl">Debug info:</h1>
      <p>{{ debug_info }}</p>
    </div>

  </div>
</template>

<style scoped></style>
