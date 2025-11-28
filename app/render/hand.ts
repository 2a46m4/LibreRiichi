import * as THREE from 'three'
import { Tile } from '../game/tile'
import { TileObject } from './tile'
import { IAnimation, IAnimationManager, quadratic_interpolator, TileAnimation } from './animation'

// Returns the offset of where the tile should be
function offset(i: number) {
    return (-7 * 0.45) + i * 0.45
}

// A renderable object representing a player's hand
export class Hand extends THREE.Group {
    public array: TileObject[] = []

    constructor(public tiles: Tile[], private animation_manager: IAnimationManager) {
        super()
        tiles.forEach(tile => this.add_tile(tile))
    }

    // Location: the index at which the new tile will end up at
    // Returns the UUID of the new tile object
    add_tile(tile: Tile, location: number = this.array.length) {
        const tile_obj = new TileObject(tile)
        const start = new THREE.Vector3(0, 2, -5)
        const end = new THREE.Vector3(offset(this.array.length), 0, 0)
        const delay = 50 * this.animation_manager.get_animations().length

        this.array.splice(location, 0, tile_obj)
        for (let i = location + 1; i < this.array.length; i++) {
            const start = this.array[i].position
            const end = new THREE.Vector3(offset(i), 0, 0)
            this.array[i].add_animation(new TileAnimation(tile_obj, start, end, quadratic_interpolator, 300, delay))
            this.animation_manager.add_animation(this.array[i])
        }
        tile_obj.add_animation(new TileAnimation(tile_obj, start, end, quadratic_interpolator, 300, delay))
        super.add(tile_obj)
        this.animation_manager.add_animation(this.array[this.array.length - 1])

        return tile_obj.uuid
    }

    find_uuid(uuid: string) {
        return this.array.find(id => id.uuid == uuid)
    }

    find_uuid_index(uuid: string) {
        return this.array.findIndex(id => id.uuid == uuid)
    }

    remove_all() {
        this.array.forEach(obj => {
            obj.removeFromParent()
        })
        this.array = []
    }

    remove_tile_id(tile_id: string) {
        const idx = this.array.findIndex((obj) => obj.uuid == tile_id)
        return this.remove_tile_idx(idx)
    }

    remove_tile_idx(idx: number) {
		if (idx < 0) {
			idx = this.array.length - 1
		}

        const obj = this.array[idx]
        obj.removeFromParent()

        this.array.splice(idx, 1)
        if (idx <= this.array.length - 1) {
            for (let i = idx; i < this.array.length; i++) {

                const tile_obj = this.array[i]
                const start = this.array[i].position
                const end = new THREE.Vector3(offset(i), 0, 0)

                this.array[i].add_animation(new TileAnimation(tile_obj, start, end, quadratic_interpolator, 300, 50 * i))
            }
        }
        return obj
    }
}
