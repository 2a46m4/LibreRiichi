import {Connection, websocket_address} from "./messaging/connection";
import {Router, useRouter} from "vue-router";
import {ApplicationState} from "./app_state/application_state";
import {LoginState} from "./app_state/login_state";
import {MessageState} from "./messaging/message_state";
import {EventHandler} from "./messaging/event_handler";

export class Application {
    connection: Connection | null;
    username: string
    public state: ApplicationState
    router: Router
    msg_state: MessageState
    handler: EventHandler

    constructor() {
        this.state = new LoginState(this)
        this.connection = null
        this.username = "No username set"
        this.router = useRouter()
        this.handler = new EventHandler()
        this.msg_state = new MessageState(this.handler)
    }

    public get action() {
        return this.state
    }

    public get conn() {
        return this.connection as Connection
    }

    public set conn(conn: Connection) {
        this.connection = conn
    }

    set_username(username: string) {
        this.username = username
    }
}
