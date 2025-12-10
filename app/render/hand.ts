import * as THREE from 'three'
import { Tile } from '../game/tile'
import { TileObject } from './tile'
import { IAnimationManager, quadratic_interpolator, TileLinearAnimation } from './animation'
import { ECS } from './ecs'
export type NakiCallType = 'pon' | 'chii' | 'ankan' | 'daiminkan'

// Returns the offset of where the tile should be
function offset(i: number) {
	return (-7 * 0.45) + i * 0.45
}

// A renderable object representing a player's hand
export class Hand extends THREE.Group {
	public array: ECS.EntityID[] = []

	constructor(public tiles: Tile[]) {
		super()
		this.add_tiles(tiles.map(t => { return { tile: t, location: this.array.length } }))
	}

	// Location: the index at which the new tile will end up at
	// Returns the UUID of the new tile object
	add_tiles(tiles: { tile: Tile, location: number }[]) {
		let tiles_added = 0
		for (let { tile: tile, location: location } of tiles) {
			const tile_entity = ECS.MakeNewTile(tile)
			ECS.GlobalRegistry.add_entity(tile_entity.entity, tile_entity.components)

			const start = new THREE.Vector3(0, 2, -5)
			const end = new THREE.Vector3(offset(this.array.length), 0, 0)
			const delay = 50 * tiles_added
			tiles_added += 1
			for (let i = location + 1; i < this.array.length; i++) {
				const start = ECS.GlobalRegistry.find_component(this.array[i], ECS.Object.ID) as ECS.Object
				const end = new THREE.Vector3(offset(i), 0, 0)
				const animation = new TileLinearAnimation(start.data.position, end, quadratic_interpolator, 300, delay)
				ECS.ApplyAnimation(this.array[i], animation.next_step.bind(animation))
			}

			this.array.splice(location, 0, tile_entity.entity.uuid)
			const animation = new TileLinearAnimation(start, end, quadratic_interpolator, 300, delay)
			ECS.ApplyAnimation(tile_entity.entity.uuid, animation.next_step.bind(animation))
			super.add(tile_entity.components[0].component.data as THREE.Mesh)
		}
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
				const start = this.array[i].position
				const end = new THREE.Vector3(offset(i), 0, 0)
				this.animation_manager.animate_object(this.array[i].uuid, new TileLinearAnimation(start, end, quadratic_interpolator, 300, 50 * i))
			}
		}
		return obj
	}
}

export class Naki extends THREE.Group {
	public groups: {
		type: NakiCallType
		tiles: TileObject[]
	}[] = []

	constructor() {
		super()
	}

	add_call(type: NakiCallType, tiles: Tile[]) {
		const tile_objs = tiles.map(tile => new TileObject(tile))
		this.groups.push({
			type: type,
			tiles: tile_objs
		})

		// TODO: Animate naki calls
	}

	convert_pon_to_shouminkan(tile: Tile) {

	}
}
