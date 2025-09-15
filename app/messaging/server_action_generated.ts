import {ArenaActionMessage} from "./arena_action_generated";

export type ServerActionMessage = InitialMessageAction | JoinArenaAction | ServerArenaAction | ListArenasAction | CreateArenaAction | ArenaInfoAction;

export enum ServerActionType {
    InitialMessageAction = 0,
    JoinArenaAction = 1,
    ServerArenaAction = 2,
    ListArenasAction = 3,
    CreateArenaAction = 4,
    ArenaInfoAction = 5,
}

export interface InitialMessageAction {
    serveraction_type: ServerActionType.InitialMessageAction
    name: string;
}

export interface JoinArenaAction {
    serveraction_type: ServerActionType.JoinArenaAction
    arena_name: string;
}

export interface ServerArenaAction {
    serveraction_type: ServerActionType.ServerArenaAction
    arena_action: ArenaActionMessage;
}

export interface ListArenasAction {
    serveraction_type: ServerActionType.ListArenasAction
}

export interface CreateArenaAction {
    serveraction_type: ServerActionType.CreateArenaAction
    arena_name: string;
}

export interface ArenaInfoAction {
    serveraction_type: ServerActionType.ArenaInfoAction
}
