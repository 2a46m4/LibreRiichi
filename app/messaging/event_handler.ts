import {EventHandlerInterface} from "./event_handler_interface";
import {Message, MessageType} from "./message";
import {ArenaMessage, ArenaMessageType} from "./arena_message";

export class EventHandler extends EventHandlerInterface {
    dispatch_message(data: Message) {
        console.log("Got event: ", data)
        switch (data.message_type) {
            case MessageType.ServerArenaEvent:
                this.handle_server_arena_event(data.data.arena_message)
                break;
            default:
                break;
        }
    }

    handle_server_arena_event(data: ArenaMessage) {
        console.log("Handling event:", data)
        switch (data.message_type) {
            case ArenaMessageType.PlayerJoinedEvent:
            case ArenaMessageType.PlayerQuitEvent:
            case ArenaMessageType.GameStartedEvent:
            case ArenaMessageType.ArenaBoardEvent:
            case ArenaMessageType.ListPlayersResponse:
                break
        }
    }
}