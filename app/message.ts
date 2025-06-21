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

type MessageEntry<T extends MessageType = MessageType, D = any> = {
	message_type: T
	data: D
}

type MessageMap = {
	[MessageType.InitialMessageEvent]: null
	[MessageType.InitialMessageAction]: { name: string }
	[MessageType.ServerArenaEvent]: { arena_message: ArenaMessage }
	[MessageType.ServerArenaAction]: {
		arena_message: ArenaMessage
	}
	[MessageType.JoinArenaAction]: {
		arena_name: string
	}
	[MessageType.JoinArenaEvent]: {
		success: boolean
	}
	[MessageType.CreateArenaEvent]: {
		success: boolean
	}
	[MessageType.CreateArenaAction]: {
		arena_name: string
	}
	[MessageType.ListArenasAction]: {}
	[MessageType.ListArenasEvent]: {
		arena_names: string[]
	}
}

type ConstrainedMap<M extends Record<MessageType, any>> = {
	[K in keyof M & MessageType]: MessageEntry<K, M[K]>
}

export type Message = ConstrainedMap<MessageMap>[keyof MessageMap]

// export type InitialMessageEvent = {
// 	message_type: MessageType.InitialMessageEvent
// 	data: null
// }

// export type InitialMessageAction = {
// 	message_type: MessageType.InitialMessageAction
// 	data: { name: string }
// }

// export type ServerArenaEvent = {
// 	message_type: MessageType.ServerArenaEvent
// 	data: { arena_message: ArenaMessage }
// }

// export type ServerArenaAction = {
// 	message_type: MessageType.ServerArenaAction
// 	data: {
// 		arena_message: ArenaMessage
// 	}
// }

// export type JoinArenaAction = {
// 	message_type: MessageType.JoinArenaAction
// 	data: {
// 		arena_name: string
// 	}
// }

// export type JoinArenaEvent = {
// 	message_type: MessageType.JoinArenaEvent
// 	data: {
// 		success: boolean
// 	}
// }

// export type CreateArenaEvent = {
//
// }
//
// export type CreateArenaAction = {
//
// }

export type ArenaMessage = {
	message_type: ArenaMessage
	data: {
		message_type: ArenaMessageType
		data: any
	}
}

