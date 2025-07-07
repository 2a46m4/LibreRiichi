import {Application} from "../application";

export abstract class ApplicationState {
    app: Application

    protected constructor(app: Application) {
        this.app = app
    }

    abstract get_state_name(): string;

    async connect() {
        throw new Error(`Wrong state: ${this.get_state_name()}`)
    }

    async connect_room(room_name: string) {
        throw new Error(`Wrong state: ${this.get_state_name()}`)
    }

    async create_room(room_name: string) {
        throw new Error(`Wrong state: ${this.get_state_name()}`)
    }

    async quit_room() {
        throw new Error(`Wrong state: ${this.get_state_name()}`)
    }

    async list_rooms(): Promise<Array<string>> {
        throw new Error(`Wrong state: ${this.get_state_name()}`)
    }

    async get_arena_info(): Promise<Arena> {
        throw new Error(`Wrong state: ${this.get_state_name()}`)
    }

    async submit_move() {
        throw new Error(`Wrong state: ${this.get_state_name()}`)
    }
}
