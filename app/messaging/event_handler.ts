import {IncomingMessage, Message, MessageType, validate_message} from "./message";
import {ArenaMessage} from "./arena_message";

export class EventHandler<TIncoming, TTransformed> {
    private listeners: Array<(data: TTransformed) => void> = []
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
            this.listeners.forEach(listener => listener(transformed_data))
        }
    }

    register(listener: (data: TTransformed) => void): number {
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
