import {ApplicationState} from "./application_state";
import {Application} from "../application";
import {ArenaMessage} from "../messaging/arena_message";

export class InGameState extends ApplicationState {
    constructor(app: Application) {
        super(app);
    }

    async submit_move() {
        // TODO: Finish
    }

    async quit_room() {
        // TODO: Finish
    }

    get_state_name(): string {
        return "in game";
    }

    async handle_arena_event(msg: ArenaMessage): Promise<void> {
        console.log("Handling arena event: ", msg)
    }

}