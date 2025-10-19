import * as THREE from 'three'
import { TileObject } from './tile'

export interface IAnimation {
    next_step(dt: number): void
    finished(): boolean
    time_before_start: number
}

export interface IAnimationManager {
    add_animation(object: IAnimation): void
    animate_step(dt: number): void
    get_animations(): IAnimation[]
}

export type Interpolator = (t: number) => number

export function linear_interpolator(t: number): number {
    return t
}

export function quadratic_interpolator(t: number): number {
    return t * t
}

export class MultipleTileAnimation implements IAnimation {
    constructor(
        public animations: TileAnimation[],
        public time_before_start: number = 0
    ) { }

    add(animation: TileAnimation) {
        this.animations.push(animation)
    }

    next_step(dt: number): void {
        console.log(this.animations, this.current_index)
        if (this.finished()) {
            return
        }
        if (this.animations[this.current_index].finished()) {
            this.current_index += 1
        }

        if (this.finished()) {
            return
        }

        this.animations[this.current_index].next_step(dt)
    }
    finished(): boolean {
        return this.current_index === this.animations.length
    }
    current_index = 0
}

export class TileAnimation implements IAnimation {
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
    animations_in_flight: IAnimation[] = []

    constructor() { }

    get_animations(): IAnimation[] {
        return this.animations_in_flight
    }

    add_animation(obj: IAnimation): void {
        this.animations_in_flight.push(obj)
    }

    animate_step(dt: number): void {
        for (let animation of this.animations_in_flight) {
            animation.next_step(dt)
        }

        this.animations_in_flight = this.animations_in_flight.filter(a => !a.finished())
    }
}


