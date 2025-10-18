import { AgentInfo } from "./agent_info";

export interface IArena {
  agents: AgentInfo[]
  game_started: boolean
  set_game_seating(player_to_order: number[]): void
  agent_to_game_seat(idx: number): number
}

export class Arena implements IArena {
  game_seating: number[] | null = null

  constructor(
    public agents: AgentInfo[],
    public game_started: boolean,
  ) { }

  set_game_seating(player_to_order: number[]): void {
    this.game_seating = player_to_order
  }

  agent_to_game_seat(idx: number): number {
    if (this.game_seating === null) {
      throw new Error("Seating is null")
    } else {
      return this.game_seating[idx]
    }
  }
}