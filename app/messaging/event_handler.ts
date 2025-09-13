import {IncomingMessage, MessageType, validate_message} from "./message";
import {ServerResponseMessage} from "./server_response_generated";
import {data} from "autoprefixer";

export class EventHandler<TIncoming, TTransformed> {
    private listeners: Array<(data: TTransformed) => boolean> = []
    private readonly transform: (data: TIncoming) => TTransformed
    private readonly conditional: (data: TTransformed) => boolean

    constructor(
        transform: (data: TIncoming) => TTransformed,
        conditional: (data: TTransformed) => boolean
    ) {
        this.transform = transform
        this.conditional = conditional
    }

    handle(data: TIncoming): void {
        const transformed_data = this.transform(data)
        if (this.conditional(transformed_data)) {
            this.listeners.filter(listener => listener(transformed_data))
        }
    }

    register(listener: (data: TTransformed) => boolean): number {
        this.listeners.push(listener)
        return this.listeners.length - 1
    }

    unregister(index: number): void {
        this.listeners.splice(index, 1)
    }
}

export function create_event_handler<TIncoming, TTransformed>(
    transform: (data: TIncoming) => TTransformed,
    conditional: (data: TTransformed) => boolean = () => true
) {
    return new EventHandler(transform, conditional)
}

export const ServerMessageBus = create_event_handler(
    (data: MessageEvent) => JSON.parse(data.data) as IncomingMessage,
    (msg: IncomingMessage) => validate_message(msg).isValid
)

export const ArenaMessageBus = create_event_handler(
    (data: IncomingMessage) => data,
    (msg: IncomingMessage)=> msg.message_type === MessageType.EVENT
)

export function register_request(msg_idx: number) : Promise<ServerResponseMessage> {
    let {promise, resolve} = Promise.withResolvers<ServerResponseMessage>();

    ServerMessageBus.register((msg)=> {
        if (msg.message_type === MessageType.RESPONSE && msg.message_index === msg_idx) {
            console.log("Matched outgoing message ", msg_idx, ", resolving")
            resolve(msg.data)
            return false
        } else {
            return true
        }
    })

    return promise
}