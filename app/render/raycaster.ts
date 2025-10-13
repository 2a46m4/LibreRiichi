import * as THREE from 'three'
import { IRenderer, ThreeJSRenderer } from "./renderer_setup";

export interface Selection {
    id: string
}

export interface Selector {
    get_selections(): Selection[]

    stop(): void
}

export function filter_tiles(mahjong_tiles: Map<string, THREE.Mesh>) {
    return (obj: THREE.Intersection) => {
        return mahjong_tiles.has(obj.object.uuid)
    }
}

export class Raycaster implements Selector {
    raycaster = new THREE.Raycaster()
    pointer = new THREE.Vector2()
    filter: (obj: THREE.Intersection) => boolean = () => true
    renderer: ThreeJSRenderer

    constructor(renderer: IRenderer) {
        if (renderer instanceof ThreeJSRenderer) {
            this.renderer = renderer
            this.filter = filter_tiles(this.renderer.tiles)
            window.addEventListener('pointermove', (event) => this.on_pointer_move(event))
        } else {
            throw new Error("Wrong renderer type")
        }
    }

    get_selections(): { data: THREE.Intersection, id: string }[] {
        this.raycaster.setFromCamera(this.pointer, this.renderer.camera)
        return this.raycaster.intersectObjects(this.renderer.scene.children)
            .filter(this.filter)
            .map(obj => ({ data: obj, id: obj.object.uuid }))
    }

    stop() {
        window.removeEventListener('pointermove', (event) => this.on_pointer_move(event))
    }

    // Update pointer's NDC coordinate system
    on_pointer_move(event: PointerEvent) {
        this.pointer.x = (event.clientX / window.innerWidth) * 2 - 1
        this.pointer.y = -(event.clientY / window.innerHeight) * 2 + 1
    }
}
