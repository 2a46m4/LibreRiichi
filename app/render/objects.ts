import * as THREE from 'three'
import { Tile } from '../game/tile'
import { tile_width, TileObject } from './tile'
import { quadratic_interpolator, TileLinearAnimation } from './animation'
import { ECS } from './ecs'

export type TileLocation = InHand
export type InHand = { location: "hand", index: number, inhand_index: number }
export type InDiscard = { location: "discard", index: number }
export type InNaki = { location: "naki", index: number, naki_index: number, in_naki_index: number, orientation: 0 | 1 | 2 }

export type NakiCallType = 'pon' | 'chii' | 'ankan' | 'daiminkan'

// Returns the offset of where the tile should be
function offset(i: number) {
	return (-7 * 0.45) + i * 0.45
}

// Creates a tile entity, adds it to the parent
function make_tile_entity(tile: Tile, parent: THREE.Group, parent_id: ECS.EntityID) {
	const tile_entity = ECS.MakeNewTile(tile)
	const parent_component = ECS.MakeComponent(ECS.Parent, tile_entity.entity.uuid, parent_id)
	tile_entity.components.push(parent_component)
	ECS.GlobalRegistry.add_entity(tile_entity.entity, tile_entity.components)

	const tile_object = ECS.GlobalRegistry.find_component(
		tile_entity.entity.uuid,
		ECS.Object.ID
	) as ECS.Object

	parent.add(tile_object.data)
	return tile_entity
}

const dora_tile_position = { x: -5, z: 5, y: -1.3 }

const discard_positions = [
	{ x: 2, z: 0, rotation: Math.PI / 2 }, // East
	{ x: 0, z: -2, rotation: Math.PI }, // North
	{ x: -2, z: 0, rotation: -Math.PI / 2 }, // West
]

const tile_width_gap = tile_width + 0.05

// A renderable object representing a player's hand
export class Hand {
	public array: ECS.EntityID[] = []
	public group_entity: ECS.Entity = new ECS.Entity()
	public group: THREE.Group = new THREE.Group()

	constructor(tiles: Tile[]) {
		const { entity: entity } = ECS.make_new_object(this.group)
		this.group_entity = entity
		ECS.GlobalRegistry.add_entity(this.group_entity, [
			ECS.MakeComponent(ECS.ObjectContainer, this.group_entity.uuid, this.group)])
		this.add_tiles(tiles.map(t => { return { tile: t, location: this.array.length } }))
	}

	set_position(position: THREE.Vector3) {
		this.group.position.copy(position)
	}

	set_rotation(rotation: number) {
		this.group.rotation.y = rotation
	}

	// Location: the index at which the new tile will end up at
	// Returns the UUID of the new tile object
	// TODO: Add a function accepting a Entity parameter
	add_tiles(tiles: { tile: Tile, location: number }[]) {
		let tiles_added = 0
		for (let { tile: tile, location: location } of tiles) {
			const tile_entity = make_tile_entity(tile, this.group, this.group_entity.uuid)
			this.array.splice(location, 0, tile_entity.entity.uuid)

			// Apply an animation to the tiles to the right of the tile
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

			// Apply an animation to the tile
			const animation = new TileLinearAnimation(start, end, quadratic_interpolator, 300, delay)
			ECS.ApplyAnimation(tile_entity.entity.uuid, (dt) => animation.next_step(dt))
		}
	}

	find_index(entity: ECS.EntityID) {
		return this.array.find(v => v == entity)
	}

	remove_all(): ECS.EntityID[] {
		const old: ECS.EntityID[] = []
		this.array.forEach(id => {
			const object = ECS.GlobalRegistry.find_component(id, ECS.Object.ID) as ECS.Object
			object.data.removeFromParent()
			old.push(id)
		})
		this.array = []
		return old
	}

	remove_tile_id(id: ECS.EntityID) {
		const idx = this.array.findIndex((obj) => obj == id)
		return this.remove_tile_idx(idx)
	}

	remove_tile_idx(idx: number) {
		if (idx < 0) {
			idx = this.array.length - 1
		}

		const object = ECS.GlobalRegistry.find_component(this.array[idx], ECS.Object.ID) as ECS.Object
		const worldMatrix = object.data.matrixWorld.clone()
		object.data.removeFromParent()
		object.data.matrixWorld.copy(worldMatrix)

		this.array.splice(idx, 1)
		if (idx <= this.array.length - 1) {
			for (let i = idx; i < this.array.length; i++) {
				const start = object.data.position
				const end = new THREE.Vector3(offset(i), 0, 0)
				const anim = ECS.GlobalRegistry.find_component(this.array[idx], ECS.Animation.ID) as ECS.Animation
				anim.data = (dt) => new TileLinearAnimation(start, end, quadratic_interpolator, 300, 50 * i).next_step(dt)
			}
		}
		return this.array[idx]
	}
}

// TODO
export class NakiGroup {

}

// TODO 
export class Naki {
	public groups: {
		type: NakiCallType
		tiles: TileObject[]
	}[] = []

	constructor() { }

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

export class DiscardPile {
	public tiles: ECS.EntityID[] = []
	public group: THREE.Group = new THREE.Group()

	constructor() { }

	add_to_pile(tile: ECS.EntityID) {
		const tile_obj = ECS.GlobalRegistry.find_component(tile, ECS.Object.ID) as ECS.Object
		// TODO: Animation

		const position = this.compute_next_tile_position()

		this.group.add(tile_obj.data)
		this.tiles.push(tile)
	}

	remove_from_pile() : ECS.EntityID {
		const tile = this.tiles.pop()
		if (tile === undefined) {
			throw new Error("DiscardPile is empty")
		} else {
			return tile
		}
	}

	private compute_next_tile_position(): THREE.Vector3 {
		const offset = ((this.tiles.length % 6) - 3) * tile_width_gap
		const vertical_offset = Math.floor(this.tiles.length / 6) * 0.3
		return new THREE.Vector3(offset, vertical_offset, 0)
	}
}

export class Dora {
	public tile_list: ECS.EntityID[] = []
	public group: THREE.Group = new THREE.Group()
	public e: ECS.Entity = new ECS.Entity()

	constructor() { }

	add_dora(tile: Tile) {
		const tile_entity = make_tile_entity(tile, this.group, this.e.uuid)
		this.tile_list.push(tile_entity.entity.uuid)

		// TODO: Animation
	}
}

export class SelectedTile {
	material: THREE.MeshLambertMaterial = new THREE.MeshLambertMaterial({
		color: 0xffff00,
		transparent: true,
		opacity: 0.5,
	})
	mesh: THREE.Mesh = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.7, 0.26), this.material)

	constructor(scene: THREE.Scene) {
		const e = new ECS.Entity()
		ECS.GlobalRegistry.add_entity(e, [ECS.MakeComponent(ECS.Object, e.uuid, this.mesh)])
		scene.add(this.mesh)
	}

	move(new_position: THREE.Vector3) {
		this.mesh.position.copy(new_position)
	}

	show() {
		this.mesh.visible = true
	}

	hide() {
		this.mesh.visible = false
	}
}
