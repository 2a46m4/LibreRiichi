import * as THREE from 'three'
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls.js";
import { HiddenTile, Tile } from "../game/tile";
import { tile_width, TileObject } from "./tile";
import { Raycaster, Selection, Selector } from "./raycaster";
import { AnimationManager, IAnimationManager, quadratic_interpolator, TileAnimation } from "./animation";
import { Hand, Naki, NakiCallType } from './hand';
import { TableIdx } from '../game/arena';

// Animates and manages actions
export interface IActionAnimator {
    clear_tiles(): void
    clear_tiles_on(idx: number): void
    add_tile(tile: Tile, player_idx?: number, location?: number): string
    remove_tile(idx: string | number, player_idx?: TableIdx): void
    add_dora(tile: Tile): void
    draw(player_idx: number, tile?: Tile): void
    toss(player_idx: number, tile: Tile): void
    select(selections: Selection[]): void
    naki_call(called_by: TableIdx, type: NakiCallType): void
}

// Callable
export interface IRenderer {
    animate_frame(dt: number): void
    stop(): void
    window_resize(width: number, height: number): void
}

export interface ISelectionManager {
    get_selection(): { tile: Tile, id: string, location: number } | null
}

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

const dora_tile_position = { x: -5, z: 5, y: -1.3 }

const player_position = { x: 0, z: 5, rotation: 0 }

const default_tiles = [0, 1, 2, 3, 4, 5, 6, 7, 8, 16, 17, 18, 19].map(
    (i) => new Tile(i),
)

export class ThreeJSRenderer implements IRenderer, IActionAnimator, ISelectionManager {

    scene: THREE.Scene
    camera: THREE.PerspectiveCamera
    renderer: THREE.WebGLRenderer
    controls: OrbitControls
    lights: THREE.Light[] = []

    table: THREE.Group
    hands: Hand[]
    naki_calls: Naki[] = []
    dora_tiles: THREE.Mesh[] = []

    discard_pile: TileObject[][]
    last_discard: [number, number] = [-1, -1]

    selection: {
        material: THREE.MeshLambertMaterial
        mesh: THREE.Mesh
        tile: TileObject | null
    }

    selector: Selector

    animation_manager: IAnimationManager = new AnimationManager()

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
            const surface_material = new THREE.MeshLambertMaterial({ color: 0x0a7c4a })
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
            this.hands[0] = new Hand([], this.animation_manager)
            this.hands[0].position.set(player_position.x, -1.3, player_position.z)
            this.hands[0].rotation.y = player_position.rotation
            default_tiles.forEach(tile => {
                this.hands[0].add_tile(tile)
            })
            this.scene.add(this.hands[0])

            // Create blank tile walls for the other players
            let blank_tile = HiddenTile
            for (let i = 1; i < 4; i++) {
                const other_hand = new Hand([], this.animation_manager)
                for (let j = 0; j < 13; j++) {
                    other_hand.add_tile(blank_tile)
                    other_hand.rotation.y = tile_positions[i - 1].rotation
                    other_hand.position.set(tile_positions[i - 1].x, -1.3, tile_positions[i - 1].z)
                }
                this.hands[i] = other_hand
                this.scene.add(other_hand)
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

        // Selection
        {
            let material = new THREE.MeshLambertMaterial({
                color: 0xffff00,
                transparent: true,
                opacity: 0.5,
            })
            let mesh = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.7, 0.26), material)
            mesh.visible = false
            let selected_tile: TileObject | null = null

            this.selection = {
                material: material,
                mesh: mesh,
                tile: selected_tile,
            }

            this.scene.add(this.selection.mesh)
        }

        // Selection manager
        {
            this.selector = new Raycaster(this)
        }
    }

    //TODO
    toss(player_idx: number, tile: Tile): void {
        let pile = this.discard_pile[player_idx]

        let tile_obj = new TileObject(tile)

        // Start a new row
        if (pile.length % 6 === 0) {
            const offset = ((pile.length % 6) - 3) * tile_width_gap
            const vertical_offset = Math.floor(pile.length / 6) * 0.3
            tile_obj.position.set(
                discard_positions[player_idx].x - offset * Math.sin(discard_positions[player_idx].rotation),
                -1.3 + vertical_offset,
                discard_positions[player_idx].z - offset * Math.cos(discard_positions[player_idx].rotation)
            )
            tile_obj.rotation.y = discard_positions[player_idx].rotation
        }
        pile.push(tile_obj)
        this.last_discard = [player_idx, pile.length - 1]
        this.scene.add(tile_obj)
    }

    animate_frame(dt: number): void {
        this.controls.update()
        this.renderer.render(this.scene, this.camera)
        this.animation_manager.animate_step(dt)
        this.select(this.selector.get_selections())
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
        this.selector.stop()
    }

    window_resize(width: number, height: number) {
        this.camera.aspect = width / height
        this.camera.updateProjectionMatrix()
        this.renderer.setSize(width, height)
    }

    select(selection: Selection[]) {
        if (selection.length > 0) {
            let tile = this.hands[0].find_uuid(selection[0].id)
            if (tile === undefined) {
                throw new Error('Tile not found')
            }

            this.selection.mesh.visible = true
            tile.getWorldPosition(this.selection.mesh.position)
            this.selection.tile = tile as TileObject
        } else {
            this.selection.mesh.visible = false
            this.selection.tile = null
            return
        }
    }

    clear_tiles(): void {
        this.hands[0].remove_all()
    }

    clear_tiles_on(idx: number): void {
        this.hands[idx].remove_all()
    }

    add_tile(tile: Tile, player_idx: number = 0, add_location: number = -1): string {
        if (add_location >= 0) {
            return this.hands[player_idx].add_tile(tile, add_location)
        } else {
            return this.hands[player_idx].add_tile(tile)
        }
    }

    remove_tile(id: string | number, player_idx: number | undefined): void {
        if (player_idx === undefined) {
            player_idx = 0
        }

        if (typeof id === "string") {
            this.hands[player_idx].remove_tile_id(id)
        } else {
            this.hands[player_idx].remove_tile_idx(id)
        }
    }

    add_dora(tile: Tile): void {
        const dora_tile = new TileObject(tile)
        this.dora_tiles.push(dora_tile)
        dora_tile.position.set(
            dora_tile_position.x, dora_tile_position.y, dora_tile_position.z
        )
        this.scene.add(dora_tile)
    }

    get_selection(): { tile: Tile; id: string, location: number } | null {
        if (this.selection.tile === null) {
            return null
        }
        return {
            tile: this.selection.tile.tile,
            id: this.selection.tile.uuid,
            location: this.hands[0].find_uuid_index(this.selection.tile.uuid)
        }
    }

    // Assumes clockwise index, with our player starting at 0
    draw(player_idx: number, tile?: Tile): void {
        if (player_idx === 0) {
            if (tile === null || tile === undefined) {
                throw new Error("Tile can't be null")
            }

            this.add_tile(tile, player_idx)
        } else {
            if (tile === undefined) {
                tile = HiddenTile
            }
            this.add_tile(tile, player_idx)
        }
    }

    // TODO
    naki_call(called_by: TableIdx, type: NakiCallType): void {

    }
}