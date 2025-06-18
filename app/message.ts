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

export class Message {
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
    }
}

export class JoinArenaEventData extends Message {
    readonly message_type = MessageType.JoinArenaEvent
    data: {
        success: boolean
    }
}

export function checkMessage(message: any) {

}

export type ArenaMessage = {
    message_type: ArenaMessage
    data: {
        message_type: ArenaMessageType
        data: any
    }
}