import { AgentInfo } from "./agent_info";

export class Arena {
  game_seating: number[] | null = null

  constructor(
    public agents: AgentInfo[],
    public game_started: boolean,
    public player_index: number
  ) { }

  set_game_seating(player_to_order: number[]): void {
    this.game_seating = player_to_order // Maps arena index to game index
  }

  // Maps arena index to game index
  agent_to_game(idx: number): number {
    if (this.game_seating === null) {
      throw new Error("Seating is null")
    } else {
      return this.game_seating[idx]
    }
  }

  // Maps game to seat index
  game_to_seating(idx: number): number {
    // Seating is always index 0, clockwise increasing
    return ((idx - this.player_index) + 4) % 4
  }
}
