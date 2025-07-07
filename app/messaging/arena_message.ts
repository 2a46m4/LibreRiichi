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

type MessageEntry<T extends ArenaMessageType = ArenaMessageType, D = any> = {
    message_type: T
    data: D
}

type MessageMap = {
    [ArenaMessageType.PlayerJoinedEvent]: {
        name: string,
        id: string
    }
    [ArenaMessageType.PlayerQuitEvent]: {
        name: string,
    }
    [ArenaMessageType.GameStartedEvent]: {}
    [ArenaMessageType.ArenaBoardEvent]: {
        // TODO
    }
    [ArenaMessageType.StartGameAction]: {}
    [ArenaMessageType.PlayerAction]: {
        // TODO
    }
    [ArenaMessageType.PlayerQuitAction]: {}

}

type ConstrainedMap<M extends Record<ArenaMessageType, any>> = {
    [K in keyof M & ArenaMessageType]: MessageEntry<K, M[K]>
}

export type ArenaMessage = ConstrainedMap<MessageMap>[keyof MessageMap]
