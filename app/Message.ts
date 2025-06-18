export enum MessageType {
    InitialMessageEvent,
    JoinArenaEvent,
    ServerArenaEvent,
    InitialMessageAction,
    JoinArenaAction,
    ServerArenaAction,
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
    data = { name: '' }

    constructor(name: string) {
        super()
        this.data = { name: name }
    }
}

export class ServerArenaEvent {
    readonly message_type = MessageType.ServerArenaEvent
    data: { arena_message: ArenaMessage }
}

export type ServerArenaAction = {
    message_type: MessageType.ServerArenaAction,
    data: {
        arena_message: ArenaMessage
    }
}

export type JoinArenaActionData = {
    message_type: MessageType.JoinArenaAction,
    data: {
        arena_name: string
    }
}

export type JoinArenaEventData = {
    message_type: MessageType.JoinArenaEvent,
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