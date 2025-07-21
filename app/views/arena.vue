<script setup lang="ts">

import {BoxStyling, ButtonStyling, FlexBox, H1Styling, Spacing, ULStyling} from "../styling";
import {useGlobalStore} from "../global_store";
import {ref, Ref} from "vue";
import ListItem from "../components/list_item.vue";
import {ArenaMessage, ArenaMessageType} from "../messaging/arena_message";
import GameBoard from "../components/game_board.vue";

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

get_arena_info()
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

async function start_game() {
  await action.start_game()
}

</script>

<template>
  <Suspense>
    <div>
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
      <div v-if="in_game">
        <GameBoard></GameBoard>
      </div>
    </div>
  </Suspense>
</template>

<style scoped>

</style>