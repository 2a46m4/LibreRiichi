import {ApplicationState} from "./application_state";
import {Connection, websocket_address} from "../messaging/connection";
import {MessageType} from "../messaging/message";
import {ConnectedState} from "./connected_state";
import {Application} from "../application";

export class LoginState extends ApplicationState {
    constructor(app: Application) {
        super(app);
    }

    get_state_name(): string {
        return "login";
    }

    async connect() {
        if (this.app.username === null) {
            throw new Error("Username is required")
        }

        this.app.conn = new Connection(
            new WebSocket(websocket_address),
            (ev: MessageEvent<any>) => {
                this.app.msg_state.handle_message_recv(ev)
            }
        )

        await this.app.conn.wait_until_ready()

        let msg_idx = this.app.conn.send({
            message_type: MessageType.InitialMessageAction,
            data: {name: this.app.username}
        })

        let ret = await this.app.msg_state.register_message(msg_idx)
        if (ret.message_type === MessageType.GenericResponse) {
            if (!ret.data.success) {
                console.log("Failed to connect: ", ret.data.fail_reason)
            }

            this.app.state = new ConnectedState(this.app)
            await this.app.router.push({name: 'connected_page'})
        } else {
            console.log("Error message type")
        }
    }
}