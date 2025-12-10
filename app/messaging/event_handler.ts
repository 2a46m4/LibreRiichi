import { IncomingMessage, MessageType } from './message'
import { ServerResponse } from './server_response_generated'
import { ServerEvent } from './server_event_generated'
import { ECS } from '../render/ecs'

export class EventHandler<TIncoming> {
  private listeners: Array<(data: TIncoming) => boolean> = []

  constructor() { }

  handle(data: TIncoming): void {
    this.listeners.filter((listener) => listener(data))
  }

  register(listener: (data: TIncoming) => boolean): number {
    this.listeners.push(listener)
    return this.listeners.length - 1
  }

  unregister(index: number): void {
    this.listeners.splice(index, 1)
  }
}

export function keep_registered<TIncoming>(
  fn: (_: TIncoming) => void,
): (_: TIncoming) => true {
  return (data: TIncoming) => {
    fn(data)
    return true
  }
}

export const ServerMessageBus = new EventHandler<IncomingMessage>()
export const ArenaMessageBus = new EventHandler<ServerEvent>()
ServerMessageBus.register(
  keep_registered((data: IncomingMessage) => {
    if (data.message_type === MessageType.EVENT) {
      ArenaMessageBus.handle(data.data as ServerEvent)
    }
  }),
)
export const ClickEventBus = new EventHandler<ECS.EntityID>()

export function register_request(
  msg_idx: number,
  bus = ServerMessageBus,
): Promise<ServerResponse> {
  let { promise, resolve } = Promise.withResolvers<ServerResponse>()
  bus.register((msg) => {
    if (
      msg.message_type === MessageType.RESPONSE &&
      msg.message_index === msg_idx
    ) {
      resolve(msg.data)
      return false
    } else {
      return true
    }
  })
  return promise
}
