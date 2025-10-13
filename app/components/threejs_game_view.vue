<script setup lang="ts">
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import {ref, onMounted, onUnmounted, watch, Ref, render} from 'vue'
import {decode, Tile} from '../game/tile'
import { initialize_tiles, TileObject } from '../render/tile'
import { BoardEvent, BoardEventType } from '../messaging/board_event_generated'
import { Setup, SetupType } from '../game/setup'
import { ArenaMessageBus } from "../messaging/event_handler";
import { ArenaEventType } from "../messaging/arena_event_generated";
import { ServerEvent } from "../messaging/server_event_generated";
import {ScoreboardState} from "../game/scoreboard";
import ScoreBoard from "../components/scoreboard.vue"
import {Action, ActionType} from "../messaging/action_generated";
import {IRenderer, ThreeJSRenderer} from "../render/renderer_setup";
import {Raycaster} from "../render/raycaster";

const props = defineProps<{in_game: boolean}>()

const three_canvas = ref<HTMLCanvasElement>()
const game_container = ref<HTMLDivElement>()
const is_fullscreen = ref(false)

let scene: THREE.Scene
let camera: THREE.PerspectiveCamera
let controls: OrbitControls
let animation_id: number

// mesh uuid to mesh, represents hand
const mahjong_tiles: Map<string, THREE.Mesh> = new Map<string, THREE.Mesh>()
const player_position = { x: 0, z: 5, rotation: 0 }
const dora_tile_position = { x: -5, z: 5, y: -1.3 }

let renderer: IRenderer
let raycaster: Raycaster

const scoreboard_state: Ref<ScoreboardState> = ref({
  scoreboard_values: [],
  player_to_order_map: [],
  round_wind: 0,
  round_number: 0,
  players: [],
  player_idx: 0,
})

let dora_tiles = []

onMounted(() => {
  if (!three_canvas.value) return
  initialize_tiles()
  renderer = new ThreeJSRenderer(three_canvas.value)
  raycaster = new Raycaster(renderer)
  animate()

  window.addEventListener('resize', on_window_resize)
  window.addEventListener('click', on_click)

  ArenaMessageBus.register(message_handler)
})

function on_click(event: MouseEvent) {
    console.log('clicked')
}

function on_window_resize() {
  if (!three_canvas.value || !game_container.value) return

  const width = game_container.value.clientWidth
  const height = game_container.value.clientHeight

  three_canvas.value.width = width
  three_canvas.value.height = height

  renderer.window_resize(width, height)
}

function animate() {
  animation_id = requestAnimationFrame(animate)

  let selections = raycaster.get_selections()
  renderer.render_selection(selections)
  renderer.animate_frame()
}

onUnmounted(() => {
  if (animation_id) {
    cancelAnimationFrame(animation_id)
  }
  window.removeEventListener('resize', on_window_resize)
  window.removeEventListener('click', on_click)

  renderer.stop()
  raycaster.stop()
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
      // Handle unknown event type
      break
  }
}

function handle_player_action_event(action: PlayerActionEvent) {

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
        mahjong_tiles.forEach((mesh) => scene.remove(mesh))
        let tiles = Tile.from(setup.data)
        tiles.sort(Tile.sort)
        let tile_objs = tiles.map((tile) => new TileObject(tile))
        tile_objs.forEach((tile_obj, i) => {
          const offsetX = (i - 6) * 0.45
          tile_obj.position.set(
            player_position.x + Math.cos(player_position.rotation) * offsetX,
            -1.3,
            player_position.z + Math.sin(player_position.rotation) * offsetX,
          )
          tile_obj.rotation.y = player_position.rotation
          tile_obj.castShadow = true
          tile_obj.receiveShadow = true
          scene.add(tile_obj)
          mahjong_tiles.set(tile_obj.uuid, tile_obj)
        })
        break

      case SetupType.DORA:
        const dora_tile = new Tile(setup.data)
        dora_tiles.push(dora_tile)
        const dora_tile_obj = new TileObject(dora_tile)
        dora_tile_obj.position.set(
          dora_tile_position.x, dora_tile_position.y, dora_tile_position.z
        )
        scene.add(dora_tile_obj)

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
  <ScoreBoard v-if="in_game"
    :scoreboard_values="scoreboard_state.scoreboard_values"
    :player_to_order_map="scoreboard_state.player_to_order_map"
    :round_wind="scoreboard_state.round_wind"
    :round_number="scoreboard_state.round_number"
    :players="scoreboard_state.players"
    :player_idx="scoreboard_state.player_idx"
  >
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
