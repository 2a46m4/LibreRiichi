<script setup lang="ts">

import {BoxStyling, ButtonStyling, FlexBox, H1Styling, Spacing, ULStyling} from "../styling";
import {ref, Ref} from "vue";
import ListItem from "../components/list_item.vue";
import GameBoard from "../components/game_board.vue";
import {router, use_room_state, use_websocket_state} from "../index";
import {MessageType} from "../messaging/message";
import {ServerActionType} from "../messaging/server_action_generated";
import {ArenaMessageBus, register_request} from "../messaging/event_handler";
import {ServerResponseType} from "../messaging/server_response_generated";
import {ServerEventMessage, ServerEventType} from "../messaging/server_event_generated";
import {ArenaEventType} from "../messaging/arena_event_generated";
import {ArenaActionType} from "../messaging/arena_action_generated";

const players: Ref<string[]> = ref([])

const room_state = use_room_state()
if (!room_state.room_set) {
  throw new Error("Room not set")
}

const websocket_state = use_websocket_state()
const error_status = ref('')

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

  let ret = await register_request(msg_idx)
  if (ret.serverresponse_type !== ServerResponseType.ArenaInfoResponse) {
    error_status.value = "Connection error: wrong type"
    return
  }

  if (!ret.success) {
    error_status.value = "Could not get arena data"
    return
  }

  players.value = ret.agents.map(x => x.name)
  room_state.room_name = ret.name
}

get_arena_info()

let callback = (data: ServerEventMessage) => {
  console.log("Arena listener called")
  switch (data.serverevent_type) {
    case ServerEventType.ServerArenaEvent:
      let message = data.arena_message
      switch (message.arenaevent_type) {
        case ArenaEventType.GameStartedEvent:
          in_game = true
          ArenaMessageBus.unregister(callback_idx)
		  // TODO
          break;
        case ArenaEventType.PlayerJoinedEvent:
          players.value.push(message.name);
          break;
        case ArenaEventType.PlayerQuitEvent:
          players.value = players.value.filter((v) => v !== message.name)
          break;
        case ArenaEventType.ArenaBoardEvent:
          if (!in_game) {
            throw new Error("Game not started")
          }
          break;
        default:
          throw new Error("Unknown arena event")
      }
      break;
    default:
      throw new Error("Unknown server event")
  }
  return true
}
let callback_idx = ArenaMessageBus.register(callback)

async function start_game() {
  let msg_idx = websocket_state.conn.send(
      {
        message_type: MessageType.REQUEST,
        data: {
          serveraction_type: ServerActionType.ServerArenaAction,
          arena_action: {
            arenaaction_type: ArenaActionType.StartGameActionData
          }
        }
      }
  )

  let ret = await register_request(msg_idx);
  if (ret.serverresponse_type !== ServerResponseType.GenericResponse) {
    throw new Error("Connection error: Wrong Type")
  }

  if (!ret.success) {
    throw new Error("Couldn't start game: " + ret.fail_reason)
  }
}

</script>

<template>
  <Suspense>
    <div>
      <div :class="BoxStyling">
        <h1 :class="H1Styling">{{ room_state.room_name }}</h1>
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
              @click="start_game">Start game
      </button>
      <div v-if="in_game">
        <GameBoard></GameBoard>
      </div>
    </div>
  </Suspense>
</template>

<style scoped>

</style>
