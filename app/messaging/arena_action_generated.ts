// ArenaAction Union Type and Enum
import {Action} from "../game/action";

export type ArenaActionMessage = StartGameActionData | PlayerQuitActionData | PlayerActionData;

export enum ArenaActionType {
    StartGameActionData = 0,
    PlayerQuitActionData = 1,
    PlayerActionData = 2,
}

// Individual struct interfaces

export interface StartGameActionData {
    arenaaction_type: ArenaActionType.StartGameActionData;
}

export interface PlayerQuitActionData {
    arenaaction_type: ArenaActionType.PlayerQuitActionData;
}

export interface PlayerActionData {
    arenaaction_type: ArenaActionType.PlayerActionData;
    action: Action;
}
