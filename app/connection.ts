import {InitialMessageAction, Message} from "./Message";

export class Connection {
    socket: WebSocket;

    constructor(name: string, websocket: WebSocket) {
        this.socket = websocket;
        this.socket.onopen = () => {
            this.Send(new InitialMessageAction(name))
        }
        this.socket.onmessage = this.Receive
        this.socket.onclose()
    }

    Send(msg: Message) {
        console.log(msg)
        this.socket.send(JSON.stringify(msg))
    }

    Receive(e: MessageEvent): any {
        console.log(e)
    }
}