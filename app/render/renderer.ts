import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import { HiddenTile, Tile } from '../game/tile'
import { tile_width, TileObject } from './tile'
import {
  DiscardPile,
  Dora,
  Hand,
  Naki,
  NakiCallType,
  SelectedTile,
} from './objects'
import { TableIdx } from '../game/arena'

const marker_positions = [
  { x: 0, z: 4, rotation: 0 },
  { x: 4, z: 0, rotation: Math.PI / 2 }, // East
  { x: 0, z: -4, rotation: Math.PI }, // North
  { x: -4, z: 0, rotation: -Math.PI / 2 }, // West
]

const tile_width_gap = tile_width + 0.05

const tile_positions = [
  { x: 5, z: 0, rotation: Math.PI / 2 }, // East
  { x: 0, z: -5, rotation: Math.PI }, // North
  { x: -5, z: 0, rotation: -Math.PI / 2 }, // West
]

const discard_positions = [
  { x: 2, z: 0, rotation: Math.PI / 2 }, // East
  { x: 0, z: -2, rotation: Math.PI }, // North
  { x: -2, z: 0, rotation: -Math.PI / 2 }, // West
]

const player_position = { x: 0, z: 5, rotation: 0 }

const default_tiles = [0, 1, 2, 3, 4, 5, 6, 7, 8, 16, 17, 18, 19].map(
  (i) => new Tile(i),
)

export class ThreeJSRenderer {
	scene: THREE.Scene
	camera: THREE.PerspectiveCamera
	renderer: THREE.WebGLRenderer
	controls: OrbitControls
	lights: THREE.Light[] = []

	table: THREE.Group

	hands: Hand[]
	naki_calls: Naki[] = []
	dora_tiles: Dora = new Dora()
	discard_pile: DiscardPile[]
	selected_tile: SelectedTile
	selector: Selector

	constructor(canvas: HTMLCanvasElement) {
		// Scene
		{
			this.scene = new THREE.Scene()
			this.scene.background = new THREE.Color(0xffffff)
		}

		// Camera
		{
			this.camera = new THREE.PerspectiveCamera(
				35,
				canvas.clientWidth / canvas.clientHeight,
				0.1,
				1000,
			)
			this.camera.position.set(0, 4, 20)
			this.camera.lookAt(0, 0, 0)
		}

		// Renderer
		{
			this.renderer = new THREE.WebGLRenderer({
				canvas: canvas,
				antialias: true,
			})
			this.renderer.setSize(canvas.clientWidth, canvas.clientHeight)
			this.renderer.setPixelRatio(window.devicePixelRatio)
			this.renderer.shadowMap.enabled = true
			this.renderer.shadowMap.type = THREE.PCFSoftShadowMap
		}

		// Controls
		{
			this.controls = new OrbitControls(this.camera, this.renderer.domElement)
			this.controls.enableDamping = true
			this.controls.dampingFactor = 0.05
			this.controls.maxPolarAngle = Math.PI / 2
			this.controls.minPolarAngle = 0
			this.controls.maxAzimuthAngle = Math.PI / 12
			this.controls.minAzimuthAngle = -Math.PI / 12
			this.controls.maxDistance = 20
			this.controls.minDistance = 12
			this.controls.enablePan = false
		}

		// Lights
		{
			this.lights.push(new THREE.AmbientLight(0x404040, 0.6))
			const directional_light = new THREE.DirectionalLight(0xffffff, 1)
			directional_light.position.set(10, 15, 5)
			directional_light.castShadow = true
			directional_light.shadow.mapSize.width = 512
			directional_light.shadow.mapSize.height = 512
			directional_light.shadow.camera.near = 15
			directional_light.shadow.camera.far = 25
			directional_light.shadow.camera.left = -5
			directional_light.shadow.camera.right = 6
			directional_light.shadow.camera.top = 4.2
			directional_light.shadow.camera.bottom = -6
			directional_light.shadow.bias = 0.0003
			this.lights.push(directional_light)
			this.scene.add(...this.lights)
		}

		// Table
		{
			const table_geometry = new THREE.BoxGeometry(13, 13, 0.5)
			const table_material = new THREE.MeshLambertMaterial({ color: 0x8b4513 }) // Brown wood
			const table = new THREE.Mesh(table_geometry, table_material)
			table.position.y = -2
			table.rotation.x = Math.PI / 2
			table.receiveShadow = true

			const surface_geometry = new THREE.BoxGeometry(12.8, 12.8, 0.1)
			const surface_material = new THREE.MeshLambertMaterial({
				color: 0x0a7c4a,
			})
			const surface = new THREE.Mesh(surface_geometry, surface_material)
			surface.position.y = -1.7
			surface.rotation.x = Math.PI / 2
			surface.receiveShadow = true

			this.table = new THREE.Group()
			this.table.add(table, surface)

			marker_positions.forEach((pos, index) => {
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
				this.table.add(marker)
			})

			this.scene.add(this.table)
		}

		// Mahjong tiles
		{
			this.hands = []

			// Create demo tiles
			this.hands[0] = new Hand([])
			this.hands[0].set_position(
				new THREE.Vector3(player_position.x, -1.3, player_position.z),
			)
			this.hands[0].set_rotation(player_position.rotation)
			this.hands[0].add_tiles(
				default_tiles.map((v, i) => {
					return { tile: v, location: i }
				}),
			)
			this.scene.add(this.hands[0].group)

			// Create blank tile walls for the other players
			let blank_tile = HiddenTile
			for (let i = 1; i < 4; i++) {
				const other_hand = new Hand([])
				for (let j = 0; j < 13; j++) {
					other_hand.add_tiles([{ tile: blank_tile, location: j }])
					other_hand.set_rotation(tile_positions[i - 1].rotation)
					other_hand.set_position(
						new THREE.Vector3(
							tile_positions[i - 1].x,
							-1.3,
							tile_positions[i - 1].z,
						),
					)
				}
				this.hands[i] = other_hand
				this.scene.add(other_hand.group)
			}
		}

		// Discard piles
		{
			this.discard_pile = new Array(4)
		}

		// Dealer marker
		{
			const dealerMarker: THREE.Mesh[] = []

			// Dealer button/marker
			const markerGeometry = new THREE.CylinderGeometry(0.3, 0.3, 0.1, 8)
			const markerMaterial = new THREE.MeshLambertMaterial({ color: 0xff6b6b })
			const marker = new THREE.Mesh(markerGeometry, markerMaterial)
			marker.position.set(2, -1.2, 2)
			marker.castShadow = true
			dealerMarker.push(marker)
			this.scene.add(marker)

			// Dealer marker text indicator
			const textGeometry = new THREE.RingGeometry(0.1, 0.2, 6)
			const textMaterial = new THREE.MeshLambertMaterial({ color: 0xffffff })
			const textRing = new THREE.Mesh(textGeometry, textMaterial)
			textRing.position.set(2, -1.1, 2)
			textRing.rotation.x = -Math.PI / 2
			dealerMarker.push(textRing)
			this.scene.add(textRing)
		}

		this.selected_tile = new SelectedTile()
		this.scene.add(this.selected_tile)
		this.selector = new Selector(this.hands[0])
	}

	get_random_tile_in_hand(player_idx: number): TileObject {
		return this.hands[player_idx].get_random_tile()
	}

	//TODO
	toss(player_idx: number, tile: TileObject): void {
		this.hands[player_idx].remove_tile(tile)
		this.discard_pile[player_idx].add_to_pile(tile)
	}

	animate_frame(dt: number): void {
		this.controls.update()
		this.renderer.render(this.scene, this.camera)

		// Animation
		for (let hand of this.hands) {
			hand.animate(dt)
		}

		for (let naki of this.naki_calls) {
			naki.animate(dt)
		}

		this.dora_tiles.animate(dt)

		// Selection
	}

	stop(): void {
		this.renderer.dispose()
		this.scene.traverse((object) => {
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

	window_resize(width: number, height: number) {
		this.camera.aspect = width / height
		this.camera.updateProjectionMatrix()
		this.renderer.setSize(width, height)
	}

	clear_tiles(): void {
		this.hands[0].remove_all()
	}

	clear_tiles_on(idx: number): void {
		this.hands[idx].remove_all()
	}

	add_tile_to_end(tile: Tile, player_idx: number = 0) {
		const end_loc = this.hands[player_idx].array.length - 1
		this.hands[player_idx].add_tiles([{ tile: tile, location: end_loc }])
	}

	remove_tile(id: number, player_idx: number | undefined): void {
		if (player_idx === undefined) {
			player_idx = 0
		}

		this.hands[player_idx].remove_tile_idx(id)
	}

	add_dora(tile: Tile): void {
		this.dora_tiles.add_dora(tile)
	}

	// TODO
	naki_call(called_by: TableIdx, type: NakiCallType): void {}
}

export class Selector {
	pointer = new THREE.Vector2()
	raycaster = new THREE.Raycaster()

	private on_move(event: MouseEvent) {
		this.pointer.x = (event.clientX / window.innerWidth) * 2 - 1
		this.pointer.y = -(event.clientY / window.innerHeight) * 2 + 1
	}

	constructor(public player_hand: Hand) {
		window.addEventListener('pointermove', this.on_move)
	}

	get_selection(): TileObject | null {
		const results = this.raycaster.intersectObjects(this.player_hand.array)
		if (results.length === 0) {
			return null
		} else {
			return results[0].object as TileObject
		}
	}
}
