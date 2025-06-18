import {Connection} from "./connection";

export enum ApplicationState {
    NOT_CONNECTED,
    CONNECTING,
    CONNECTED,
    JOINING_ROOM,
    JOINED_ROOM,
    CREATING_ROOM,
    IN_GAME,
}

export enum GameState {
    OUT_OF_TURN,
    IN_TURN
}

export class Application {
    connection: Connection;
    username: string

    constructor() {

    }

    connect() {

    }
}