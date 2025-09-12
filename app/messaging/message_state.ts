import {IncomingMessage, Message, MessageType} from "./message";
import {EventHandler} from "./event_handler";

type MessageResolver = (v: Message) => void

export class MessageState {
    outgoing_messages: Map<number, {
        resolve: MessageResolver,
        reject: MessageResolver,
    }>

    constructor() {
        this.outgoing_messages = new Map();
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

    match_message(data: Message) {
        if (data.message_type === MessageType.RESPONSE && this.outgoing_messages.has(data.message_index)) {
            console.log("Matched outgoing message")
            this.outgoing_messages.get(data.message_index)?.resolve(data)
            this.outgoing_messages.delete(data.message_index)
        }
    }
}