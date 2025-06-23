import {IncomingMessage} from "../message";

export abstract class EventHandlerInterface<T> {
    abstract handle_message(ev: MessageEvent);
}