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
    }
    [ActionType.TSUMO]: {
        tile_to_tsumo: Tile
    },
    [ActionType.RIICHI]: {
        tile_to_riichi: Tile
    },
    [ActionType.TOSS]: {
        tile_to_toss: Tile
    },
    [ActionType.SKIP]: {
        action_to_skip: Action
    },
    [ActionType.PON]: {
        tile_to_pon: Tile
    },
    [ActionType.KAN]: {
        tile_to_kan: Tile
    },
    [ActionType.CHII]: {
        tile_to_chii: Tile
        tiles_in_hand: Tile[]
    },
    [ActionType.DRAW]: {
        drawn_tile: Tile
    }
}

type ConstrainedMap<M extends Record<ActionType, any>> = {
    [K in keyof M & ActionType]: MessageEntry<K, M[K]>
}

export type Action = ConstrainedMap<MessageMap>[keyof MessageMap]
