import {ApplicationState} from "./application_state";
import {Application} from "../application";
import app from "../views/app.vue";

export class JoinedRoomState extends ApplicationState {
    constructor(app: Application) {
        super(app);
    }

    async start_game() {

    }

    async submit_move() {
        // TODO: Finish
    }

    async quit_room() {
        // TODO: Finish
    }

    get_state_name(): string {
        return "joined_room";
    }

}