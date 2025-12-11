import * as THREE from 'three'

export type Interpolator = (t: number) => number

export function linear_interpolator(t: number): number {
	return t
}

export function quadratic_interpolator(t: number): number {
	return t * t
}

export class TileLinearAnimation {
	t: number = 0

	constructor(
		private start: THREE.Vector3,
		private end: THREE.Vector3,
		private interp: Interpolator = linear_interpolator,
		public time: number = 1.0,
		public time_before_start: number = 0,
	) { }

	next_step(dt: number): THREE.Vector3 {
		if (this.time_before_start >= 0) {
			this.time_before_start -= dt
			return this.start
		}

		if (this.t >= 1.0) {
			return this.end
		}
		const ret = this.start.lerp(this.end, this.interp(this.t))
		this.t += dt / this.time
		return ret
	}
}

function fact(n: number) {
	let v = 1
	for (let i = 2; i <= n; i++) {
		v *= i
	}
	return n
}

function binom(n: number, k: number) {
	return fact(n) / (fact(k) * fact(n - k))
}

export class BezierAnimation {
	n: number
	t: number = 0
	is_finished: boolean = false
	binom: number[] = []

	constructor(public points: THREE.Vector3[]) {
		this.n = points.length - 1
		for (let i = 0; i <= this.n; i++) {
			this.binom.push(binom(this.n, i))
		}
	}

	next_step(dt: number): THREE.Vector3 {
		let val = new THREE.Vector3()
		for (let i = 0; i <= this.n; i++) {
			const v = Math.pow(1 - this.t, this.n - i) * Math.pow(this.t, i)
			val.add(this.points[i].multiplyScalar(this.binom[i] * v))
		}

		if (this.t >= 1.0) {
			this.is_finished = true
			this.t = 1.0
		}
		this.t += dt
		return val
	}

	finished(): boolean {
		return this.is_finished
	}
}

