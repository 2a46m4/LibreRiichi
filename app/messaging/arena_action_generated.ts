// ArenaAction Union Type and Enum
import {Action} from "../game/action";

export type ArenaActionMessage =  ArenaAction & {
    arenaaction_type: ArenaActionType;
}

export type ArenaAction = StartGameActionData | PlayerQuitActionData | PlayerActionData;

export enum ArenaActionType {
    StartGameActionData = 0,
    PlayerQuitActionData = 1,
    PlayerActionData = 2,
}

// Individual struct interfaces

export interface StartGameActionData {
}

export interface PlayerQuitActionData {
}

export interface PlayerActionData {
    action: Action;
}
