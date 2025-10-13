<script setup lang="ts">
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import {ref, onMounted, onUnmounted, watch, Ref} from 'vue'
import {decode, Tile, TileValue} from '../game/tile'
import { initialize_tiles, TileObject } from '../render/tile'
import { BoardEvent, BoardEventType } from '../messaging/board_event_generated'
import { Setup, SetupType } from '../game/setup'
import { ArenaMessageBus } from "../messaging/event_handler";
import { ArenaEventType } from "../messaging/arena_event_generated";
import { ServerEvent } from "../messaging/server_event_generated";
import {ScoreboardState} from "../game/scoreboard";
import ScoreBoard from "../components/scoreboard.vue"
import {Action, ActionType} from "../messaging/action_generated";

const props = defineProps<{in_game: boolean}>()

const three_canvas = ref<HTMLCanvasElement>()
const gameContainer = ref<HTMLDivElement>()
const isFullscreen = ref(false)

let scene: THREE.Scene
let camera: THREE.PerspectiveCamera
let renderer: THREE.WebGLRenderer
let controls: OrbitControls
let animation_id: number

const raycaster = new THREE.Raycaster()
const pointer = new THREE.Vector2()

// mesh uuid to mesh, represents hand
const mahjong_tiles: Map<string, THREE.Mesh> = new Map<string, THREE.Mesh>()
const player_position = { x: 0, z: 5, rotation: 0 }
const dora_tile_position = { x: -5, z: 5, y: -1.3 }

const scoreboard_state: Ref<ScoreboardState> = ref({
  scoreboard_values: [],
  player_to_order_map: [],
  round_wind: 0,
  round_number: 0,
  players: [],
  player_idx: 0,
})

let selection: {
  material: THREE.MeshLambertMaterial
  mesh: THREE.Mesh
  tile: TileObject | null
}
function init_selection() {
  let material = new THREE.MeshLambertMaterial({
    color: 0xffff00,
    transparent: true,
    opacity: 0.5,
  })
  let mesh = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.7, 0.26), material)
  mesh.visible = false
  let selected_tile: TileObject | null = null

  selection = {
    material: material,
    mesh: mesh,
    tile: selected_tile,
  }

  scene.add(selection.mesh)
}

const dealerMarker: THREE.Mesh[] = []

let dora_tiles = []

onMounted(() => {
  if (!three_canvas.value) return

  setup_scene()
  initialize_tiles()
  init_selection()
  create_mahjong_table()
  create_mahjong_tiles()
  create_dealer_marker()
  animate()

  window.addEventListener('resize', on_window_resize)
  window.addEventListener('pointermove', on_pointer_move)
  window.addEventListener('click', on_click)

  ArenaMessageBus.register(message_handler)
})

function on_click(event: MouseEvent) {
  if (selection.tile !== null) {
    console.log('clicked on', selection.tile)
  }
}

// ndc coords, update pointer
function on_pointer_move(event: PointerEvent) {
  pointer.x = (event.clientX / window.innerWidth) * 2 - 1
  pointer.y = -(event.clientY / window.innerHeight) * 2 + 1
}

function on_window_resize() {
  if (!three_canvas.value || !gameContainer.value) return

  const width = gameContainer.value.clientWidth
  const height = gameContainer.value.clientHeight

  three_canvas.value.width = width
  three_canvas.value.height = height

  camera.aspect = width / height
  camera.updateProjectionMatrix()
  renderer.setSize(width, height)
}

function setup_scene() {
  // Scene setup
  scene = new THREE.Scene()
  scene.background = new THREE.Color(0xffffff) // Mahjong table green

  // Camera setup
  camera = new THREE.PerspectiveCamera(
    35,
    three_canvas.value!.clientWidth / three_canvas.value!.clientHeight,
    0.1,
    1000,
  )
  camera.position.set(0, 4, 20)
  camera.lookAt(0, 0, 0)

  // Renderer setup
  renderer = new THREE.WebGLRenderer({
    canvas: three_canvas.value!,
    antialias: true,
  })
  renderer.setSize(
    three_canvas.value!.clientWidth,
    three_canvas.value!.clientHeight,
  )
  renderer.shadowMap.enabled = true
  renderer.shadowMap.type = THREE.PCFShadowMap

  // Controls setup
  controls = new OrbitControls(camera, renderer.domElement)
  controls.enableDamping = true
  controls.dampingFactor = 0.05
  controls.maxPolarAngle = Math.PI / 2
  controls.minPolarAngle = 0
  controls.maxAzimuthAngle = Math.PI / 12
  controls.minAzimuthAngle = -Math.PI / 12
  controls.maxDistance = 20
  controls.minDistance = 12
  controls.enablePan = false

  // Add lighting
  const ambientLight = new THREE.AmbientLight(0x404040, 0.6)
  scene.add(ambientLight)

  const directionalLight = new THREE.DirectionalLight(0xffffff, 1)
  directionalLight.position.set(10, 15, 5)
  directionalLight.castShadow = true
  directionalLight.shadow.mapSize.width = 512
  directionalLight.shadow.mapSize.height = 512
  directionalLight.shadow.camera.near = 15
  directionalLight.shadow.camera.far = 25
  directionalLight.shadow.camera.left = -5
  directionalLight.shadow.camera.right = 6
  directionalLight.shadow.camera.top = 4.2
  directionalLight.shadow.camera.bottom = -6
  directionalLight.shadow.bias = 0.0003
  scene.add(directionalLight)
}

function create_mahjong_table() {
  // Table base
  const tableGeometry = new THREE.BoxGeometry(13, 13, 0.5)
  const tableMaterial = new THREE.MeshLambertMaterial({ color: 0x8b4513 }) // Brown wood
  const table = new THREE.Mesh(tableGeometry, tableMaterial)
  table.position.y = -2
  table.rotation.x = Math.PI / 2
  table.receiveShadow = true
  scene.add(table)

  // Table surface (green felt)
  const surfaceGeometry = new THREE.BoxGeometry(12.8, 12.8, 0.1)
  const surfaceMaterial = new THREE.MeshLambertMaterial({ color: 0x0a7c4a })
  const surface = new THREE.Mesh(surfaceGeometry, surfaceMaterial)
  surface.position.y = -1.7
  surface.rotation.x = Math.PI / 2
  surface.receiveShadow = true
  scene.add(surface)

  // Player positions (4 sides of a square)
  const positions = [
    { x: 0, z: 4, rotation: 0 },
    { x: 4, z: 0, rotation: Math.PI / 2 }, // East
    { x: 0, z: -4, rotation: Math.PI }, // North
    { x: -4, z: 0, rotation: -Math.PI / 2 }, // West
  ]

  positions.forEach((pos, index) => {
    // Player area markers
    const markerGeometry = new THREE.PlaneGeometry(3, 0.5)
    const markerMaterial = new THREE.MeshLambertMaterial({
      color: index === 0 ? 0xff6b6b : 0x4a90e2,
      transparent: true,
      opacity: 0.7,
    })
    const marker = new THREE.Mesh(markerGeometry, markerMaterial)
    marker.position.set(pos.x, -1.6, pos.z)
    marker.rotation.x = -Math.PI / 2
    marker.rotation.z = pos.rotation
    scene.add(marker)
  })
}

function create_mahjong_tiles() {

  const tiles = [0, 1, 2, 3, 4, 5, 6, 7, 8, 16, 17, 18, 19].map(
    (i) => new Tile(i),
  )
  for (const [i, tile] of tiles.entries()) {
    let tile_obj = new TileObject(tile)

    const offsetX = (i - 6) * 0.45
    tile_obj.position.set(
      player_position.x + Math.cos(player_position.rotation) * offsetX,
      -1.3,
      player_position.z + Math.sin(player_position.rotation) * offsetX,
    )
    tile_obj.rotation.y = player_position.rotation
    tile_obj.castShadow = true
    tile_obj.receiveShadow = true

    mahjong_tiles.set(tile_obj.uuid, tile_obj)
    scene.add(tile_obj)
  }

  // Create tile walls for each player
  const positions = [
    { x: 5, z: 0, rotation: Math.PI / 2 }, // East
    { x: 0, z: -5, rotation: Math.PI }, // North
    { x: -5, z: 0, rotation: -Math.PI / 2 }, // West
  ]
  let blank_tile = new Tile(TileValue.Hidden)
  positions.forEach((pos) => {
    for (let i = 0; i < 13; i++) {
      let blank_tile_obj = new TileObject(blank_tile)

      // Position tiles in a row
      const offsetX = (i - 6) * 0.45
      blank_tile_obj.position.set(
        pos.x + Math.cos(pos.rotation) * offsetX,
        -1.3,
        pos.z + Math.sin(pos.rotation) * offsetX,
      )
      blank_tile_obj.rotation.y = pos.rotation
      blank_tile_obj.castShadow = true
      blank_tile_obj.receiveShadow = true

      scene.add(blank_tile_obj)
    }
  })

  // Center tiles (discarded pile)
  for (let i = 0; i < 12; i++) {
    const tile_obj = new TileObject(new Tile(TileValue.Invalid))
    const angle = (i / 12) * Math.PI * 2
    const radius = 1.5
    tile_obj.position.set(
      Math.cos(angle) * radius,
      -1.3,
      Math.sin(angle) * radius,
    )
    tile_obj.rotation.y = angle + Math.PI / 2
    tile_obj.castShadow = true
    tile_obj.receiveShadow = true
    scene.add(tile_obj)
  }
}

function create_dealer_marker() {
  // Dealer button/marker
  const markerGeometry = new THREE.CylinderGeometry(0.3, 0.3, 0.1, 8)
  const markerMaterial = new THREE.MeshLambertMaterial({ color: 0xff6b6b })
  const marker = new THREE.Mesh(markerGeometry, markerMaterial)
  marker.position.set(2, -1.2, 2)
  marker.castShadow = true
  dealerMarker.push(marker)
  scene.add(marker)

  // Dealer marker text indicator
  const textGeometry = new THREE.RingGeometry(0.1, 0.2, 6)
  const textMaterial = new THREE.MeshLambertMaterial({ color: 0xffffff })
  const textRing = new THREE.Mesh(textGeometry, textMaterial)
  textRing.position.set(2, -1.1, 2)
  textRing.rotation.x = -Math.PI / 2
  dealerMarker.push(textRing)
  scene.add(textRing)
}

function animate() {
  animation_id = requestAnimationFrame(animate)

  // Gentle rotation of dealer marker
  dealerMarker.forEach((marker) => {
    marker.rotation.y += 0.005
  })

  raycaster.setFromCamera(pointer, camera)
  const intersects = raycaster.intersectObjects(scene.children)

  let intersect_occurred = false
  for (let i = 0; i < intersects.length; i++) {
    if (mahjong_tiles.has(intersects[i].object.uuid)) {
      intersect_occurred = true
      let tile = mahjong_tiles.get(intersects[i].object.uuid)
      if (tile === undefined) {
        throw new Error('Tile not found')
      }

      selection.mesh.visible = true
      selection.mesh.position.copy(tile.position)
      selection.tile = tile as TileObject

      break
    }
  }

  if (!intersect_occurred) {
    selection.mesh.visible = false
  }

  controls.update()
  renderer.render(scene, camera)
}

onUnmounted(() => {
  if (animation_id) {
    cancelAnimationFrame(animation_id)
  }
  window.removeEventListener('resize', on_window_resize)
  window.removeEventListener('pointermove', on_pointer_move)
  window.removeEventListener('click', on_click)

  // Cleanup Three.js resources
  if (renderer) {
    renderer.dispose()
  }
  if (scene) {
    scene.traverse((object) => {
      if (object.type === 'Mesh') {
        const mesh = object as THREE.Mesh
        mesh.geometry.dispose()
        if (Array.isArray(mesh.material)) {
          mesh.material.forEach((material) => material.dispose())
        } else {
          mesh.material.dispose()
        }
      }
    })
  }
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
  <div ref="gameContainer" class="game-view-container" :class="{ fullscreen: isFullscreen }">
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
