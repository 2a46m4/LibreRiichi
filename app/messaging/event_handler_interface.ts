import {Message} from "./message";

export abstract class EventHandlerInterface {
    abstract dispatch_message(data: Message): void;
}