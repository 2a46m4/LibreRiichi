import {EventHandler} from "./event_handler";
import {Message, MessageType} from "./message";
import {ArenaMessage} from "./arena_message";
export type ArenaListener = (data: ArenaMessage) => void

export default class ArenaHandler {
    arena_message_listeners: ArenaListener[]

    constructor(event_handler: EventHandler) {
        event_handler.register_server_listener(this.handle_arena_message.bind(this))
        this.arena_message_listeners = []
    }

    handle_arena_message(msg: Message) {
        if (msg.message_type === MessageType.ServerArenaEvent) {
            console.log("Handling event:", msg)
            for (let i = 0; i < this.arena_message_listeners.length; i++) {
                this.arena_message_listeners[i](msg.data.arena_message);
            }
        }
    }

    register_arena_listener(listenerFn: ArenaListener): number {
        this.arena_message_listeners.push(listenerFn);
        return this.arena_message_listeners.length - 1
    }
    unregister_arena_listener(index: number): void {
        this.arena_message_listeners.splice(index, 1)
    }
}