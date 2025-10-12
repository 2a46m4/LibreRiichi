<script setup lang="ts">
import { ref, Ref } from 'vue'
import BoxElement from '../components/box_element.vue'
import TitleBoxElement from '../components/title_box_element.vue'
import ThreeJSGameView from '../components/threejs_game_view.vue'
import List from '../components/list.vue'
import { use_room_state, use_websocket_state } from '../index'
import { MessageType } from '../messaging/message'
import { ServerActionType } from '../messaging/server_action_generated'
import { ArenaMessageBus, register_request } from '../messaging/event_handler'
import { ServerResponseType } from '../messaging/server_response_generated'
import {
  ServerEvent,
  ServerEventType,
} from '../messaging/server_event_generated'
import { ArenaEventType } from '../messaging/arena_event_generated'
import { ArenaActionType } from '../messaging/arena_action_generated'
import { BoardEvent, BoardEventType } from '../messaging/board_event_generated'
import { SetupType } from '../game/setup'
import Button from '../components/button.vue'
import {Arena} from "../game/arena";

const num_ai: Ref<number> = ref(0)
const error_status = ref('')
const in_game = ref(false)
const arena: Ref<Arena> = ref({
  agents: [],
  game_started: false
})
const events: BoardEvent[] = []

const room_state = use_room_state()
if (!room_state.room_set) {
  error_status.value = 'Room not set'
}

const websocket_state = use_websocket_state()

async function get_arena_info() {
  let msg_idx = websocket_state.conn.send({
    message_type: MessageType.REQUEST,
    data: {
      serveraction_type: ServerActionType.ArenaInfoAction,
    },
  })

  let ret = await register_request(msg_idx)
  if (ret.serverresponse_type !== ServerResponseType.ArenaInfoResponse) {
    error_status.value = 'Connection error: wrong type'
    return
  }

  if (!ret.success) {
    error_status.value = 'Could not get arena data'
    return
  }

  arena.value.agents = ret.agents
  room_state.room_name = ret.name
}

get_arena_info()

let callback = (data: ServerEvent) => {
  console.log('Arena listener called: ', data)
  switch (data.serverevent_type) {
    case ServerEventType.ServerArenaEvent:
      let message = data.arena_message
      switch (message.arenaevent_type) {
        case ArenaEventType.GameStartedEvent:
          in_game.value = true
          break
        case ArenaEventType.PlayerJoinedEvent:
          arena.value.agents.push(message.agent_info)
          break
        case ArenaEventType.PlayerQuitEvent:
          arena.value.agents = arena.value.agents.filter((v) => v.name !== message.name)
          break
        case ArenaEventType.ArenaBoardEvent:
          if (!in_game.value) {
            error_status.value = 'Game not started'
            return true
          }
          break
        default:
          error_status.value = 'Unknown arena event'
          return true
      }
      break
    default:
      error_status.value = 'Unknown server event'
      return true
  }
  return true
}
let callback_idx = ArenaMessageBus.register(callback)
</script>

<template>
  <Suspense>
    <div>
      <div class="game-container">
        <div v-if="!in_game" id="ui">
          <TitleBoxElement :text="'Room name: ' + room_state.room_name" />
          <div class="container outline bg-white rounded shadow-md pb-5 mb-5">
            <h1 class="font-bold text-xl text-center pt-2">
              Players {{ arena.agents.length }} / 4
            </h1>
            <List :items="arena.agents.map(a=>a.name)" class="p-1"></List>
          </div>

          <Button text="Start game" :condition="!in_game" :on_click="async () => {
            let msg_idx = websocket_state.conn.send({
              message_type: MessageType.REQUEST,
              data: {
                serveraction_type: ServerActionType.ServerArenaAction,
                arena_action: {
                  arenaaction_type: ArenaActionType.StartGameActionData,
                },
              },
            })

            let ret = await register_request(msg_idx)
            if (
              ret.serverresponse_type !== ServerResponseType.GenericResponse
            ) {
              error_status = 'Connection error: Wrong Type'
              return
            }

            if (!ret.success) {
              error_status = 'Couldn\'t start game: ' + ret.fail_reason
              return
            }

            in_game = true
          }
            " />

          <Button text="Add AI" :condition="!in_game" :on_click="async () => {
            let msg_idx = websocket_state.conn.send({
              message_type: MessageType.REQUEST,
              data: {
                serveraction_type: ServerActionType.ServerArenaAction,
                arena_action: {
                  arenaaction_type: ArenaActionType.AddAIArenaAction,
                },
              },
            })

            let ret = await register_request(msg_idx)
            if (
              ret.serverresponse_type !== ServerResponseType.GenericResponse
            ) {
              error_status = 'Connection error: Wrong Type'
              return
            }

            if (!ret.success) {
              error_status = 'Couldn\'t add AI: ' + ret.fail_reason
            }

            num_ai = num_ai + 1
          }
            " />

          <Button text="Remove AI" :condition="num_ai > 0" :on_click="async () => {
            let msg_idx = websocket_state.conn.send({
              message_type: MessageType.REQUEST,
              data: {
                serveraction_type: ServerActionType.ServerArenaAction,
                arena_action: {
                  arenaaction_type: ArenaActionType.RemoveAIArenaAction,
                },
              },
            })

            let ret = await register_request(msg_idx)
            if (
              ret.serverresponse_type !== ServerResponseType.GenericResponse
            ) {
              throw new Error('Connection error: Wrong Type')
            }

            if (!ret.success) {
              throw new Error('Couldn\'t remove AI: ' + ret.fail_reason)
            }

            num_ai = num_ai - 1
          }
            " />
          <BoxElement :text="error_status" v-if="error_status.length !== 0" />
        </div>
        <ThreeJSGameView :events="events"
        :in_game="in_game"/>
      </div>
    </div>
  </Suspense>
</template>

<style scoped>
.game-container {
  position: relative;
  width: 100%;
  height: 100%;
}

#ui {
  position: absolute;
  /* let us position ourself inside the container */
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  z-index: 100;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}
</style>
