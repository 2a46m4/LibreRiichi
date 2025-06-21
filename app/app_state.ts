import {Connection, websocket_address} from "./connection";
import {Router, useRouter} from "vue-router";
import {getMatchingMessageType, Message, MessageType} from "./message";

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

class OutgoingMessageState {
    return: Promise<Message | undefined>
    resolve: MessageResolver
    outgoing: Message
}

export class Application {
    connection: Connection;
    username: string
    state: ApplicationState
    router: Router
    outgoing_messages: {
        promise: Promise<Message>,
        resolve: MessageResolver,
        outgoing: Message
    }[]

    constructor() {
        this.state = ApplicationState.NOT_CONNECTED
        this.router = useRouter();
        this.outgoing_messages = [];
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
            throw new Error("Username is required");
        }

        this.check_state(ApplicationState.NOT_CONNECTED)

        this.connection = new Connection(
            this.username,
            new WebSocket(websocket_address),
            (ev: MessageEvent<any>) => this.handle_message_recv(ev)
        )
        this.state = ApplicationState.CONNECTING

        const promiseMsg = this.push_outgoing_message(
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
        this.connection.Send({
            message_type: MessageType.JoinArenaAction,
            data: {
                arena_name: room_name,
            }
        })

        const promiseMsg = this.push_outgoing_message(
            {
                message_type: MessageType.JoinArenaAction,
                data: {
                    arena_name: room_name,
                }
            }
        )

        let msg = await promiseMsg

        if (msg === undefined) {
            // TODO: Give a reason for why
            throw new Error("Connection error: Failed to join room")
        }

        if (msg.message_type !== MessageType.JoinArenaEvent) {
            throw new Error("Connection error: Wrong type")
        }

        if (msg.data.success === false) {
            throw new Error("Could not join room")
        }

    }

    quit_room() {
        // TODO: Finish
    }

    list_rooms() {
        // TODO: Finish
    }

    submit_move() {
        // TODO: Finish
    }

	push_outgoing_message(msg: Message) {
		let resolveMsg: MessageResolver = (_: Message) => {};
		let promiseMsg = new Promise((resolve: MessageResolver) => {
			resolveMsg = resolve
		});

		this.outgoing_messages.push({
			promise: promiseMsg,
			resolve: resolveMsg,
			outgoing: msg
		})

		return promiseMsg
	}

    handle_message_recv(ev: MessageEvent) {
        try {
            let data = JSON.parse(ev.data)
            this.match_outgoing_message(data)
        } catch (e) {
            console.log(e)
        }

    }

    // If it matches an outgoing message, it's a response to the outgoing message. Otherwise, it's a fresh event from the server.
    match_outgoing_message(data: Message) {
            console.log("Got return: ", data)

	    const matching = getMatchingMessageType(data.message_type)
	    for (let i = 0; i < this.outgoing_messages.length; i++) {
		    if (this.outgoing_messages[i].outgoing.message_type === matching) {
			    this.outgoing_messages.splice(i)
		    }
	    }
            this.outgoing_messages[0].resolve(data)
    }
}
