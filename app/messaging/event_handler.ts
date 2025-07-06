import {EventHandlerInterface} from "./event_handler_interface";
import {Message, MessageType} from "./message";
import {ArenaMessage} from "./arena_message";

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
            // Messages that are sent from game (server) to player (client)
            PlayerJoinedEvent:
                PlayerQuitEvent:
                    GameStartedEvent:
                    ArenaBoardEvent:
                    ListPlayersResponse:
        }
    }
}