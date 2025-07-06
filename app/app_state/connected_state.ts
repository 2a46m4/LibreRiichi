import {ApplicationState} from "./application_state";
import {Message, MessageType} from "../messaging/message";
import {JoinedRoomState} from "./joined_room_state";
import {Application} from "../application";

export class ConnectedState extends ApplicationState {
    constructor(app: Application) {
        super(app);
    }

    get_state_name(): string {
        return "connected";
    }

    async connect_room(room_name: string): Promise<void> {
        let msg_idx = this.app.conn.send(
            {
                message_type: MessageType.JoinArenaAction,
                data: {
                    arena_name: room_name,
                }
            }
        )

        let ret = await this.app.msg_state.register_message(msg_idx)
        if (ret.message_type !== MessageType.GenericResponse) {
            throw new Error("Connection error: Wrong type")
        }

        if (ret.data.success === false) {
            throw new Error("Could not join room: " + ret.data.fail_reason)
        }

        console.log("Joined room")
        this.app.state = new JoinedRoomState(this.app)
        await this.app.router.push({name: 'arena_page'})
    }

    async create_room(room_name: string) {
        let msg_idx = this.app.conn.send({
            message_type: MessageType.CreateArenaAction,
            data: {
                arena_name: room_name,
            }
        })

        let msg = await this.app.msg_state.register_message(msg_idx)

        if (msg === undefined) {
            // TODO: Give a reason for why
            throw new Error("Connection error: Failed to create room")
        }

        if (msg.message_type !== MessageType.GenericResponse) {
            throw new Error("Connection error: Wrong type")
        }

        if (msg.data.success === false) {
            throw new Error("Could not create room: " + msg.data.fail_reason)
        }
    }

    async list_rooms() : Promise<Array<string>> {
        let msg_idx = this.app.conn.send({
            message_type: MessageType.ListArenasAction,
            data: {}
        })

        let msg = await this.app.msg_state.register_message(msg_idx)

        if (msg.message_type !== MessageType.ListArenasResponse) {
            throw new Error("Wrong type")
        }

        if (msg.data.success === false) {
            throw new Error("Failed to list rooms")
        }

        return msg.data.arena_list.sort()
    }
}