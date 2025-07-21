<script setup lang="ts">

import {BoxStyling, ButtonStyling, FlexBox, H1Styling, Spacing, ULStyling} from "../styling";
import {useGlobalStore} from "../global_store";
import {ref, Ref} from "vue";
import ListItem from "../components/list_item.vue";
import {ArenaMessage, ArenaMessageType} from "../messaging/arena_message";

const store = useGlobalStore()
const app = store.application
const action = app.action
const handler = app.handler

const players: Ref<string[]> = ref([])
const room_name = ref('')

let in_game = false

async function get_arena_info() {
  let arena = await action.get_arena_info()
  players.value = arena.agents.map(x => x.name)
  room_name.value = arena.name
}

await get_arena_info()
let listener_idx = handler.register_arena_listener((data: ArenaMessage) => {
  console.log("Arena listener called")
  switch (data.message_type) {
    case ArenaMessageType.PlayerJoinedEvent:
      players.value.push(data.data.name);
      break;
    case ArenaMessageType.PlayerQuitEvent:
        players.value = players.value.filter((v) => v !== data.data.name)
      break;
    case ArenaMessageType.GameStartedEvent:
      in_game = true
      break;
    case ArenaMessageType.ArenaBoardEvent:
      if (!in_game) {
        throw new Error("Not in correct state")
      } else {
        throw new Error("NYI")
      }
      break;
    default:
      throw new Error("Unexpected message")
  }
})

function start_game() {
  action.
}

</script>

<template>
  <div :class="BoxStyling">
    <h1 :class="H1Styling">{{ room_name }}</h1>
  </div>
  <div :class="BoxStyling">
    <div :class="FlexBox">
      <h1 :class="H1Styling">Players {{ players.length }} / 4</h1>
    </div>
    <ul :class="ULStyling" v-if="players.length !== 0">
      <ListItem v-for="player in players">{{ player }}</ListItem>
    </ul>
  </div>
  <button :class="ButtonStyling + Spacing"
          @click="start_game">Start game</button>
</template>

<style scoped>

</style>