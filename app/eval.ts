import {Connection, websocket_address} from "./messaging/connection";
import {IncomingMessage, MessageType, validate_message} from "./messaging/message";
import {EventHandler, keep_registered, register_request} from "./messaging/event_handler";
import {ServerActionType} from "./messaging/server_action_generated";
import {ServerResponseType} from "./messaging/server_response_generated";
import {ServerEventMessage, ServerEventType} from "./messaging/server_event_generated";
import {ArenaEventType} from "./messaging/arena_event_generated";
import {ArenaActionType} from "./messaging/arena_action_generated";
import {BoardEvent, BoardEventType} from "./messaging/board_event_generated";
import { SetupType } from "./types/setup";
import { decode, Tile } from "./game/tile";

let busses: EventHandler<IncomingMessage>[] = Array(4).fill(0).map(() => new EventHandler())

let arena_busses = Array(4).fill(0).map((_, idx)=>{
  let arena_bus: EventHandler<ServerEventMessage> = new EventHandler()
  busses[idx].register(keep_registered((data: IncomingMessage) => {
    if (data.message_type === MessageType.EVENT) {
      arena_bus.handle(data.data as ServerEventMessage)
    }
  }))

  return arena_bus
})

let game_state = [{},{},{},{}]

let callback = (i: number, data: ServerEventMessage) => {
  console.log("Arena listener called")
  switch (data.serverevent_type) {
	case ServerEventType.ServerArenaEvent:
	  let message = data.arena_message
	  switch (message.arenaevent_type) {
		case ArenaEventType.GameStartedEvent:
		  break;
		case ArenaEventType.PlayerJoinedEvent:
		  break;
		case ArenaEventType.PlayerQuitEvent:
		  break;
		case ArenaEventType.ArenaBoardEvent:
		  handle_game(i, message.board_event)
		  break;
		default:
		  throw new Error("Unknown arena event")
	  }
	  break;
	default:
	  throw new Error("Unknown server event")
  }
  return true
}

for (let i = 0; i < 4; i++) {
  arena_busses[i].register((data)=>callback(i, data))
}

let conns = []
for (let i = 0; i < 4; i++) {
  conns.push(new Connection(new WebSocket(websocket_address), (data: MessageEvent) => {
	let msg = JSON.parse(data.data) as IncomingMessage
	if (validate_message(msg).isValid) {
	  busses[i].handle(msg)
	}
  }))
}

Promise.all(conns.map(async (conn, idx) => {
  await conn.wait_until_ready()
  console.log("Connection ready!")
  let return_index = conn.send({
	message_type: MessageType.REQUEST,
	data: {
	  serveraction_type: ServerActionType.InitialMessageAction,
	  name: Date().toString()
	}
  })

  let message_return = await register_request(return_index, busses[idx])
  if (message_return.serverresponse_type !== ServerResponseType.GenericResponse) {
	console.error("Unexpected message type")
	return
  }

  if (!message_return.success) {
	console.error("Failed to connect: " + message_return.fail_reason)
	return
  }
})).then(async () => {
  let msg_idx = conns[0].send({
	message_type: MessageType.REQUEST,
	data: {
	  serveraction_type: ServerActionType.CreateArenaAction,
	  arena_name: "room_test",
	}
  })

  let msg = await register_request(msg_idx, busses[0])

  if (msg === undefined) {
	console.error("Connection error: Failed to create room")
  }

  if (msg.serverresponse_type !== ServerResponseType.GenericResponse) {
	throw new Error("Connection error: Wrong type")
  }

  if (!msg.success) {
	  console.error("Could not create room: " + msg.fail_reason)
  }
}).then(()=>{
  return Promise.all(conns.map(async (conn, idx)=>{

	let msg_idx = conn.send(
      {
        message_type: MessageType.REQUEST,
        data: {
          serveraction_type: ServerActionType.JoinArenaAction,
          arena_name: "room_test",
        }
      }
    )

	let ret = await register_request(msg_idx, busses[idx])
    if (ret.serverresponse_type !== ServerResponseType.GenericResponse) {
      throw new Error("Connection error: Wrong type")
    }

    if (!ret.success) {
      throw new Error("Could not join room: " + ret.fail_reason)
    }

    console.log("Joined room: ", idx)
  }))
}).then(async () => {
  let msg_idx = conns[0].send(
      {
        message_type: MessageType.REQUEST,
        data: {
          serveraction_type: ServerActionType.ServerArenaAction,
          arena_action: {
            arenaaction_type: ArenaActionType.StartGameActionData
          }
        }
      }
	)

	let ret = await register_request(msg_idx);
	if (ret.serverresponse_type !== ServerResponseType.GenericResponse) {
      throw new Error("Connection error: Wrong Type")
	}

	if (!ret.success) {
      throw new Error("Couldn't start game: " + ret.fail_reason)
	}

})

let game_started = false
let initial_hand = [new Array(), new Array(), new Array(), new Array()]
let dora = new Tile(0)
let points = Array(4).fill(0)
let player_n = new Uint8Array()
let player_order = Array(4).fill(0)
let round_wind = 0
let round_number = 0

function handle_game(i: number, event: BoardEvent) {
	game_started = true

    if (event.boardevent_type !== BoardEventType.GameSetupEvent) {
        throw new Error("Wrong event type")
    }

    for (let setup of event.setup) {
        switch (setup.setup_type) {
			case SetupType.INITIAL_TILES:
				initial_hand[i] = Tile.from(setup.data)
				break
			case SetupType.DORA:
				dora = new Tile(setup.data)
				break
			case SetupType.STARTING_POINTS:
				points[i] = setup.data[i]
				break
			case SetupType.PLAYER_NUMBER:
				player_n = setup.data
				break
			case SetupType.PLAYER_ORDER:
				player_order = decode(setup.data)
				break
			case SetupType.ROUND_WIND:
				round_wind = setup.data
				break
			case SetupType.ROUND_NUMBER:
				round_number = setup.data
				break
        }
    }
}
