import * as THREE from 'three'
import { TileObject } from './tile'

export interface IAnimation {
    next_step(dt: number): void
    finished(): boolean
}

export interface IAnimationManager {
    add_animation(object: IAnimation): void
    animate_step(dt: number): void
}

export type Interpolator = (t: number) => number

export function linear_interpolator(t: number): number {
    return t
}

export function quadratic_interpolator(t: number): number {
    return t * t
}

export class TileAnimation implements IAnimation {
    is_finished: boolean = false
	t: number = 0

    constructor(
        private tile: TileObject,
        private start: THREE.Vector3,
        private end: THREE.Vector3,
        private interp: Interpolator = linear_interpolator
    ) { }

    next_step(dt: number): void {
        if (this.is_finished)
            throw new Error("Animation already finished")
		
		this.tile.position.lerpVectors(
			this.start, this.end, this.interp(this.t)
		)
		
		if (this.t >= 1.0) {
			this.is_finished = true
		} else {
			this.t += dt
		}
    }

    finished(): boolean {
		return this.is_finished
    }
}

export class AnimationManager implements IAnimationManager {
    animations_in_flight: IAnimation[] = []

    constructor() { }

    add_animation(obj: IAnimation): void {
		this.animations_in_flight.push(obj)
    }

    animate_step(dt: number): void {
        for (let animation of this.animations_in_flight) {
			animation.next_step(dt)
		}
			
		this.animations_in_flight.filter(a=>!a.finished())
    }
}


