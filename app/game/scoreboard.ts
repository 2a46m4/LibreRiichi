import { AgentInfo } from "./agent_info";

export type ScoreboardState = {
    scoreboard_values: number[],
    player_to_order_map: number[],
    round_wind: number,
    round_number: number,
    players: AgentInfo[]
    // The game index of the player
    player_idx: number,
}