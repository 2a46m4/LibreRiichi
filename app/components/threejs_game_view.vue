<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, Ref, render } from 'vue'
import { decode, Tile } from '../game/tile'
import { initialize_tiles } from '../render/tile'
import { BoardEvent, BoardEventType, PlayerActionEvent } from '../messaging/board_event_generated'
import { Setup, SetupType } from '../game/setup'
import { ArenaMessageBus } from "../messaging/event_handler";
import { ArenaEventType } from "../messaging/arena_event_generated";
import { ServerEvent } from "../messaging/server_event_generated";
import { ScoreboardState } from "../game/scoreboard";
import ScoreBoard from "../components/scoreboard.vue"
import { Action, ActionType } from "../messaging/action_generated";
import { IActionAnimator, IRenderer, ISelectionManager, ThreeJSRenderer } from "../render/renderer";
import { IArena } from '../game/arena'

const props = defineProps<{ in_game: boolean, arena: IArena }>()

const three_canvas = ref<HTMLCanvasElement>()
const game_container = ref<HTMLDivElement>()
const is_fullscreen = ref(false)

let animation_id: number

let renderer: IRenderer
let selection_manager: ISelectionManager
let action_animator: IActionAnimator

const scoreboard_state: Ref<ScoreboardState> = ref({
  scoreboard_values: [],
  player_to_order_map: [],
  round_wind: 0,
  round_number: 0,
  players: [],
  player_idx: 0,
})

onMounted(() => {
  if (!three_canvas.value) return
  initialize_tiles()

  const manager = new ThreeJSRenderer(three_canvas.value)
  renderer = manager
  selection_manager = manager
  action_animator = manager
  animate(0)

  window.addEventListener('resize', on_window_resize)
  window.addEventListener('click', on_click)

  ArenaMessageBus.register(message_handler)
})

function on_click(event: MouseEvent) {
  const selection = selection_manager.get_selection()
  if (selection === null) {
    return
  } else {
    console.warn("Not yet implemented")
  }
}

function on_window_resize() {
  if (!three_canvas.value || !game_container.value) return

  const width = game_container.value.clientWidth
  const height = game_container.value.clientHeight

  three_canvas.value.width = width
  three_canvas.value.height = height

  renderer.window_resize(width, height)
}

function animate(t: number) {
  animation_id = requestAnimationFrame(animate)
  renderer.animate_frame(t)
}

onUnmounted(() => {
  if (animation_id) {
    cancelAnimationFrame(animation_id)
  }
  window.removeEventListener('resize', on_window_resize)
  window.removeEventListener('click', on_click)

  renderer.stop()
})

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
      handle_player_action_event(new_event)
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

function handle_player_action_event(action: PlayerActionEvent) {
  switch (action.action_data.action_type) {
    case ActionType.Tsumo:
      throw new Error("Win")
    case ActionType.Ron:
      throw new Error("Win")
    case ActionType.Riichi:
      break;
    case ActionType.Toss:
      break;
    case ActionType.Skip:
      break;
    case ActionType.Pon:
      break;
    case ActionType.Kan:
      break;
    case ActionType.Chii:
      break;
    case ActionType.Draw:
      action_animator.draw(action.action_data.drawn_tile)
      break;
  }
}

function handle_potential_action_event(actions: Action[]) {
  for (let action of actions) {
    switch (action.action_type) {
      case ActionType.Tsumo:
        break;
      case ActionType.Ron:
        break;
      case ActionType.Riichi:
        break;
      case ActionType.Toss:
        break;
      case ActionType.Skip:
        break;
      case ActionType.Pon:
        break;
      case ActionType.Kan:
        break;
      case ActionType.Chii:
        break;
      case ActionType.Draw:
        break;
    }
  }
}

function handle_game_setup_event(setups: Setup[]) {
  for (let setup of setups) {
    switch (setup.setup_type) {
      case SetupType.INITIAL_TILES:
        action_animator.clear_tiles()
        Tile.from(setup.data)
          .sort(Tile.sort)
          .forEach((tile) => action_animator.add_tile(tile))
        break
      case SetupType.DORA:
        action_animator.add_dora(new Tile(setup.data))
        break
      case SetupType.STARTING_POINTS:
        scoreboard_state.value.scoreboard_values = setup.data
        break
      case SetupType.PLAYER_NUMBER:
        scoreboard_state.value.player_idx = setup.data
        break
      case SetupType.PLAYER_ORDER:
        scoreboard_state.value.player_to_order_map = Array.from(decode(setup.data))
        break
      case SetupType.ROUND_WIND:
        scoreboard_state.value.round_wind = setup.data
        break
      case SetupType.ROUND_NUMBER:
        scoreboard_state.value.round_number = setup.data
        break
    }
  }
}
</script>

<template>
  <ScoreBoard v-if="in_game" :scoreboard_values="scoreboard_state.scoreboard_values"
    :player_to_order_map="scoreboard_state.player_to_order_map" :round_wind="scoreboard_state.round_wind"
    :round_number="scoreboard_state.round_number" :players="scoreboard_state.players"
    :player_idx="scoreboard_state.player_idx">
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
