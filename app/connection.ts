import {Message, MessageType} from "./message";

export const websocket_address = "ws://localhost:3000/game";

export class Connection {
	socket: WebSocket;
	ready: Promise<void>;
	handler: (e: MessageEvent) => any

	constructor(name: string, websocket: WebSocket, handler: (e: MessageEvent) => any) {
		this.socket = websocket;
		this.handler = handler

		this.ready = new Promise((resolve) => {
			this.socket.onopen = (ev) => {
				console.log("Connection opened");
				this.Send({
					message_type: MessageType.InitialMessageAction,
					data: {name: name}
				})
				resolve();
			}
		})
		this.socket.onmessage = handler
		this.socket.onclose = () => {}
	}

	async WaitUntilReady(): Promise<void> {
		await this.ready
	}

	Send(msg: Message) {
		console.log(msg)
		this.socket.send(JSON.stringify(msg))
	}

}
