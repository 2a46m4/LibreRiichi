import { AgentInfo } from "./agent_info";
import {create_event, create_fsm_builder, create_state, FSM} from "../fsm";
import {Tile} from "./tile";

interface NoData {}

interface DiscardData {
    tile?: Tile
}

interface OtherPlayerTurn {
    player_id: number
}

export class Arena {
  constructor(
    public agents: AgentInfo[],
    public game_started: boolean,
    public player_index: number,
    public dealer: number
  ) {
      const awaiting_discard_state = create_state<NoData>("awaiting_discard", {})
      const discarded_state = create_state<DiscardData>("discarded_state", {})
      const other_player_state = create_state<OtherPlayerTurn>("other_player_turn", {
          player_id: dealer,
      })
      const naki_call_state = create_state<NoData>("naki_call_state", {})

      const draw_event = create_event("draw_event", awaiting_discard_state, )

      create_fsm_builder()


  }
}
