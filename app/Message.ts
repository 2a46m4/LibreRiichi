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

export type Message = {
    message_type: MessageType
    data: any
}

export const InitialMessageEvent: Message = {
    message_type: MessageType.InitialMessageEvent,
    data: null
}

export class InitialMessageAction {
    message_type = MessageType.InitialMessageAction
    data : { name: string } = { name: '' }
}

export type ServerArenaEvent = {
    message_type: MessageType.ServerArenaEvent,
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