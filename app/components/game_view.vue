<script setup lang="ts">
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps<{}>()

const threeCanvas = ref<HTMLCanvasElement>()
const gameContainer = ref<HTMLDivElement>()
const isFullscreen = ref(false)

let scene: THREE.Scene
let camera: THREE.PerspectiveCamera
let renderer: THREE.WebGLRenderer
let controls: OrbitControls
let animationId: number

const mahjongTiles: THREE.Mesh[] = []
const dealerMarker: THREE.Mesh[] = []


onMounted(() => {
  if (!threeCanvas.value) return

  setupScene()
  createMahjongTable()
  createMahjongTiles()
  createDealerMarker()
  animate()

  window.addEventListener('resize', onWindowResize)
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
  renderer.shadowMap.type = THREE.PCFSoftShadowMap

  // Controls setup
  controls = new OrbitControls(camera, renderer.domElement)
  controls.enableDamping = true
  controls.dampingFactor = 0.05
  controls.maxPolarAngle = Math.PI / 2
  controls.minPolarAngle = Math.PI / 5
  controls.maxAzimuthAngle = Math.PI/12
  controls.minAzimuthAngle = -Math.PI/12
  controls.maxDistance = 20
  controls.minDistance = 15

  // Add lighting
  const ambientLight = new THREE.AmbientLight(0x404040, 0.6)
  scene.add(ambientLight)

  const directionalLight = new THREE.DirectionalLight(0xffffff, 1)
  directionalLight.position.set(10, 15, 5)
  directionalLight.castShadow = true
  directionalLight.shadow.mapSize.width = 2048
  directionalLight.shadow.mapSize.height = 2048
  directionalLight.shadow.camera.near = 0.1
  directionalLight.shadow.camera.far = 50
  directionalLight.shadow.camera.left = -20
  directionalLight.shadow.camera.right = 20
  directionalLight.shadow.camera.top = 20
  directionalLight.shadow.camera.bottom = -20
  scene.add(directionalLight)
}

function createMahjongTable() {
  // Table base
  const tableGeometry = new THREE.CylinderGeometry(7, 7, 0.5, 8)
  const tableMaterial = new THREE.MeshLambertMaterial({ color: 0x8b4513 }) // Brown wood
  const table = new THREE.Mesh(tableGeometry, tableMaterial)
  table.position.y = -2
  table.rotation.y = Math.PI/8
  table.receiveShadow = true
  scene.add(table)

  // Table surface (green felt)
  const surfaceGeometry = new THREE.CylinderGeometry(6.8, 6.8, 0.1, 8)
  const surfaceMaterial = new THREE.MeshLambertMaterial({ color: 0x0a7c4a })
  const surface = new THREE.Mesh(surfaceGeometry, surfaceMaterial)
  surface.position.y = -1.7
  surface.rotation.y = Math.PI/8
  surface.receiveShadow = true
  scene.add(surface)

  // Player positions (4 sides of a square)
  const positions = [
    { x: 0, z: 4, rotation: 0 }, // South (player)
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

function createMahjongTiles() {
  // Create sample mahjong tiles around the table
  const tileGeometry = new THREE.BoxGeometry(0.4, 0.6, 0.25)
  const tileMaterials = [
    new THREE.MeshLambertMaterial({ color: 0xf5f5dc }), // Ivory
    new THREE.MeshLambertMaterial({ color: 0xe6e6fa }), // Lavender
    new THREE.MeshLambertMaterial({ color: 0xffe4e1 }), // Misty rose
  ]

  // Create tile walls for each player
  const positions = [
    { x: 0, z: 5, rotation: 0 }, // South
    { x: 5, z: 0, rotation: Math.PI / 2 }, // East
    { x: 0, z: -5, rotation: Math.PI }, // North
    { x: -5, z: 0, rotation: -Math.PI / 2 }, // West
  ]

  positions.forEach((pos) => {
    for (let i = 0; i < 13; i++) {
      const tile = new THREE.Mesh(
        tileGeometry,
        tileMaterials[i % tileMaterials.length],
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

      mahjongTiles.push(tile)
      scene.add(tile)
    }
  })

  // Center tiles (discarded pile)
  for (let i = 0; i < 12; i++) {
    const tile = new THREE.Mesh(tileGeometry, tileMaterials[0])
    const angle = (i / 12) * Math.PI * 2
    const radius = 1.5
    tile.position.set(Math.cos(angle) * radius, -1.3, Math.sin(angle) * radius)
    tile.rotation.y = angle + Math.PI / 2
    tile.castShadow = true
    tile.receiveShadow = true
    mahjongTiles.push(tile)
    scene.add(tile)
  }
}

function createDealerMarker() {
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
  animationId = requestAnimationFrame(animate)

  // Gentle rotation of dealer marker
  dealerMarker.forEach((marker) => {
    marker.rotation.y += 0.005
  })

  // Subtle floating animation for some tiles
  mahjongTiles.forEach((tile, index) => {
    if (index < 4) {
      // Only animate a few tiles
      tile.position.y = -1.3 + Math.sin(Date.now() * 0.001 + index) * 0.05
    }
  })

  controls.update()
  renderer.render(scene, camera)
}

function onWindowResize() {
  if (!threeCanvas.value || !gameContainer.value) return

  const width = gameContainer.value.clientWidth
  const height = gameContainer.value.clientHeight

  camera.aspect = width / height
  camera.updateProjectionMatrix()
  renderer.setSize(width, height)
}

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId)
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
    <div class="game-ui">
      <div class="game-info">
        <h3>Mahjong Game</h3>
        <p>{{ isFullscreen ? 'Mouse to orbit • Scroll to zoom' : 'Mouse to orbit • Scroll to zoom' }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.game-view-container {
  position: relative;
  width: 100%;
  height: 600px;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
}

.game-view-container.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 1000;
  border-radius: 0;
  cursor: auto;
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
