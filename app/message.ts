export enum MessageType {
    // Messages that are sent from server to client in response to an event
    ServerArenaEvent,

	// Messages sent in response to an action
    GenericResponse,
    ListArenasResponse,

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

export type IncomingMessage = Message & {
    message_index: number
}

type MessageEntry<T extends MessageType = MessageType, D = any> = {
    message_type: T
    data: D
}

type MessageMap = {
    [MessageType.InitialMessageAction]: { name: string }
    [MessageType.ServerArenaEvent]: { arena_message: ArenaMessage }
    [MessageType.ServerArenaAction]: {
        arena_message: ArenaMessage
    }
    [MessageType.JoinArenaAction]: {
        arena_name: string
    }
    [MessageType.CreateArenaAction]: {
        arena_name: string
    }
    [MessageType.ListArenasAction]: {}
	[MessageType.GenericResponse]: {success: boolean, fail_reason: string}
	[MessageType.ListArenasResponse]: {
		success: boolean,
		arena_names: string[]
	}
}

type ConstrainedMap<M extends Record<MessageType, any>> = {
    [K in keyof M & MessageType]: MessageEntry<K, M[K]>
}

export type Message = ConstrainedMap<MessageMap>[keyof MessageMap]

export type ArenaMessage = {
    message_type: ArenaMessage
    data: {
        message_type: ArenaMessageType
        data: any
    }
}

