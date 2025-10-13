import * as THREE from 'three'
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls.js";
import { Tile, TileValue } from "../game/tile";
import { TileObject } from "./tile";
import { Selection } from "./raycaster";
import { Animation, defaultAnimationManager } from "./animation";

export interface IRenderer {
    animate_frame(): void

    stop(): void

    window_resize(width: number, height: number): void

    render_selection(selection: Selection[]): void

    clear_tiles(): void

    add_tile(tile: Tile, location?: number): string

    remove_tile(id: string): void

    add_dora(tile: Tile): void

    get_selection(): { tile: Tile, id: string } | null

    render_draw(player_idx: number, tile?: Tile): void
}

const marker_positions = [
    { x: 0, z: 4, rotation: 0 },
    { x: 4, z: 0, rotation: Math.PI / 2 }, // East
    { x: 0, z: -4, rotation: Math.PI }, // North
    { x: -4, z: 0, rotation: -Math.PI / 2 }, // West
]

const tile_positions = [
    { x: 5, z: 0, rotation: Math.PI / 2 }, // East
    { x: 0, z: -5, rotation: Math.PI }, // North
    { x: -5, z: 0, rotation: -Math.PI / 2 }, // West
]

const dora_tile_position = { x: -5, z: 5, y: -1.3 }

const player_position = { x: 0, z: 5, rotation: 0 }

const default_tiles = [0, 1, 2, 3, 4, 5, 6, 7, 8, 16, 17, 18, 19].map(
    (i) => new Tile(i),
)

export class ThreeJSRenderer implements IRenderer {

    scene: THREE.Scene
    camera: THREE.PerspectiveCamera
    renderer: THREE.WebGLRenderer
    controls: OrbitControls
    lights: THREE.Light[] = []

    table: THREE.Group
    tiles: Map<string, TileObject>
    other_tiles: THREE.Group[]
    dora_tiles: THREE.Mesh[] = []

    animations_in_flight: Animation[] = []

    animation_id = 0
    selection: {
        material: THREE.MeshLambertMaterial
        mesh: THREE.Mesh
        tile: TileObject | null
    }

    constructor(canvas: HTMLCanvasElement) {
        // Scene
        {
            this.scene = new THREE.Scene()
            this.scene.background = new THREE.Color(0xffffff) // Mahjong table green
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
            // Create demo tiles
            this.tiles = new Map<string, TileObject>()
            for (const [i, tile] of default_tiles.entries()) {
                let tile_obj = new TileObject(tile)

                const offset_x = (i - 6) * 0.45
                tile_obj.position.set(
                    player_position.x + Math.cos(player_position.rotation) * offset_x,
                    -1.3,
                    player_position.z + Math.sin(player_position.rotation) * offset_x,
                )
                tile_obj.rotation.y = player_position.rotation

                this.tiles.set(tile_obj.uuid, tile_obj)
                this.scene.add(tile_obj)
            }

            // Create blank tile walls for the other players
            this.other_tiles = []
            let blank_tile = new Tile(TileValue.Hidden)
            tile_positions.forEach((pos) => {
                const blank_tiles = new THREE.Group()
                for (let i = 0; i < 13; i++) {
                    const blank_tile_obj = new TileObject(blank_tile)
                    const offset_x = (i - 6) * 0.45
                    blank_tile_obj.position.set(
                        pos.x + Math.cos(pos.rotation) * offset_x,
                        -1.3,
                        pos.z + Math.sin(pos.rotation) * offset_x,
                    )
                    blank_tile_obj.rotation.y = pos.rotation
                    blank_tiles.add(blank_tile_obj)
                }
                this.other_tiles.push(blank_tiles)
                this.scene.add(blank_tiles)
            })
        }

        // Discard piles
        {
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
                this.scene.add(tile_obj)
            }
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
    }

    animate_frame(): void {
        this.controls.update()
        this.renderer.render(this.scene, this.camera)
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

    render_selection(selection: Selection[]) {
        if (selection.length > 0) {
            let tile = this.tiles.get(selection[0].id)
            if (tile === undefined) {
                throw new Error('Tile not found')
            }

            this.selection.mesh.visible = true
            this.selection.mesh.position.copy(tile.position)
            this.selection.tile = tile as TileObject
        } else {
            this.selection.mesh.visible = false
            return
        }
    }

    clear_tiles(): void {
        this.tiles.forEach(mesh => this.scene.remove(mesh))
        this.tiles.clear()
    }

    add_tile(tile: Tile, location: number = this.tiles.size): string {
        let tile_obj = new TileObject(tile)
        let offset_x = (location - 6) * 0.45
        tile_obj.position.set(
            player_position.x + Math.cos(player_position.rotation) * offset_x,
            -1.3,
            player_position.z + Math.sin(player_position.rotation) * offset_x,
        )
        tile_obj.rotation.y = player_position.rotation
        this.tiles.set(tile_obj.uuid, tile_obj)
        this.scene.add(tile_obj)
        return tile_obj.uuid
    }

    remove_tile(id: string): void {
        let tile_obj = this.tiles.get(id)
        if (tile_obj === undefined) {
            throw new Error('Tile not found')
        }
        this.scene.remove(tile_obj)
        this.tiles.delete(id)
    }

    add_dora(tile: Tile): void {
        const dora_tile = new TileObject(tile)
        this.dora_tiles.push(dora_tile)
        dora_tile.position.set(
            dora_tile_position.x, dora_tile_position.y, dora_tile_position.z
        )
        this.scene.add(dora_tile)
    }

    get_selection(): { tile: Tile; id: string } | null {
        if (this.selection.tile === null) {
            return null
        }
        return {
            tile: this.selection.tile.tile,
            id: this.selection.tile.uuid,
        }
    }

    // Assumes clockwise index, with our player starting at 0
    render_draw(player_idx: number, tile?: Tile): void {
        if (player_idx === 0) {
            if (tile === null || tile === undefined) {
                throw new Error("Tile can't be null")
            }

            this.add_tile(tile)
        } else {

        }
    }
}
