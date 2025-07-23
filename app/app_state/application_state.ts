import {Application} from "../application";
import {Arena} from "../types/arena";
import {ArenaMessage} from "../messaging/arena_message";

export abstract class ApplicationState {
    app: Application

    protected constructor(app: Application) {
        this.app = app
    }

    abstract get_state_name(): string;

    transition(new_state: ApplicationState): void {
        this.app.state = new_state
    }

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

    async start_game() {
        throw new Error(`Wrong state: ${this.get_state_name()}`)
    }

    async handle_arena_event(msg: ArenaMessage) {
        throw new Error(`Wrong state: ${this.get_state_name()}`)
    }
}
