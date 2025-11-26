import {
  ServerAction as ServerActionData,
  ServerActionType,
} from './server_action_generated'
import {
  ServerResponse as ServerResponseData,
  ServerResponseType,
} from './server_response_generated'
import { ServerEvent as ServerEventData, ServerEventType } from './server_event_generated'

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
