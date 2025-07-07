import {ArenaMessage} from "./arena_message";

export enum MessageType {
    // Messages that are sent from server to client in response to an event
    ServerArenaEvent,

	// Messages sent in response to an action
    GenericResponse,
    ListArenasResponse,
    ArenaInfoResponse,

	// Messages that are sent from client to server
    InitialMessageAction,
    JoinArenaAction,
    ServerArenaAction,
    ListArenasAction,
    CreateArenaAction,
    ArenaInfoAction,
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
		arena_list: string[]
	}
    [MessageType.ArenaInfoResponse]: {
        success: boolean,
        name: string,
        agents: Array<{
            name: string
        }>,
        game_started: boolean,
        date_created: string
    }
    [MessageType.ArenaInfoAction]: {}
}

type ConstrainedMap<M extends Record<MessageType, any>> = {
    [K in keyof M & MessageType]: MessageEntry<K, M[K]>
}

export type Message = ConstrainedMap<MessageMap>[keyof MessageMap]
