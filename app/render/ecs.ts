import { v4 as uuidv4 } from 'uuid'

namespace ECS {
	type EntityID = string & { readonly __brand: unique symbol }
	type ComponentID = string & { readonly __brand: unique symbol }
	type SystemID = string & { readonly __brand: unique symbol }

	export class Registry {
		entities: Map<EntityID, {
			components: ComponentID[],
		}> = new Map()
		components: Map<ComponentID, ComponentArray> = new Map()
		systems: Map<SystemID, System> = new Map()

		add_entity(entity: Entity, components: { id: ComponentID, component: Component }[]) {
			this.entities.set(entity.uuid, { components: components.map(c => c.id) })
			for (let component of components) {
				const c = this.components.get(component.id)
				c?.add_entity(entity, component.component)
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
			this.components.forEach(c=>c.drop_entity(entity))
		}

		add_component(): ComponentID {
			const arr = new ComponentArray()
			this.components.set(arr.component_uuid, arr)
			return arr.component_uuid
		}

		add_system(system: System) {
			this.systems.set(system.uuid, system)
		}

		run_system(system: System | SystemID): void {
			let id: SystemID
			if (typeof system === "object") {
				id = system.uuid
			} else {
				id = system
			}

			const sys = this.systems.get(id)
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
				.forEach(val => sys.apply(val))
		}

		run_all_systems() {
			this.systems.values().forEach(s => this.run_system(s))
		}
	}

	export class System {
		public uuid: SystemID = uuidv4() as SystemID
		public component_ids: ComponentID[] = []
		constructor(public apply: (components: Component[]) => void) { }
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

		add_entity(entity: Entity, value: any) {
			const idx = this.array.push(new Component(entity.uuid, value)) - 1
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
			this.map.set(last.entity_uuid, idx)
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

	export class Component {
		constructor(public entity_uuid: EntityID, public data: any) { }
	}
}