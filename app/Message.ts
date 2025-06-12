export enum MessageType {
    InitialMessage,
    ArenaMessage,
    InitialMessageReturn
}

export enum ArenaMessageType {

}

type Message = {
    message_type: MessageType
    data: any
}

type ArenaMessage = {
    message_type: ArenaMessage
    data: {
        message_type: ArenaMessageType
        arena_data: any
    }
}