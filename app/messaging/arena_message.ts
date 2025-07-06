export enum ArenaMessageType {
    // Messages that are sent from game (server) to player (client)
    PlayerJoinedEvent,
    PlayerQuitEvent,
    GameStartedEvent,
    ArenaBoardEvent,
    ListPlayersResponse,

    // Messages that are sent from player (client) to game (server)
    StartGameAction,
    PlayerAction,
    PlayerQuitAction,
    ListPlayersAction,
}

export type ArenaMessage = {
    message_type: ArenaMessage
    data: ArenaMessageData
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
    [ArenaMessageType.ListPlayersResponse]: {
        success: boolean,
        player_list: string[]
    }
    [ArenaMessageType.StartGameAction]: {}
    [ArenaMessageType.PlayerAction]: {
        // TODO
    }
    [ArenaMessageType.PlayerQuitAction]: {}
    [ArenaMessageType.ListPlayersAction]: {}

}

type ConstrainedMap<M extends Record<ArenaMessageType, any>> = {
    [K in keyof M & ArenaMessageType]: MessageEntry<K, M[K]>
}

export type ArenaMessageData = ConstrainedMap<MessageMap>[keyof MessageMap]
