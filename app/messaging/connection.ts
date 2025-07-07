import {Message, MessageType} from "./message";

export const websocket_address = "ws://localhost:3000/game";

export class Connection {
	socket: WebSocket;
	ready_promise: Promise<void>;
	ready: boolean;
	index: number
	handler: (e: MessageEvent) => any

	constructor(websocket: WebSocket, handler: (e: MessageEvent) => any) {
		this.socket = websocket;
		this.handler = handler
		this.ready = false
		this.index = 0
		this.ready_promise = new Promise((resolve) => {
			this.socket.onopen = (ev) => {
				this.ready = true;
				resolve();
			}
		})

		this.socket.onmessage = handler
		this.socket.onclose = () => {}
	}

	async wait_until_ready(): Promise<void> {
		await this.ready_promise
	}

	send(msg: Message): number {
		if (!this.ready) {
			throw new Error("Not ready")
		}
		this.socket.send(JSON.stringify({...msg, message_index: this.index}))
		this.index += 1
		return this.index - 1
	}

}
