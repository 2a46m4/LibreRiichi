import { v4 as uuidv4 } from 'uuid'
import * as THREE from "three"
import { Tile } from '../game/tile'
import { ClickEventBus } from '../messaging/event_handler'
import { TileObject } from './tile'

export namespace ECS {
    export type EntityID = string & { readonly __brand: unique symbol }
    export type ComponentID = string & { readonly __brand: unique symbol }
    export type SystemID = string & { readonly __brand: unique symbol }

    export class Registry {
	entities: Map<EntityID, {
	    components: ComponentID[],
	}> = new Map()
	components: Map<ComponentID, ComponentArray> = new Map()

	add_entity(entity: Entity, components: { id: ComponentID, component: Component }[]) {
	    this.entities.set(entity.uuid, { components: components.map(c => c.id) })
	    for (let component of components) {
		const c = this.components.get(component.id)
		c?.add_entity(entity, component.component)
	    }
	}

	update_entity_component(entity: EntityID, ...components: { id: ComponentID, component: Component }[]) {
	    const existing = this.entities.get(entity)
	    if (existing === undefined) {
		throw new Error("Entity not found")
	    }
	    for (let component of components) {
		const c = this.components.get(component.id)
		c?.update_entity(entity, component.component)
	    }
	}

	add_component_to_entity(entity: Entity, ...components: { id: ComponentID, component: Component }[]) {
	    const existing = this.entities.get(entity.uuid)
	    if (existing === undefined) {
		this.entities.set(entity.uuid, { components: components.map(c => c.id) })
		return
	    }
	    const list = new Set(existing.components).union(new Set(components.map(c => c.id))).values().toArray()
	    this.entities.set(entity.uuid, { components: list })
	    for (let component of components) {
		const c = this.components.get(component.id)
		c?.add_entity(entity, component.component)
	    }
	}

	remove_component_from_entity(entity: Entity, ...components: { id: ComponentID, component: Component }[]) {
	    const existing = this.entities.get(entity.uuid)
	    if (existing === undefined) {
		return
	    }
	    // Assumes that parameter passed in is a subset of the existing components
	    const list = new Set(existing.components).difference(new Set(components.map(c => c.id))).values().toArray()
	    this.entities.set(entity.uuid, { components: list })
	    for (let component of components) {
		const c = this.components.get(component.id)
		c?.drop_entity(entity.uuid)
	    }
	}

	remove_entity(entity: EntityID) {
	    this.components.forEach(c => c.drop_entity(entity))
	}

	add_component(): ComponentID {
	    const arr = new ComponentArray()
	    this.components.set(arr.component_uuid, arr)
	    return arr.component_uuid
	}

	find_component(entity: EntityID, component: ComponentID): Component {
	    const c = this.components.get(component)?.get_entity(entity)
	    if (c === undefined)
		throw new Error("Couldn't find component in entity")
	    return c
	}

	run_system(sys: System): void {
	    if (sys === undefined) {
		throw new Error("Could not find system")
	    }

	    sys.component_ids
		.map(c => new Set(this.components.get(c)?.get_entities()))
		.reduce((acc, cur) => acc.intersection(cur))
		.values() // All entities we need to iterate through
		.map(e => sys.component_ids.map(c => { // Get all components of an entity
		const cmp = this.components.get(c)
		if (cmp === undefined) {
		    throw new Error("Cannot find component")
		}
		return cmp.get_entity(e)
	    }))
		.forEach(val => sys.apply(...val))
	}
    }

    export class System {
	public uuid: SystemID = uuidv4() as SystemID
	constructor(public component_ids: ComponentID[], public apply: (...components: Component[]) => void) { }
    }

    export class Entity {
	public uuid: EntityID = uuidv4() as EntityID
    }

    export class ComponentArray {
	public component_uuid: ComponentID = uuidv4() as ComponentID
	public array: Component[] = []
	public map: Map<EntityID, number> = new Map()

	get_component(entity: Entity): Component {
	    const idx = this.map.get(entity.uuid)
	    if (idx === undefined) {
		throw new Error("Couldn't find entity")
	    }
	    return this.array[idx]
	}

	update_entity(entity: EntityID, value: any) {
	    const idx = this.map.get(entity)
	    if (idx === undefined) {
		throw new Error("Couldn't find entity")
	    }
	    this.array[idx] = value
	}

	add_entity(entity: Entity, component: Component) {
	    const idx = this.array.push(component) - 1
	    this.map.set(entity.uuid, idx)
	}

	drop_entity(entity: EntityID) {
	    const idx = this.map.get(entity)
	    if (idx === undefined) {
		throw new Error("Couldn't find entity")
	    }
	    const last = this.array.pop()
	    if (last === undefined) {
		throw new Error("Empty list")
	    }
	    this.map.delete(entity)
	    if (this.array.length === 0) {
		return
	    }
	    this.array[idx] = last
	    this.map.set(last.eid, idx)
	}

	get_entities(): MapIterator<EntityID> {
	    return this.map.keys()
	}

	get_entity(entity: EntityID): Component {
	    const idx = this.map.get(entity)
	    if (idx === undefined) {
		throw new Error("Couldn't find entity")
	    }
	    return this.array[idx]
	}
    }

    export const GlobalRegistry: Registry = new Registry()

    export abstract class Component {
	constructor() { }
	public abstract eid: EntityID
	public abstract data: any
	public static ID: ComponentID
    }

    export class Object extends Component {
	public static ID = GlobalRegistry.add_component()
	constructor(public eid: EntityID, public data: THREE.Mesh) { super() }
    }

    export class Animation extends Component {
	public static ID = GlobalRegistry.add_component()
	constructor(public eid: EntityID, public data: (dt: number) => THREE.Vector3) { super() }
    }

    export class TileType extends Component {
	public static ID = GlobalRegistry.add_component()
	constructor(public eid: EntityID, public data: Tile) { super() }
    }

    export class Selectable extends Component {
	public static ID = GlobalRegistry.add_component()
	constructor(public eid: EntityID, public data: {}) { super() }
    }

    export class ObjectContainer extends Component {
	public static ID = GlobalRegistry.add_component()
	constructor(public eid: EntityID, public data: THREE.Group) { super() }
    }

    export class Parent extends Component {
	public static ID: ComponentID = GlobalRegistry.add_component()
	constructor(public eid: EntityID, public data: EntityID) { super() }
    }

    export function MakeComponent<T extends Component, D>(type: {new(id: EntityID, data: D): T, ID: ComponentID}, entity: EntityID, data: D) {
	const component = new type(entity, data)
	return {
	    id: type.ID,
	    component: component
	}
    }

    export function AnimateSystem(dt: number) {
	return new System([Object.ID, Animation.ID], (...components: Component[]) => {
	    const obj = components[0] as Object
	    const anim = components[1] as Animation
	    obj.data.position.copy(anim.data(dt))
	})
    }

    export function SelectSystem(camera: THREE.Camera, pointer: THREE.Vector2) {
	const raycaster = new THREE.Raycaster()
	return new System([Object.ID, Selectable.ID], (obj: Component, _: {}) => {
	    const object = obj as Object
	    raycaster.setFromCamera(pointer, camera)
	    const results = raycaster.intersectObject(object.data)
	    if (results.length > 0)
		ClickEventBus.handle(object.eid)
	})
    }

    export function ApplyAnimation(target: EntityID, anim: (dt: number) => THREE.Vector3) {
	GlobalRegistry.update_entity_component(target, { id: Animation.ID, component: new Animation(target, anim) })
    }

    export function Unparent() {

    }

    export function MakeNewTile(tile: Tile): {
	entity: Entity,
	components: {
	    id: ComponentID,
	    component: Component
	}[]
    } {
	const obj = new TileObject(tile)
	const entity = new ECS.Entity()

	return {
	    entity: entity,
	    components: [
		{
		    id: ECS.Object.ID,
		    component: new ECS.Object(entity.uuid, obj)
		},
		{
		    id: ECS.Animation.ID,
		    component: new ECS.Animation(entity.uuid, () => obj.position)
		}
	    ]
	}
    }

    export function MakeNewObject(object: any) {
	const obj = new ECS.Entity()
	return {
	    entity: obj,
	    components: [ { id: ECS.Object.ID, component: new ECS.Object(obj.uuid, object) } ]
	}
    }
}
