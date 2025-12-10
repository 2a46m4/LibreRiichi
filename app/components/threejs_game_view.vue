<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { HiddenTile, Tile } from '../game/tile'
import { initialize_tiles } from '../render/tile'
import { BoardEvent, BoardEventType } from '../messaging/board_event_generated'
import { Setup, SetupType } from '../game/setup'
import { ArenaMessageBus, register_request } from "../messaging/event_handler";
import { ArenaEventType } from "../messaging/arena_event_generated";
import { ServerEvent, ServerEventType } from "../messaging/server_event_generated";
import ScoreBoard from "../components/scoreboard.vue"
import { Action, ActionType } from "../messaging/action_generated";
import { IActionAnimator, IRenderer, ISelectionManager, Scene } from "../render/renderer";
import { GameIdx, TableIdx } from '../game/arena'
import { create_event, create_fsm_builder, create_state } from "../fsm";
import { use_websocket_state } from '..'
import { MessageType } from '../messaging/message'
import { ServerActionType } from '../messaging/server_action_generated'
import { ArenaActionType } from '../messaging/arena_action_generated'

// Data flow in this file:
// Event comes in from the server -> message_handler
// Message handler dispatches the correct state transition in the FSM
// State transition callbacks handle the animations and game state updates

const props = defineProps<{ in_game: boolean }>()
const websocket_state = use_websocket_state()

interface Add {
  tag: "add",
  add: Tile,
  location: number,
  table_idx: TableIdx,
}
interface Remove {
  tag: "remove",
  location: number,
  table_idx: TableIdx,
  riichi_called: boolean
}
interface NewSet {
  tag: "newset",
  set: Tile[],
  table_idx: TableIdx,
}
interface NakiCall {
  tag: "nakicall",
  call: Action,
  table_idx: TableIdx,
}
type Animation = Add | Remove | NewSet | NakiCall

const arena_data = ref({
  dealer_idx: 0,
  round_wind: 0,
  round_number: 0,
  players: [],
  player_idx: 0 as GameIdx, // Game idx
  scores: [] as number[]
})

const game_data = {
  tiles: new Array<{ tile: Tile, id: string }>(),
  selected_tile_location: 0
}

function handle_animation(action: Animation) {
  if (action === undefined) {
    return
  }

  switch (action.tag) {
    case "newset":
      action_animator.clear_tiles_on(action.table_idx)
      const uuids = action.set.map(tile => action_animator.add_tile(tile, action.table_idx))

      if (action.table_idx === 0) {
        game_data.tiles = []
        for (let i = 0; i < uuids.length; i++) {
          game_data.tiles.push({
            id: uuids[i],
            tile: action.set[i]
          })
        }
      }
      break;
    case "add":
      const uuid = action_animator.add_tile(action.add, 0, action.location)
      if (action.table_idx === 0) {
        game_data.tiles.splice(action.location, 0, {
          id: uuid,
          tile: action.add,
        })
      }
      break;
    case "remove":
      if (action.table_idx === 0) {
        game_data.tiles.splice(action.location, 1)
      } else {
        action_animator.remove_tile(action.location, action.table_idx)
      }
      break;
    case "nakicall":
      if (action.table_idx === 0) {
        game_data.tiles
      }

      switch (action.call.action_type) {
        case ActionType.Kan:
          action_animator.naki_call(action.table_idx, 'ankan')
        case ActionType.Pon:
          action_animator.naki_call(action.table_idx, 'pon')
        case ActionType.Chii:
          action_animator.naki_call(action.table_idx, 'chii')
      }

  }
}

// Self is table idx 0
// So if we are player 3, player 0's offset is (4 + 0 - 3) % 4 = 1
function offset_to_self(other_index: GameIdx): TableIdx {
  return (4 + arena_data.value.player_idx + other_index) % 4 as TableIdx
}

const make_fsm = () => {

  const out_of_game = create_state("out_of_game", {})
  const awaiting_discard = create_state("awaiting_discard", {
    player_idx: 0 as GameIdx
  })
  const discarded = create_state("discarded", {
    discarded_by: 0 as GameIdx,
    tile_discarded: HiddenTile,
  })
  // If the player has calls that he needs to make
  const awaiting_naki_calls = create_state("awaiting_naki_calls", {})
  const naki_called = create_state("naki_called", {})
  const round_finished = create_state("round_finished", {})
  const game_finished = create_state("game_finished", {})

  const draw_event = create_event("draw_event", out_of_game, awaiting_discard, {
    callback: (_, awaiting_discard, tile_received: Tile, index: GameIdx) => {
      handle_animation({
        tag: 'add',
        add: tile_received,
        location: game_data.tiles.length,
        table_idx: offset_to_self(index),
      })

      awaiting_discard.player_idx = index
    }
  })
  const discard_event = create_event("discard_event", awaiting_discard, discarded, {
    callback: (_, discarded, action: Action, from_player: GameIdx) => {
      discarded.discarded_by = from_player
      let discard_location: number
      if (from_player === arena_data.value.player_idx) {
        discard_location = game_data.selected_tile_location
      } else {
        discard_location = -1
      }

      switch (action.action_type) {
        case ActionType.Riichi:
          discarded.tile_discarded = new Tile(action.tile_to_riichi)
          break;
        case ActionType.Toss:
          discarded.tile_discarded = new Tile(action.tile_to_toss)
      }

      handle_animation({
        tag: "remove",
        location: discard_location,
        table_idx: offset_to_self(from_player),
        riichi_called: action.action_type == ActionType.Riichi
      })
    }
  })
  // TODO
  const receive_naki_event = create_event("receive_naki_event", discarded, awaiting_naki_calls, {
    callback: (discarded, __, naki_calls: ActionType[]) => {
      console.log("Naki calls: ", naki_calls)
    }
  })
  // TODO
  const naki_called_event = create_event("naki_called_event", discarded, naki_called, {
    callback: (from, _, naki_calls: Action, from_player: GameIdx) => {

    }
  })

  const fsm = create_fsm_builder()
    .add_state(out_of_game)
    .add_state(awaiting_discard)
    .add_state(discarded)
    .add_state(awaiting_naki_calls)
    .add_state(naki_called)
    .add_state(round_finished)
    .add_state(game_finished)
    .add_event(draw_event)
    .add_event(discard_event)
    .add_event(receive_naki_event)
    .add_event(naki_called_event)
    .build(out_of_game)

  ArenaMessageBus.register((data: ServerEvent) => {
    if (data.arena_message.arenaevent_type != ArenaEventType.ArenaBoardEvent) {
      return true
    }

    const board_event = data.arena_message.board_event
    switch (board_event.boardevent_type) {
      case BoardEventType.PotentialActionEvent:
        handle_potential_action_event(board_event.actions)
        break
      case BoardEventType.PlayerActionEvent:
        handle_player_action_event(board_event.action_data, board_event.from_player)
        break
      case BoardEventType.GameSetupEvent:
        handle_game_setup_event(board_event.setup)
        break
      case BoardEventType.GameEndEvent:
		throw new Error("Not yet implemented: GameEndEvent")
        break;
    }

    return true
  })

  return fsm
}


const three_canvas = ref<HTMLCanvasElement>()
const game_container = ref<HTMLDivElement>()
const is_fullscreen = ref(false)

const discard_required = ref(false)

let animation_id: number

let renderer: IRenderer
let selection_manager: ISelectionManager
let action_animator: IActionAnimator

let fsm: ReturnType<typeof make_fsm>

onMounted(() => {
  if (!three_canvas.value) return
  initialize_tiles()

  const manager = new Scene(three_canvas.value)
  renderer = manager
  selection_manager = manager
  action_animator = manager
  animate(0, 0)

  window.addEventListener('resize', on_window_resize)
  window.addEventListener('click', on_click)

  fsm = make_fsm()

  ArenaMessageBus.register(message_handler)
  ArenaMessageBus.register(debug_message_printer)
})

function on_click(event: MouseEvent) {
  const selection = selection_manager.get_selection()
  if (selection === null) {
    return
  }

  if (!discard_required.value) {
    console.log("Discard not required")
    return
  }

  game_data.selected_tile_location = selection.location
  discard_required.value = false

  let msg_idx = websocket_state.conn.send({
    message_type: MessageType.REQUEST,
    data: {
      serveraction_type: ServerActionType.ServerArenaAction,
      arena_action: {
        arenaaction_type: ArenaActionType.PlayerActionData,
        action: {
          action_type: ActionType.Toss,
          tile_to_toss: game_data.tiles[selection.location].tile.value,
        }
      }
    }
  })

  let ret = register_request(msg_idx)
  ret.then(() => {
    console.log("Discard action acknowledged by server")
  }).catch(() => {
    console.error("Discard action was not acknowledged by server")
  })
}

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
  animation_id = requestAnimationFrame(new_t => { animate(new_t, new_t - t) })
}

onUnmounted(() => {
  if (animation_id) {
    cancelAnimationFrame(animation_id)
  }
  window.removeEventListener('resize', on_window_resize)
  window.removeEventListener('click', on_click)

  renderer.stop()
})

function debug_message_printer(event: ServerEvent) {
  const msg = `ServerEvent message: ${ServerEventType[event.serverevent_type]}\n\t`

  const debug_board_event = function (board_event: BoardEvent, msg: string) {
    msg += `BoardEvent message: ${BoardEventType[board_event.boardevent_type]}\n\t\t`
    switch (board_event.boardevent_type) {
      case BoardEventType.PlayerActionEvent:
        msg += `PlayerActionEvent message: ${ActionType[board_event.action_data.action_type]} from player ${board_event.from_player}`
        break
      case BoardEventType.PotentialActionEvent:
        msg += `PotentialActionEvent message: ${board_event.actions.map((action) => ActionType[action.action_type])}`
        break
      case BoardEventType.GameSetupEvent:
        msg += `GameSetupEvent message: ${board_event.setup.map((setup) => SetupType[setup.setup_type])}`
        break
      case BoardEventType.GameEndEvent:
        msg += `GameEndEvent message: ${board_event.result}`

    }
  }

  switch (event.arena_message.arenaevent_type) {
    case ArenaEventType.ArenaBoardEvent:
      debug_board_event(event.arena_message.board_event, msg)
  }

  return true
}

function message_handler(event: ServerEvent) {
  switch (event.arena_message.arenaevent_type) {
    case ArenaEventType.ArenaBoardEvent:
      handle_board_event(event.arena_message.board_event)
  }
  return true
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
      break
    default:
      throw new Error("Unexpected")
  }
}

function handle_player_action_event(action: Action, from_player: number) {
  switch (action.action_type) {
    case ActionType.Tsumo:
      // TODO
      throw new Error("Win")
    case ActionType.Ron:
      // TODO
      throw new Error("Win")
    case ActionType.Riichi:
    case ActionType.Toss:
      fsm.trigger_event("discard_event", action, from_player as GameIdx)
      break;
    case ActionType.Pon:
    case ActionType.Kan:
    case ActionType.Chii:
      fsm.trigger_event("naki_called_event", action, from_player as GameIdx)
      break;
    case ActionType.Draw:
      fsm.trigger_event("draw_event", new Tile(action.drawn_tile), from_player as GameIdx)
      break;
    default:
      console.error("Unexpected action performed")
  }
}

function handle_potential_action_event(actions: Action[]) {
  for (let action of actions) {
    switch (action.action_type) {
      case ActionType.Tsumo:
        console.log("Tsumo possible")
        break;
      case ActionType.Ron:
        console.log("Ron possible")
        break;
      case ActionType.Riichi:
        console.log("Riichi possible")
        break;
      case ActionType.Toss:
        discard_required.value = true
        break;
      case ActionType.Skip:
        // need to show the player a prompt to skip
        console.log("Skip possible")
        break;
      case ActionType.Pon:
        console.log("Pon possible")
        break;
      case ActionType.Kan:
        console.log("Kan possible")
        break;
      case ActionType.Chii:
        console.log("Chii possible")
        break;
      case ActionType.Draw:
        throw new Error("Unexpected draw")
        break;
    }
  }
}

function handle_game_setup_event(setups: Setup[]) {
  for (let setup of setups) {
    switch (setup.setup_type) {
      case SetupType.INITIAL_TILES:
        handle_animation({
          tag: "newset",
          set: Tile.from(setup.data).sort(Tile.sort),
          table_idx: 0 as TableIdx,
        })
        break
      case SetupType.DORA:
        action_animator.add_dora(new Tile(setup.data))
        break
      case SetupType.STARTING_POINTS:
        arena_data.value.scores = setup.data
        break
      case SetupType.PLAYER_NUMBER:
		console.log("Player number setup received: ", setup.data)
        arena_data.value.player_idx = setup.data as GameIdx
        break
      case SetupType.PLAYER_ORDER:
        throw new Error("Should be removed")
      case SetupType.ROUND_WIND:
        arena_data.value.round_wind = setup.data
        break
      case SetupType.ROUND_NUMBER:
        arena_data.value.round_number = setup.data
        break
    }
  }
}
</script>

<template>
  <ScoreBoard v-if="in_game" :scoreboard_values="arena_data.scores" :round_wind="arena_data.round_wind"
    :round_number="arena_data.round_number" :players="arena_data.players" :player_idx="arena_data.player_idx">
  </ScoreBoard>
  <div ref="game_container" class="game-view-container" :class="{ fullscreen: is_fullscreen }">
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
