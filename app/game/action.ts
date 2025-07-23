import {Tile} from "./tile";

export enum ActionType {
    RON,
    TSUMO,
    RIICHI,
    TOSS,
    SKIP,
    PON,
    KAN,
    CHII,
    DRAW
}

type MessageEntry<T extends ActionType = ActionType, D = any> = {
    action_type: T
    data: D
}

type MessageMap = {
    [ActionType.RON]: {
        tile_to_ron: Tile
        win_result:
    }
    [ActionType.TSUMO]: {},
    [ActionType.RIICHI]: {},
    [ActionType.TOSS]: {},
    [ActionType.SKIP]: {},
    [ActionType.PON]: {},
    [ActionType.KAN]: {},
    [ActionType.CHII]: {},
    [ActionType.DRAW]: {}
}

type ConstrainedMap<M extends Record<ActionType, any>> = {
    [K in keyof M & ActionType]: MessageEntry<K, M[K]>
}

export type Action = ConstrainedMap<MessageMap>[keyof MessageMap]
