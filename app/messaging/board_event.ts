import {ArenaMessageType} from "./arena_message";

export enum BoardEventType {
    // An action that a player performed, affecting the board state
    PlayerActionEventType,
    // A potential action available to the player
    PotentialActionEventType,
    // A setup event
    GameSetupEventType,
    // A game end event
    GameEndEventType
}

type MessageEntry<T extends BoardEventType = BoardEventType, D = any> = {
    event_type: T
    data: D
}

type MessageMap = {
    [BoardEventType.PlayerActionEventType]: {
        action_data: any,
        from_player: number
    }
    [BoardEventType.PotentialActionEventType]: {
        action_data: any
    }
    [BoardEventType.GameSetupEventType]: {
        setup: {setup_type: number, data: any}[]
    }
    [BoardEventType.GameEndEventType]: {
        result: any
    }
}

type ConstrainedMap<M extends Record<BoardEventType, any>> = {
    [K in keyof M & BoardEventType]: MessageEntry<K, M[K]>
}

export type BoardEvent = ConstrainedMap<MessageMap>[keyof MessageMap]
