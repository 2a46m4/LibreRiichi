import {IncomingMessage, Message, MessageType} from "./message";
import {ArenaMessage} from "./arena_message";

export type MessageListener = (data: IncomingMessage) => void

export class EventHandler<
    Incoming,
    Transformed,
    Transform extends (data: Incoming) => Transformed,
    Conditional extends (data: Transformed) => boolean,
    Listener extends (data: Transformed) => void,
> {

    listeners: Listener[]
    transform: Transform
    conditional: Conditional

    constructor(transform: Transform, conditional: Conditional) {
        this.listeners = []
        this.transform = transform
        this.conditional = conditional
    }

    handle(data: Incoming): void {
        let transformed_data = this.transform(data)
        if (this.conditional(transformed_data)) {
            this.listeners.forEach(l=>l(transformed_data))
        }
    }

    register(listenerFn: Listener): number {
        this.listeners.push(listenerFn);
        return this.listeners.length - 1
    }

    unregister(index: number): void {
        this.listeners.splice(index, 1)
    }
}


//
// export class EventHandler {
//
//     server_message_listeners: MessageListener[]
//
//     constructor() {
//         this.server_message_listeners = []
//     }
//
//     handle_server_message(event: MessageEvent): void {
//         console.log("Got event: ", event)
//         let data: IncomingMessage = JSON.parse(event.data);
//         for (let i = 0; i < this.server_message_listeners.length; i++) {
//             this.server_message_listeners[i](data);
//         }
//     }
//
//     register_server_listener(listenerFn: MessageListener): number {
//         this.server_message_listeners.push(listenerFn);
//         return this.server_message_listeners.length - 1
//     }
//
//     unregister_server_listener(index: number): void {
//         this.server_message_listeners.splice(index, 1)
//     }
// }