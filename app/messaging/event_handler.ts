import {Message, MessageType} from "./message";
import {ArenaMessage} from "./arena_message";

export type ArenaListener = (data: ArenaMessage) => void
export type MessageListener = (data: Message) => void

export class EventHandler {

    arena_message_listeners: ArenaListener[]
    server_message_listeners: MessageListener[]

    constructor() {
        this.arena_message_listeners = []
        this.server_message_listeners = []
    }

    handle_server_message(data: Message): void {
        console.log("Got event: ", data)
        for (let i = 0; i < this.server_message_listeners.length; i++) {
            this.server_message_listeners[i](data);
        }

        if (data.message_type === MessageType.ServerArenaEvent) {
            this.handle_arena_message(data.data.arena_message)
        }
    }

    handle_arena_message(data: ArenaMessage): void {
        console.log("Handling event:", data)
        for  (let i = 0; i < this.arena_message_listeners.length; i++) {
            this.arena_message_listeners[i](data);
        }
    }

    register_arena_listener(listenerFn: ArenaListener): number {
        this.arena_message_listeners.push(listenerFn);
        return this.arena_message_listeners.length - 1
    }

    unregister_arena_listener(index: number): void {
        this.arena_message_listeners.splice(index, 1)
    }

    register_server_listener(listenerFn: MessageListener): void {
        this.server_message_listeners.push(listenerFn);
    }

    unregister_server_listener(index: number): void {
        this.server_message_listeners.splice(index, 1)
    }
}