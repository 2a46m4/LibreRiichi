<script setup lang="ts">

import {BoxStyling, ButtonStyling, FlexBox, H1Styling, Spacing, ULStyling} from "../styling";
import {ref, Ref} from "vue";
import ListItem from "../components/list_item.vue";
import {ArenaMessage, ArenaMessageType} from "../messaging/arena_message";
import GameBoard from "../components/game_board.vue";
import {use_room_state, use_websocket_state} from "../index";
import {MessageType} from "../messaging/message";
import {ServerActionType} from "../messaging/server_action_generated";

const players: Ref<string[]> = ref([])

const room_state = use_room_state()
if (!room_state.room_set) {
  throw new Error("Room not set")
}

const room_name = room_state.room_name

const websocket_state = use_websocket_state()

let in_game = false

async function get_arena_info() {
  let msg_idx = websocket_state.conn.send(
      {
        message_type: MessageType.REQUEST,
        data: {
          serveraction_type: ServerActionType.ArenaInfoAction,
        }
      }
  )

  let ret = await websocket_state.msg_router.register_message(msg_idx)
  if (ret.message_type !== MessageType.ArenaInfoResponse) {
    throw new Error("Connection error: wrong type")
  }

  if (!ret.data.success) {
    throw new Error("Could not get arena data")
  }

  return ret.data
  players.value = arena.agents.map(x => x.name)
  room_name.value = arena.name
}

await get_arena_info()
let callback = (data: ArenaMessage) => {
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
        throw new Error("Game not started")
      }


      break;
    default:
      throw new Error("Unexpected message")
  }
}
let arena_handler = new ArenaHandler(handler)
let callback_idx = arena_handler.register_arena_listener(callback)

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