import * as THREE from 'three'
import { Tile } from '../game/tile'
import { TileObject } from './tile'
import { IAnimationManager, quadratic_interpolator, TileAnimation } from './animation'

function offset(i: number) {
    return i * 0.45
}

export class Hand extends THREE.Group {
    public array: TileObject[] = []

    constructor(public tiles: Tile[], animation_manager: IAnimationManager) {
        super()
        tiles.forEach(tile => this.add_tile(tile, animation_manager))
    }

    add_tile(tile: Tile, animation_manager: IAnimationManager) {
        const tile_obj = new TileObject(tile)
        this.array.push(tile_obj)
        super.add(tile_obj)

        const start = new THREE.Vector3(0, 2, -5)
        const end = new THREE.Vector3(offset(this.array.length), 0, 0)
        const delay = 50 * animation_manager.get_animations().length

        animation_manager.add_animation(new TileAnimation(tile_obj, start, end, quadratic_interpolator, 100, delay))
    }

	find_uuid(uuid:string) {
		return this.array.find(id=>id.uuid==uuid)
	}

	remove_all() {
		this.array.forEach(obj=>{
			obj.removeFromParent()
		})
		this.array = []
	}

    remove_tile_id(tile_id: string) {
        const idx = this.array.findIndex((obj) => obj.uuid == tile_id)
        return this.remove_tile_idx(idx)
    }

    remove_tile_idx(idx: number) {
        const obj = this.array[idx]
        obj.removeFromParent()
        this.array.splice(idx, 1)
        if (idx <= this.array.length - 1) {
            for (let i = idx; i < this.array.length; i++) {
                this.array[i].position.set(offset(i), 0, 0)
            }
        }
        return obj
    }
}
