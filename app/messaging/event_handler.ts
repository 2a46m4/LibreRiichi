import { IncomingMessage, MessageType } from './message'
import { ServerResponse, ServerResponseType } from './server_response_generated'
import { ServerEvent } from './server_event_generated'

export class EventHandler<TIncoming> {
  private listeners: Array<(data: TIncoming) => boolean> = []

  constructor() {}

  handle(data: TIncoming): void {
    this.listeners = this.listeners.filter((listener) => listener(data))
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
export const ClickBus = new EventHandler<MouseEvent>()
window.addEventListener('click', (event) => {
  ClickBus.handle(event)
})

// Registers a wait for a request
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

// Registers a wait for a request, and throws an error if the response failed
export async function get_response(
  idx: number,
  bus = ServerMessageBus,
): Promise<ServerResponse> {
  let response = await register_request(idx, bus)
  switch (response.serverresponse_type) {
    case ServerResponseType.GenericResponse:
      if (response.success) {
        console.log('Response success')
      } else {
        throw new Error(response.fail_reason)
      }
    default:
      throw new Error('Expected default response type')
  }
}

// Generator from an async function
export async function* make_async_generator_from_event<T>(
  event_handler: EventHandler<T>,
) {
  while (true) {
    const promise = new Promise<T>((res) => {
      event_handler.register((data: T) => {
        res(data)
        return false
      })
    })
    yield await promise
  }
}
