import {IncomingMessage, Message, MessageType} from "./message";
import {EventHandler} from "./event_handler";

type MessageResolver = (v: Message) => void

export class MessageState {
    outgoing_messages: Map<number, {
        resolve: MessageResolver,
        reject: MessageResolver,
    }>
    event_handler: EventHandler

    constructor(handler: EventHandler) {
        this.outgoing_messages = new Map();
        this.event_handler = handler
    }

    handle_message_event(ev: MessageEvent) {
        let data = JSON.parse(ev.data)
        this.match_message(data)
    }

    // TODO: Timeout option
    register_message(msg_idx: number): Promise<Message> {
        let {promise, resolve, reject} = Promise.withResolvers<Message>();

        this.outgoing_messages.set(msg_idx, {
            resolve: resolve,
            reject: reject,
        })

        return promise
    }

    match_message(data: IncomingMessage) {
        console.log("Got a message: ", data)

        if (this.outgoing_messages.has(data.message_index)) {
            console.log("Matched outgoing message")
            this.outgoing_messages.get(data.message_index)?.resolve(data)
            this.outgoing_messages.delete(data.message_index)
            return
        }

        else {
            console.log("No match for message: ", data)
            console.log("Calling event handler")

            this.event_handler.handle_server_message(data)
        }
    }
}