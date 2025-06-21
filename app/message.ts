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

export function getMatchingMessageType(type: MessageType): MessageType {
	const map = {
		[MessageType.InitialMessageAction]: MessageType.InitialMessageEvent,
		[MessageType.JoinArenaAction]: MessageType.JoinArenaEvent,
		[MessageType.ServerArenaAction]: MessageType.ServerArenaEvent,
		[MessageType.ListArenasAction]: MessageType.ListArenasEvent,
		[MessageType.CreateArenaAction]: MessageType.CreateArenaEvent,

		[MessageType.InitialMessageEvent]: MessageType.InitialMessageAction,
		[MessageType.JoinArenaEvent]: MessageType.JoinArenaAction,
		[MessageType.ServerArenaEvent]: MessageType.ServerArenaAction,
		[MessageType.ListArenasEvent]: MessageType.ListArenasAction,
		[MessageType.CreateArenaEvent]: MessageType.CreateArenaAction,
	} as const

	return map[type]
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

export type ArenaMessage = {
	message_type: ArenaMessage
	data: {
		message_type: ArenaMessageType
		data: any
	}
}

