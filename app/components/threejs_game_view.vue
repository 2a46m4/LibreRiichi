<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { use_websocket_state } from '..'
import ScoreBoard from '../components/scoreboard.vue'
import { GameIdx, TableIdx } from '../game/arena'
import { Setup, SetupType } from '../game/setup'
import { Tile } from '../game/tile'
import { Action, ActionType } from '../messaging/action_generated'
import { ArenaEvent, ArenaEventType } from '../messaging/arena_event_generated'
import { BoardEvent, BoardEventType } from '../messaging/board_event_generated'
import {
  ArenaMessageBus,
  ClickBus,
  make_async_generator_from_event,
  register_request,
  select,
} from '../messaging/event_handler'
import { send_action_request } from '../messaging/message'
import { ServerEvent } from '../messaging/server_event_generated'
import { ThreeJSRenderer } from '../render/renderer'
import { initialize_tiles, TileObject } from '../render/tile'

import { ServerResponseType } from '../messaging/server_response_generated'

const props = defineProps<{ in_game: boolean }>()
const websocket_state = use_websocket_state()

const arena_data = ref({
  dealer_idx: 0,
  round_wind: 0,
  round_number: 0,
  players: [],
  player_idx: 0 as GameIdx, // Game idx
  scores: [] as number[],
  game_should_end: false,
  round_should_end: false,
  // Hacky, but we need to keep track of this to remember which tile to discard when the player event comes
  // A better way could be to remove when the server returns a positive result to our request, although that
  // 	also means that we need to ignore our result when the server event comes.
  selected: null as TileObject | null,
})

// Self is table idx 0
// So if we are player 3, player 0's offset is (4 + 0 - 3) % 4 = 1
function offset_to_self(other_index: GameIdx): TableIdx {
  return ((4 + arena_data.value.player_idx + other_index) % 4) as TableIdx
}

const three_canvas = ref<HTMLCanvasElement>()
const game_container = ref<HTMLDivElement>()
const is_fullscreen = ref(false)

let animation_id: number
let renderer: ThreeJSRenderer
let click_channel = make_async_generator_from_event(ClickBus)
let server_event_channel = make_async_generator_from_event(ArenaMessageBus)

onMounted(() => {
  if (!three_canvas.value) return
  initialize_tiles()

  renderer = new ThreeJSRenderer(three_canvas.value)
  animate(0, 0)

  window.addEventListener('resize', on_window_resize)
})

function on_window_resize() {
  if (!three_canvas.value || !game_container.value) return

  const width = game_container.value.clientWidth
  const height = game_container.value.clientHeight

  three_canvas.value.width = width
  three_canvas.value.height = height

  renderer.window_resize(width, height)
}

function animate(t: number, dt: number) {
  renderer.animate_frame(dt)
  animation_id = requestAnimationFrame((new_t) => {
    animate(new_t, new_t - t)
  })
}

onUnmounted(() => {
  if (animation_id) {
    cancelAnimationFrame(animation_id)
  }
  window.removeEventListener('resize', on_window_resize)

  renderer.stop()
})

function handle_player_action_event(action: Action, from_player: number) {
  switch (action.action_type) {
    case ActionType.Tsumo:
      // TODO
      throw new Error('Win')
    case ActionType.Ron:
      // TODO
      throw new Error('Win')
    case ActionType.Riichi: // Show the toss animation, handle riichi case
    case ActionType.Toss: // Special case if we are the one that tossed the tile
      const tile_value =
        action.action_type === ActionType.Riichi
          ? new Tile(action.tile_to_riichi)
          : new Tile(action.tile_to_toss)
      let tile_obj: TileObject
      if (from_player === arena_data.value.player_idx) {
        if (arena_data.value.selected === null) {
          throw new Error('Selection is null')
        }
        tile_obj = arena_data.value.selected
      } else {
        tile_obj = renderer.get_random_tile_in_hand(from_player)
      }
      tile_obj.value = tile_value
      renderer.toss(
        from_player,
        tile_obj,
        action.action_type === ActionType.Riichi,
      )
      break
    case ActionType.Pon:
      renderer.naki_call(from_player as TableIdx, 'pon')
    case ActionType.Kan:
      renderer.naki_call(from_player as TableIdx, 'ankan')
    case ActionType.Chii:
      renderer.naki_call(from_player as TableIdx, 'chii')
      break
    case ActionType.Draw:
      renderer.add_tile_to_end(new Tile(action.drawn_tile), from_player)
      break
    default:
      console.error('Unexpected action performed')
  }
}

function handle_potential_action_event(actions: Action[]) {
  for (let action of actions) {
    switch (action.action_type) {
      case ActionType.Tsumo:
        console.log('Tsumo possible')
        break
      case ActionType.Ron:
        console.log('Ron possible')
        break
      case ActionType.Riichi:
        console.log('Riichi possible')
        break
      case ActionType.Toss:
        break
      case ActionType.Skip:
        // need to show the player a prompt to skip
        console.log('Skip possible')
        break
      case ActionType.Pon:
        console.log('Pon possible')
        break
      case ActionType.Kan:
        console.log('Kan possible')
        break
      case ActionType.Chii:
        console.log('Chii possible')
        break
      case ActionType.Draw:
        throw new Error('Unexpected draw')
        break
    }
  }
}

function handle_game_setup_event(setups: Setup[]) {
  for (let setup of setups) {
    switch (setup.setup_type) {
      case SetupType.INITIAL_TILES:
        renderer.clear_tiles()
        Tile.from(setup.data)
          .sort(Tile.sort)
          .map((t) => renderer.add_tile_to_end(t))
        break
      case SetupType.DORA:
        renderer.add_dora(new Tile(setup.data))
        break
      case SetupType.STARTING_POINTS:
        arena_data.value.scores = setup.data
        break
      case SetupType.PLAYER_NUMBER:
        console.log('Player number setup received: ', setup.data)
        arena_data.value.player_idx = setup.data as GameIdx
        break
      case SetupType.PLAYER_ORDER:
        throw new Error('Should be removed')
      case SetupType.ROUND_WIND:
        arena_data.value.round_wind = setup.data
        break
      case SetupType.ROUND_NUMBER:
        arena_data.value.round_number = setup.data
        break
    }
  }
}

function handle_arena_event(arena_message: ArenaEvent) {
  switch (arena_message.arenaevent_type) {
    case ArenaEventType.ArenaBoardEvent:
      handle_board_event(arena_message.board_event)
      break
    default:
      console.log(`Not handling event: ${arena_message}`)
  }
}

function handle_board_event(new_event: BoardEvent) {
  if (new_event === undefined) {
    throw new Error('Event is undefined')
  }

  switch (new_event.boardevent_type) {
    case BoardEventType.PlayerActionEvent:
      handle_player_action_event(new_event.action_data, new_event.from_player)
      break
    case BoardEventType.PotentialActionEvent:
      handle_potential_action_event(new_event.actions)
      break
    case BoardEventType.GameSetupEvent:
      handle_game_setup_event(new_event.setup)
      break
    case BoardEventType.GameEndEvent:
      // Handle game end event
      throw new Error('TODO')
      break
    default:
      throw new Error('Unexpected')
  }
}

function handle_mouse_event(_: MouseEvent) {
  const selected = renderer.selector.get_selection()

  if (selected === null) {
    return
  }

  let msg_idx = websocket_state.conn.send(
    send_action_request({
      action_type: ActionType.Toss,
      tile_to_toss: selected.value.value,
    }),
  )

  register_request(msg_idx)
    .then((response) => {
      switch (response.serverresponse_type) {
        case ServerResponseType.GenericResponse:
          if (!response.success) {
            console.error(`Couldn't satisfy request: ${response.fail_reason}`)
          } else {
            console.log('Success in discarding tile')
          }
        default:
          console.error('Unexpected response')
      }
    })
    .catch(() => {
      console.error('Discard action was not acknowledged by server')
    })
}

const unvoid = function <T>(x: T | void): T {
  if (!x) throw new Error('Void')
  else return x
}

while (!arena_data.value.game_should_end) {
  const data = await server_event_channel.next()
  if (!data.value) {
    throw new Error('Server event returned void')
  }

  while (!arena_data.value.round_should_end) {
    const selection = select(server_event_channel, click_channel)
    const _msg = await selection.next()
    const msg = unvoid(_msg.value)
    if (msg.index == 0) {
      handle_arena_event((msg.value as ServerEvent).arena_message)
    } else {
      handle_mouse_event(msg.value as MouseEvent)
    }
  }

  // Game end event
}
</script>

<template>
  <ScoreBoard
    v-if="in_game"
    :scoreboard_values="arena_data.scores"
    :round_wind="arena_data.round_wind"
    :round_number="arena_data.round_number"
    :players="arena_data.players"
    :player_idx="arena_data.player_idx"
  >
  </ScoreBoard>
  <div
    ref="game_container"
    class="game-view-container"
    :class="{ fullscreen: is_fullscreen }"
  >
    <canvas ref="three_canvas" class="game-canvas"></canvas>
    <div class="game-ui"></div>
  </div>
</template>

<style scoped>
.game-view-container {
  position: relative;
  width: 100%;
  height: 100%;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
}

.game-canvas {
  width: 100%;
  height: 100%;
  display: block;
}

.game-ui {
  position: absolute;
  top: 16px;
  left: 16px;
  z-index: 10;
}

.game-info {
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(10px);
  padding: 12px 16px;
  border-radius: 8px;
  color: white;
  font-size: 14px;
}

.game-info h3 {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: 600;
}

.game-info p {
  margin: 0;
  opacity: 0.8;
  font-size: 12px;
}
</style>
