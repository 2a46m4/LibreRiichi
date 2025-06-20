import {z} from "zod/v4"

export enum MessageType {
    // Messages that are sent from server to client
    InitialMessageEvent,
    JoinArenaEvent,
    ServerArenaEvent,
    ListArenasEvent,
    CreateArenaEvent,

// Messages that are sent from client to server
    InitialMessageAction,
    JoinArenaAction,
    ServerArenaAction,
    ListArenasAction,
    CreateArenaAction,
}

export enum ArenaMessageType {
    // Messages that are sent from game (server) to player (client)
    PlayerJoinedEvent,
    PlayerQuitEvent,
    GameStartedEvent,
    ArenaBoardEvent,

    // Messages that are sent from player (client) to game (server)
    StartGameAction,
    PlayerAction,
    PlayerQuitAction,

}

export type Message = {
    message_type: MessageType
    data: any
}

export class InitialMessageEvent extends Message {
    readonly message_type = MessageType.InitialMessageEvent
    readonly data: null
}

export class InitialMessageAction extends Message {
    readonly message_type = MessageType.InitialMessageAction
    data = {name: ''}

    constructor(name: string) {
        super()
        this.data = {name: name}
    }
}

export class ServerArenaEvent extends Message {
    readonly message_type = MessageType.ServerArenaEvent
    data: { arena_message: ArenaMessage }
}

export class ServerArenaAction extends Message {
    readonly message_type = MessageType.ServerArenaAction
    data: {
        arena_message: ArenaMessage
    }
}

export class JoinArenaActionData extends Message {
    readonly message_type = MessageType.JoinArenaAction
    data: {
        arena_name: string
        arena_id: Uint8Array
    }

    constructor(arena_name: string, arena_id: Uint8Array) {
        super();
        this.data = {arena_name: arena_name, arena_id: arena_id}
    }
}

export class JoinArenaEventData extends Message {
    readonly message_type = MessageType.JoinArenaEvent
    data: {
        success: boolean
    }
}

export type ArenaMessage = {
    message_type: ArenaMessage
    data: {
        message_type: ArenaMessageType
        data: any
    }
}