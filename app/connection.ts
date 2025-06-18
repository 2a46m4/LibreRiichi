import {InitialMessageAction, Message} from "./message";

export const websocket_address = "ws://localhost:3000/game";

export class Connection {
    socket: WebSocket;
    ready: Promise<void>;

    constructor(name: string, websocket: WebSocket) {
        this.socket = websocket;
        this.ready = new Promise((resolve) => {
            this.socket.onopen = (ev) => {
                console.log("Connection opened");
                this.Send(new InitialMessageAction(name))
                resolve();
            }
        })
        this.socket.onmessage = this.Receive
        this.socket.onclose = () => {}
    }

    async WaitUntilReady(): Promise<void> {
        await this.ready
    }

    Send(msg: Message) {
        console.log(msg)
        this.socket.send(JSON.stringify(msg))
    }

    Receive(e: MessageEvent): any {
        console.log(e)
    }
}