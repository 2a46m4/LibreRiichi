import {Connection, websocket_address} from "./connection";
import {useRouter, Router} from "vue-router";
import {Message} from "./message";

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

export class Application {
	connection: Connection;
	username: string
	state: ApplicationState
	router: Router
	outgoing_messages: Message[]

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
			(ev) => this.handle_message_recv(ev)
		)
		this.state = ApplicationState.CONNECTING
		await this.connection.WaitUntilReady()
		await this.router.push({name: 'connected_page'})
		this.state = ApplicationState.CONNECTED
		console.log("Connected!")
	}

	connect_room(room_name: string) {
		// TODO: Finish
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

	handle_message_recv(ev: MessageEvent) {
		try {
			let data = JSON.parse(ev.data)
			this.match_outgoing_message(data)
		} catch (e) {
			console.log(e)
		}
		
	}

	// If it matches an outgoing message, it's a response to the outgoing message. Otherwise, it's a fresh event from the server.
	match_outgoing_message(data: any) {

		console.log(data)
	}
}
