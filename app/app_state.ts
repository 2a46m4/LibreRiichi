import {Connection, websocket_address} from "./connection";
import {Router, useRouter} from "vue-router";
import {IncomingMessage, Message, MessageType} from "./message";

export enum ApplicationState {
    NOT_CONNECTED,
    CONNECTING,
    CONNECTED,
    JOINING_ROOM,
    JOINED_ROOM,
    CREATING_ROOM,
    IN_GAME,
}

export enum GameState {
    OUT_OF_TURN,
    IN_TURN
}

type MessageResolver = (v: Message) => void

export class Application {
    connection: Connection;
    username: string
    state: ApplicationState
    router: Router
    outgoing_messages: Map<number, {
        promise: Promise<Message>,
        resolve: MessageResolver,
    }>

    constructor() {
        this.state = ApplicationState.NOT_CONNECTED
        this.router = useRouter();
        this.outgoing_messages = new Map();
    }

    private check_state(expected: ApplicationState) {
        if (this.state !== expected) {
            throw new Error("Expected state to be " + expected + ", but got " + this.state);
        }
    }

    set_username(username: string) {
        this.username = username
    }

    async connect() {
        console.log("connecting...")
        if (this.username === null) {
            // TODO: Give the user some message
            throw new Error("Username is required");
        }

        this.check_state(ApplicationState.NOT_CONNECTED)

        this.connection = new Connection(
            new WebSocket(websocket_address),
            (ev: MessageEvent<any>) => this.handle_message_recv(ev)
        )
        this.state = ApplicationState.CONNECTING
        await this.connection.wait_until_ready()

        const promiseMsg = this.send_message(
            {
                message_type: MessageType.InitialMessageAction,
                data: {name: this.username}
            }
        )

        await promiseMsg
        await this.router.push({name: 'connected_page'})
        this.state = ApplicationState.CONNECTED
        console.log("Connected!")
    }

    async connect_room(room_name: string) {
        let msg = await this.send_message(
            {
                message_type: MessageType.JoinArenaAction,
                data: {
                    arena_name: room_name,
                }
            }
        )

        if (msg === undefined) {
            // TODO: Give a reason for why
            throw new Error("Connection error: Failed to join room")
        }

        if (msg.message_type !== MessageType.GenericResponse) {
            throw new Error("Connection error: Wrong type")
        }

        if (msg.data.success === false) {
            throw new Error("Could not join room: " + msg.data.fail_reason)
        }

        console.log("Joined room")
        this.state = ApplicationState.JOINED_ROOM
    }

    async create_room(room_name: string) {
        let msg = await this.send_message({
            message_type: MessageType.CreateArenaAction,
            data: {
                arena_name: room_name,
            }
        })

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

    quit_room() {
        // TODO: Finish
    }

    async list_rooms() : Promise<Array<string>> {
        let msg = await this.send_message({
            message_type: MessageType.ListArenasAction,
            data: {}
        })

        if (msg.message_type !== MessageType.ListArenasResponse) {
            throw new Error("Wrong type")
        }

        if (msg.data.success === false) {
            throw new Error("Failed to list rooms")
        }

        return msg.data.arena_list
    }

    submit_move() {
        // TODO: Finish
    }

    send_message(msg: Message) {
        console.log("Sending: ", msg)
        let resolveMsg: MessageResolver = (_: Message) => {
        };
        let promiseMsg = new Promise((resolve: MessageResolver) => {
            resolveMsg = resolve
        });
        let msg_index = this.connection.send(msg)
        this.outgoing_messages.set(msg_index, {
            promise: promiseMsg,
            resolve: resolveMsg,
        })

        return promiseMsg
    }

    handle_message_recv(ev: MessageEvent) {
        let data = JSON.parse(ev.data)
        console.log("Got message: ", data)
        this.match_outgoing_message(data)
    }

    // If it matches an outgoing message, it's a response to the outgoing message. Otherwise, it's a fresh event from the server.
    match_outgoing_message(data: IncomingMessage) {
        console.log("Got return: ", data)

        if (this.outgoing_messages.has(data.message_index)) {
            console.log("Matched outgoing message")
            this.outgoing_messages.get(data.message_index).resolve(data)
            this.outgoing_messages.delete(data.message_index)
            return
        }

        this.handle_event(data)
    }

    handle_event(data: Message) {
        console.log("Got event: ", data)
    }

}
