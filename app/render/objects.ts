import * as THREE from 'three'
import { Tile } from '../game/tile'
import { tile_width, TileObject } from './tile'
import { quadratic_interpolator, TileLinearAnimation } from './animation'

export type TileLocation = InHand
export type InHand = { location: 'hand'; index: number; inhand_index: number }
export type InDiscard = { location: 'discard'; index: number }
export type InNaki = {
	location: 'naki'
	index: number
	naki_index: number
	in_naki_index: number
	orientation: 0 | 1 | 2
}

export type NakiCallType = 'pon' | 'chii' | 'ankan' | 'daiminkan'

// Returns the offset of where the tile should be
function offset(i: number) {
	return -7 * 0.45 + i * 0.45
}

const dora_tile_position = { x: -5, z: 5, y: -1.3 }

const discard_positions = [
	{ x: 2, z: 0, rotation: Math.PI / 2 }, // East
	{ x: 0, z: -2, rotation: Math.PI }, // North
	{ x: -2, z: 0, rotation: -Math.PI / 2 }, // West
]

const tile_width_gap = tile_width + 0.05

export namespace Move {
	function hand_to_discard(hand: Hand, discard: DiscardPile) {
		throw new Error("TODO")
	}

	function discard_to_naki(discard: DiscardPile, naki_group: SingleNakiGroup) {
		throw new Error("TODO")
	}
}

// A renderable object representing a player's hand
export class Hand {
	public array: TileObject[] = []
	public group = new THREE.Group()

	constructor(tiles: Tile[]) {
		this.add_tiles_to_last(tiles)
	}

	set_position(position: THREE.Vector3) {
		this.group.position.copy(position)
	}

	set_rotation(rotation: number) {
		this.group.rotation.y = rotation
	}

	add_tiles_to_last(tiles: Tile[]) {
		this.add_tiles(
			tiles.map((t) => {
				return { tile: t, location: this.array.length }
			})
		)
	}

	// Location: the index at which the new tile will end up at
	// Returns the UUID of the new tile object
	add_tiles(tiles: { tile: Tile; location: number }[]) {
		let tiles_added = 0
		for (let { tile: tile, location: location } of tiles) {
			const tile_obj = new TileObject(tile)
			this.array.splice(location, 0, tile_obj)

			// Apply an animation to the tiles to the right of the tile
			const start = new THREE.Vector3(0, 2, -5)
			const end = new THREE.Vector3(offset(this.array.length), 0, 0)
			const delay = 50 * tiles_added
			tiles_added += 1
			for (let i = location + 1; i < this.array.length; i++) {
				const start = this.array[i].position
				const end = new THREE.Vector3(offset(i), 0, 0)
				this.array[i].animation = new TileLinearAnimation(
					start,
					end,
					quadratic_interpolator,
					300,
					delay,
				)
			}

			tile_obj.animation = new TileLinearAnimation(
				start,
				end,
				quadratic_interpolator,
				300,
				delay,
			)
			this.group.add(tile_obj)
		}
	}

	find_index(entity: TileObject) {
		return this.array.find((v) => v == entity)
	}

	remove_all(): TileObject[] {
		const old: TileObject[] = []
		this.array.forEach((obj) => {
			obj.removeFromParent()
			old.push(obj)
		})
		this.array = []
		return old
	}

	remove_tile(obj: TileObject) {
		const idx = this.array.findIndex((obj_arr) => obj_arr == obj)
		return this.remove_tile_idx(idx)
	}

	remove_tile_id(id: string) {
		const idx = this.array.findIndex((obj) => obj.uuid == id)
		return this.remove_tile_idx(idx)
	}

	remove_tile_idx(idx: number) {
		if (idx < 0) {
			idx = this.array.length - 1
		}

		const object = this.array[idx]
		const worldMatrix = object.matrixWorld.clone()
		object.removeFromParent()
		object.matrixWorld.copy(worldMatrix)

		this.array.splice(idx, 1)
		if (idx <= this.array.length - 1) {
			for (let i = idx; i < this.array.length; i++) {
				const start = this.array[i].position
				const end = new THREE.Vector3(offset(i), 0, 0)
				this.array[i].animation
				new TileLinearAnimation(start, end, quadratic_interpolator, 300, 50 * i)
			}
			return this.array[idx]
		}
	}


	get_random_tile(): TileObject {
		const rand_idx = Math.floor(Math.random() * this.array.length)
		return this.array[rand_idx]
	}

	animate(dt: number) {
		for (let obj of this.array) {
			obj.animate(dt)
		}
	}
}

export class SingleNakiGroup {
	group = new THREE.Group()

	constructor(public type: NakiCallType, public tiles: TileObject[]) {
		// TODO: Animate naki calls
	}

	animate(dt: number) {
		this.tiles.forEach(t => t.animate(dt))
	}
}

export class Naki {
	public groups: SingleNakiGroup[] = []

	constructor() { }

	convert_pon_to_shouminkan(tile: Tile) { }

	animate(dt: number) {
		for (let obj of this.groups) {
			obj.animate(dt)
		}
	}
}

export class DiscardPile {
	public tiles: TileObject[] = []
	public group: THREE.Group = new THREE.Group()

	constructor() { }

	add_to_pile(tile: TileObject) {
		// TODO: Animation

		const position = this.compute_next_tile_position()

		this.group.add(tile)
		this.tiles.push(tile)
	}

	remove_from_pile(): TileObject {
		const tile = this.tiles.pop()
		if (tile === undefined) {
			throw new Error('DiscardPile is empty')
		} else {
			return tile
		}
	}

	private compute_next_tile_position(): THREE.Vector3 {
		const offset = ((this.tiles.length % 6) - 3) * tile_width_gap
		const vertical_offset = Math.floor(this.tiles.length / 6) * 0.3
		return new THREE.Vector3(offset, vertical_offset, 0)

		// // Start a new row
		// if (pile.length % 6 === 0) {
		//   const offset = ((pile.length % 6) - 3) * tile_width_gap
		//   const vertical_offset = Math.floor(pile.length / 6) * 0.3
		//   tile_obj.position.set(
		//     discard_positions[player_idx].x -
		//       offset * Math.sin(discard_positions[player_idx].rotation),
		//     -1.3 + vertical_offset,
		//     discard_positions[player_idx].z -
		//       offset * Math.cos(discard_positions[player_idx].rotation),
		//   )
		//   tile_obj.rotation.y = discard_positions[player_idx].rotation
		// }
		// pile.push(tile_obj)
		// this.last_discard = [player_idx, pile.length - 1]
		// this.scene.add(tile_obj)
	}
}

export class Dora {
	public tile_list: TileObject[] = []
	public group: THREE.Group = new THREE.Group()

	constructor() { }

	add_dora(tile: Tile) {
		const tile_obj = new TileObject(tile)
		this.tile_list.push(tile_obj)

		// TODO: Animation
	}

	animate(dt: number) {
		this.tile_list.forEach(t => t.animate(dt))
	}
}

export class SelectedTile extends THREE.Mesh {
	constructor() {
		const material = new THREE.MeshLambertMaterial({
			color: 0xffff00,
			transparent: true,
			opacity: 0.5,
	})
	super(
		new THREE.BoxGeometry(0.5, 0.7, 0.26),
		material)
}

	move(new_position: THREE.Vector3) {
		this.position.copy(new_position)
	}

	show() {
		this.visible = true
	}

	hide() {
		this.visible = false
	}
}
