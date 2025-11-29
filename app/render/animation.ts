import * as THREE from 'three'
import { TileObject } from './tile'

// An object that exists in the scene and can be animated
export interface IAnimatable {
    animate(dt: number): void
    add_animation(animation: IAnimation): void
}

// A type of animation, to be called by an IAnimatable object
export interface IAnimation {
    next_step(dt: number): void
    finished(): boolean
}

export interface IAnimationManager {
    add_animation(object: IAnimatable): void
    remove_animation(object: IAnimatable): void
    animate_step(dt: number): void
    get_animations(): IAnimatable[]
}

export type Interpolator = (t: number) => number

export function linear_interpolator(t: number): number {
    return t
}

export function quadratic_interpolator(t: number): number {
    return t * t
}

export class TileLinearAnimation implements IAnimation {
    is_finished: boolean = false
    t: number = 0

    constructor(
        private tile: TileObject,
        private start: THREE.Vector3,
        private end: THREE.Vector3,
        private interp: Interpolator = linear_interpolator,
        public time: number = 1.0,
        public time_before_start: number = 0,
    ) { }

    next_step(dt: number): void {

        if (this.time_before_start >= 0) {
            this.time_before_start -= dt
            return
        }

        if (this.t >= 1.0) {
            this.is_finished = true
            this.t = 1.0
        }

        this.tile.position.lerpVectors(
            this.start, this.end, this.interp(this.t)
        )

        this.t += dt / this.time
    }

    finished(): boolean {
        return this.is_finished
    }
}

export class AnimationManager implements IAnimationManager {
    animations_in_flight: IAnimatable[] = []

    constructor() { }

    remove_animation(object: IAnimatable): void {
        const index = this.animations_in_flight.indexOf(object)
        if (index !== -1) {
            this.animations_in_flight.splice(index, 1)
        }
    }

    get_animations(): IAnimatable[] {
        return this.animations_in_flight
    }

    add_animation(obj: IAnimatable): void {
        this.animations_in_flight.push(obj)
    }

    animate_step(dt: number): void {
        for (let animation of this.animations_in_flight) {
            animation.animate(dt)
        }
    }
}
