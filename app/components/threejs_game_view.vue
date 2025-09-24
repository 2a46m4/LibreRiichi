<script setup lang="ts">
import * as THREE from 'three'
import {OrbitControls} from 'three/examples/jsm/controls/OrbitControls.js'
import {ref, onMounted, onUnmounted} from 'vue'
import {Tile, TileValue} from "../game/tile";
import {alphatest_colour, load_texture} from "../render/texture";
import {MeshLambertMaterial} from "three";
import {load_all_textures} from "../render/tile";

const props = defineProps<{
  tiles: Tile[]
}>()

const threeCanvas = ref<HTMLCanvasElement>()
const gameContainer = ref<HTMLDivElement>()
const isFullscreen = ref(false)

let scene: THREE.Scene
let camera: THREE.PerspectiveCamera
let renderer: THREE.WebGLRenderer
let controls: OrbitControls
let animation_id: number

const raycaster = new THREE.Raycaster();
const pointer = new THREE.Vector2();

const mahjong_tiles: Map<string, THREE.Mesh> = new Map<string, THREE.Mesh>()
let selected: THREE.Mesh | null = null
let selected_old_mat: THREE.Material
let selected_mat = new THREE.MeshLambertMaterial({color: 0xff0000})

const dealerMarker: THREE.Mesh[] = []

let textures: Map<number, THREE.Texture>

onMounted(() => {
  if (!threeCanvas.value) return

  setupScene()
  textures = load_all_textures()
  createMahjongTable()
  createMahjongTiles()
  createDealerMarker()
  animate()

  window.addEventListener('resize', onWindowResize)


  function onPointerMove(event: PointerEvent) {

    // calculate pointer position in normalized device coordinates
    // (-1 to +1) for both components

    pointer.x = (event.clientX / window.innerWidth) * 2 - 1;
    pointer.y = -(event.clientY / window.innerHeight) * 2 + 1;

  }

  window.addEventListener('pointermove', onPointerMove);
})

function setupScene() {
  // Scene setup
  scene = new THREE.Scene()
  scene.background = new THREE.Color(0xffffff) // Mahjong table green

  // Camera setup
  camera = new THREE.PerspectiveCamera(
      35,
      threeCanvas.value!.clientWidth / threeCanvas.value!.clientHeight,
      0.1,
      1000,
  )
  camera.position.set(0, 4, 20)
  camera.lookAt(0, 0, 0)

  // Renderer setup
  renderer = new THREE.WebGLRenderer({
    canvas: threeCanvas.value!,
    antialias: true,
  })
  renderer.setSize(
      threeCanvas.value!.clientWidth,
      threeCanvas.value!.clientHeight,
  )
  renderer.shadowMap.enabled = true
  renderer.shadowMap.type = THREE.PCFShadowMap

  // Controls setup
  controls = new OrbitControls(camera, renderer.domElement)
  controls.enableDamping = true
  controls.dampingFactor = 0.05
  controls.maxPolarAngle = Math.PI / 2
  controls.minPolarAngle = Math.PI / 5
  controls.maxAzimuthAngle = Math.PI / 12
  controls.minAzimuthAngle = -Math.PI / 12
  controls.maxDistance = 20
  controls.minDistance = 15
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

function createMahjongTable() {
  // Table base
  const tableGeometry = new THREE.BoxGeometry(13, 13, 0.5)
  const tableMaterial = new THREE.MeshLambertMaterial({color: 0x8b4513}) // Brown wood
  const table = new THREE.Mesh(tableGeometry, tableMaterial)
  table.position.y = -2
  table.rotation.x = Math.PI / 2
  table.receiveShadow = true
  scene.add(table)

  // Table surface (green felt)
  const surfaceGeometry = new THREE.BoxGeometry(12.8, 12.8, 0.1)
  const surfaceMaterial = new THREE.MeshLambertMaterial({color: 0x0a7c4a})
  const surface = new THREE.Mesh(surfaceGeometry, surfaceMaterial)
  surface.position.y = -1.7
  surface.rotation.x = Math.PI / 2
  surface.receiveShadow = true
  scene.add(surface)

  // Player positions (4 sides of a square)
  const positions = [
    {x: 0, z: 4, rotation: 0},
    {x: 4, z: 0, rotation: Math.PI / 2}, // East
    {x: 0, z: -4, rotation: Math.PI}, // North
    {x: -4, z: 0, rotation: -Math.PI / 2}, // West
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

function createMahjongTiles() {
  // Create sample mahjong tiles around the table
  const tileGeometry = new THREE.BoxGeometry(0.4, 0.6, 0.25)

  const player_position = {x: 0, z: 5, rotation: 0}

  // Create tile walls for each player
  const positions = [
    {x: 5, z: 0, rotation: Math.PI / 2}, // East
    {x: 0, z: -5, rotation: Math.PI}, // North
    {x: -5, z: 0, rotation: -Math.PI / 2}, // West
  ]

  for (const [i, tile] of props.tiles.entries()) {
    console.log(tile, tile.value)
    const texture = textures.get(tile.value)!
    texture.colorSpace = THREE.SRGBColorSpace
    const material = new THREE.MeshLambertMaterial({
      map: texture,
      color: 0xffffff,
      transparent: false,
      side: THREE.DoubleSide,
      alphaTest: 0.9
    })
    material.onBeforeCompile = alphatest_colour

    const tile_mesh = new THREE.Mesh(tileGeometry, material)

    // Position tiles in a row
    const offsetX = (i - 6) * 0.45
    tile_mesh.position.set(
        player_position.x + Math.cos(player_position.rotation) * offsetX,
        -1.3,
        player_position.z + Math.sin(player_position.rotation) * offsetX,
    )
    tile_mesh.rotation.y = player_position.rotation
    tile_mesh.castShadow = true
    tile_mesh.receiveShadow = true

    mahjong_tiles.set(tile_mesh.uuid, tile_mesh)
    scene.add(tile_mesh)
  }

  let blank_tile = new Tile(TileValue.Hidden)
  const blank_texture = textures.get(blank_tile.value)!
  blank_texture.colorSpace = THREE.SRGBColorSpace
  const blank_material = new THREE.MeshLambertMaterial({
    map: blank_texture,
    color: 0xffffff,
    transparent: false,
    side: THREE.DoubleSide,
    alphaTest: 0.9
  })
  blank_material.onBeforeCompile = alphatest_colour

  positions.forEach((pos) => {
    for (let i = 0; i < 13; i++) {
      const tile = new THREE.Mesh(
          tileGeometry, blank_material
      )

      // Position tiles in a row
      const offsetX = (i - 6) * 0.45
      tile.position.set(
          pos.x + Math.cos(pos.rotation) * offsetX,
          -1.3,
          pos.z + Math.sin(pos.rotation) * offsetX,
      )
      tile.rotation.y = pos.rotation
      tile.castShadow = true
      tile.receiveShadow = true

      scene.add(tile)
    }
  })

  // Center tiles (discarded pile)
  for (let i = 0; i < 12; i++) {
    const tile = new THREE.Mesh(tileGeometry, new MeshLambertMaterial())
    const angle = (i / 12) * Math.PI * 2
    const radius = 1.5
    tile.position.set(Math.cos(angle) * radius, -1.3, Math.sin(angle) * radius)
    tile.rotation.y = angle + Math.PI / 2
    tile.castShadow = true
    tile.receiveShadow = true
    scene.add(tile)
  }
}

function createDealerMarker() {
  // Dealer button/marker
  const markerGeometry = new THREE.CylinderGeometry(0.3, 0.3, 0.1, 8)
  const markerMaterial = new THREE.MeshLambertMaterial({color: 0xff6b6b})
  const marker = new THREE.Mesh(markerGeometry, markerMaterial)
  marker.position.set(2, -1.2, 2)
  marker.castShadow = true
  dealerMarker.push(marker)
  scene.add(marker)

  // Dealer marker text indicator
  const textGeometry = new THREE.RingGeometry(0.1, 0.2, 6)
  const textMaterial = new THREE.MeshLambertMaterial({color: 0xffffff})
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

  // update the picking ray with the camera and pointer position
  raycaster.setFromCamera(pointer, camera);

  // calculate objects intersecting the picking ray
  const intersects = raycaster.intersectObjects(scene.children);
  console.log(intersects.length)

  let intersect_occurred = false
  for (let i = 0; i < intersects.length; i++) {
    if (mahjong_tiles.has(intersects[i].object.uuid)) {
      intersect_occurred = true

      if (selected !== null) {
        selected.material = selected_old_mat
        selected = null
      }

      let tile = mahjong_tiles.get(intersects[i].object.uuid)
      if (tile === undefined) {
        throw new Error("Tile not found")
      }

      selected = tile
      selected_old_mat = tile.material as THREE.Material
      tile.material = selected_mat
      break
    }
  }

  if (!intersect_occurred && selected !== null) {
    selected.material = selected_old_mat
    selected = null
  }

  controls.update()
  renderer.render(scene, camera)
}

function onWindowResize() {
  if (!threeCanvas.value || !gameContainer.value) return

  const width = gameContainer.value.clientWidth
  const height = gameContainer.value.clientHeight

  threeCanvas.value.width = width
  threeCanvas.value.height = height

  camera.aspect = width / height
  camera.updateProjectionMatrix()
  renderer.setSize(width, height)
}

onUnmounted(() => {
  if (animation_id) {
    cancelAnimationFrame(animation_id)
  }
  window.removeEventListener('resize', onWindowResize)

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
</script>

<template>
  <div ref="gameContainer" class="game-view-container" :class="{ fullscreen: isFullscreen }">
    <canvas ref="threeCanvas" class="game-canvas"></canvas>
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
