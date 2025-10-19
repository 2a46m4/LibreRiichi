import * as THREE from 'three'
import { Tile } from '../game/tile'
import { TileObject } from './tile'
import { IAnimation, IAnimationManager, MultipleTileAnimation, quadratic_interpolator, TileAnimation } from './animation'

function offset(i: number) {
    return (-7 * 0.45) + i * 0.45
}

export class Hand extends THREE.Group {
    public array: {
        tile: TileObject,
        animation: MultipleTileAnimation
    }[] = []

    constructor(public tiles: Tile[], private animation_manager: IAnimationManager) {
        super()
        tiles.forEach(tile => this.add_tile(tile))
    }

    add_tile(tile: Tile) {
        const tile_obj = new TileObject(tile)
        const start = new THREE.Vector3(0, 2, -5)
        const end = new THREE.Vector3(offset(this.array.length), 0, 0)
        const delay = 50 * this.animation_manager.get_animations().length

        this.array.push({
            tile: tile_obj,
            animation: new MultipleTileAnimation([
                new TileAnimation(tile_obj, start, end, quadratic_interpolator, 300, delay)
            ])
        })
        super.add(tile_obj)
        this.animation_manager.add_animation(this.array[this.array.length - 1].animation)
    }

    find_uuid(uuid: string) {
        return this.array.find(id => id.tile.uuid == uuid)
    }

    remove_all() {
        this.array.forEach(obj => {
            obj.tile.removeFromParent()
        })
        this.array = []
    }

    remove_tile_id(tile_id: string) {
        const idx = this.array.findIndex((obj) => obj.tile.uuid == tile_id)
        return this.remove_tile_idx(idx)
    }

    remove_tile_idx(idx: number) {
        const obj = this.array[idx]
        obj.tile.removeFromParent()

        this.array.splice(idx, 1)
        if (idx <= this.array.length - 1) {
            for (let i = idx; i < this.array.length; i++) {

                const tile_obj = this.array[i].tile
                const start = this.array[i].tile.position
                const end = new THREE.Vector3(offset(i), 0, 0)

                if (this.array[i].animation.finished()) {
                    this.array[i].animation = new MultipleTileAnimation([
                        new TileAnimation(tile_obj, start, end, quadratic_interpolator, 300, 50 * i)
                    ])
                    this.animation_manager.add_animation(this.array[i].animation)
                } else {
                    this.array[i].animation.add(
                        new TileAnimation(tile_obj, start, end, quadratic_interpolator, 300, 50 * i)
                    )
                }

            }
        }
        return obj
    }
}
