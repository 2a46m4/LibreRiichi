import {IncomingMessage, MessageType} from "./message";
import {ServerResponseMessage} from "./server_response_generated";

type MessageResolver = (v: any) => void

export class MessageRouter {
    outgoing_messages: Map<number, {
        resolve: MessageResolver,
        reject: MessageResolver,
    }>

    constructor() {
        this.outgoing_messages = new Map();
    }

    // TODO: Timeout option
    register_message(msg_idx: number): Promise<ServerResponseMessage> {
        let {promise, resolve, reject} = Promise.withResolvers<ServerResponseMessage>();

        this.outgoing_messages.set(msg_idx, {
            resolve: resolve,
            reject: reject,
        })

        return promise
    }

    match_message(data: IncomingMessage) {
        if (data.message_type === MessageType.RESPONSE && this.outgoing_messages.has(data.message_index)) {
            console.log("Matched outgoing message, resolving")
            this.outgoing_messages.get(data.message_index)?.resolve(data.data)
            this.outgoing_messages.delete(data.message_index)
        }
    }
}