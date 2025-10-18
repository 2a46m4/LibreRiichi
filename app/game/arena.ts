import { AgentInfo } from "./agent_info";

export interface IArena {
  agents: AgentInfo[]
  game_started: boolean
}

export class Arena implements IArena {
  constructor(public agents: AgentInfo[], public game_started: boolean) { }

  map_idx() {

  }

}