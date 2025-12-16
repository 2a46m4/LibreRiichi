import { ServerAction as ServerActionData, ServerActionType } from './server_action_generated'
import { ServerResponse as ServerResponseData } from './server_response_generated'
import { ServerEvent as ServerEventData } from './server_event_generated'
import { ArenaActionType } from './arena_action_generated'
import { Action } from './action_generated'

export enum MessageType {
  RESPONSE = 0,
  REQUEST = 1,
  EVENT = 2,
}

export type Message = ServerAction | ServerResponse | ServerEvent

export type ServerAction = {
  message_type: MessageType.REQUEST
  data: ServerActionData
}

export type ServerResponse = {
  message_type: MessageType.RESPONSE
  data: ServerResponseData
}

export type ServerEvent = {
  message_type: MessageType.EVENT
  data: ServerEventData
}

export type IncomingMessage = Message & { message_index: number }

export function send_action_request(action: Action): Message {
	return {
		message_type: MessageType.REQUEST,
		data: {
			serveraction_type: ServerActionType.ServerArenaAction,
			arena_action: {
				arenaaction_type: ArenaActionType.PlayerActionData,
				action: action,
			},
		},
	}
}
