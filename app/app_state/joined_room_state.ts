import {ApplicationState} from "./application_state";
import {Application} from "../application";
import {MessageType} from "../messaging/message";
import {Arena} from "../types/arena";
import {ArenaMessageType} from "../messaging/arena_message";

export class JoinedRoomState extends ApplicationState {
    constructor(app: Application) {
        super(app);
    }

    async start_game() {
        let msg_idx = this.app.conn.send(
            {
                message_type: MessageType.ServerArenaAction,
                data: {
                    arena_message: {
                        message_type: ArenaMessageType.StartGameAction,
                        data: {}
                    }
                }
            }
        )

        let ret = await this.app.msg_state.register_message(msg_idx);
        if (ret.message_type !== MessageType.GenericResponse) {
            throw new Error("Connection error: Wrong Type")
        }

        if (!ret.data.success) {
            throw new Error("Couldn't start game: " + ret.data.fail_reason)
        }
    }

    async get_arena_info(): Promise<Arena> {
        let msg_idx = this.app.conn.send(
            {
                message_type: MessageType.ArenaInfoAction,
                data: {}
            }
        )

        let ret = await this.app.msg_state.register_message(msg_idx)
        if (ret.message_type !== MessageType.ArenaInfoResponse) {
            throw new Error("Connection error: wrong type")
        }

        if (!ret.data.success) {
            throw new Error("Could not get arena data")
        }

        return ret.data
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