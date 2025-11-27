<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { Tile } from '../game/tile'
import { initialize_tiles } from '../render/tile'
import { BoardEvent, BoardEventType } from '../messaging/board_event_generated'
import { Setup, SetupType } from '../game/setup'
import { ArenaMessageBus } from "../messaging/event_handler";
import { ArenaEventType } from "../messaging/arena_event_generated";
import { ServerEvent, ServerEventType } from "../messaging/server_event_generated";
import ScoreBoard from "../components/scoreboard.vue"
import { Action, ActionType } from "../messaging/action_generated";
import { IActionAnimator, IRenderer, ISelectionManager, ThreeJSRenderer } from "../render/renderer";
import { Arena } from '../game/arena'
import { create_event, create_fsm_builder, create_state } from "../fsm";

const props = defineProps<{ in_game: boolean, arena: Arena }>()

interface Add {
	tag: "add",
	add: Tile,
	location: number,
}
interface Remove {
	tag: "remove",
	id: string,
	location: number,
}
interface NewSet {
	tag: "newset",
	set: Tile[],
}
type HandAction = Add | Remove | NewSet

type GameIdx = (0 | 1 | 2 | 3) & { __brand: 'GameIdx' }
type TableIdx = (0 | 1 | 2 | 3) & { __brand: 'TableIdx' }

const arena_data = ref<{
	dealer_idx: number,
	round_wind: number,
	round_number: number,
	players: string[],
	player_idx: GameIdx,
}>({
	dealer_idx: 0,
	round_wind: 0,
	round_number: 0,
	players: [],
	player_idx: 0 as GameIdx, // Game idx
})

const game_data = {
	tiles: new Array<{ tile: Tile, id: string }>(),
	trigger_action: ref<HandAction>()
}

// Self is table idx 0
// So if we are player 3, our offset is (4 - 3) % 4 = 1
function offset_to_self(other_index: GameIdx): TableIdx {
	return (arena_data.value.player_idx + other_index) % 4 as TableIdx
}

const make_fsm = () => {

	const out_of_game = create_state("out_of_game", {})
	const awaiting_discard = create_state("awaiting_discard", {})
	const other_players_waiting_discard = create_state("other_players_waiting_discard", {})
	const discarded = create_state("discarded", {})
	const awaiting_other_player = create_state("awaiting_other_player", {
		player_id: 0,
	})
	const awaiting_naki_calls = create_state("awaiting_naki_calls", {})
	const naki_called = create_state("naki_called", {})
	const round_finished = create_state("round_finished", {})
	const game_finished = create_state("game_finished", {})

	const player_starts = create_event("player_starts", out_of_game, awaiting_discard, {
		callback: (_, __, tile_received: Tile) => {
			game_data.trigger_action.value = {
				tag: "add",
				add: tile_received,
				location: game_data.tiles.length,
			}
		}
	})
	const player_waiting = create_event("player_waiting", out_of_game, awaiting_other_player, {
		callback: (_, __, player_idx: number) => {
			// Do nothing
			console.log("Waiting for player ", player_idx)
		}
	})
	const draw_event = create_event("draw_event", awaiting_other_player, awaiting_discard, {
		callback: (_, __, tile_received: Tile, from_player: number) => {
			const table_idx = offset_to_self(from_player as GameIdx)
			action_animator.draw(table_idx, tile_received)
		}
	})
	const own_draw_event = create_event("own_draw_event", naki_called, awaiting_discard, {
		callback: (_, __, tile_received: Tile) => {
			action_animator.draw(0, tile_received)
		}
	})

	const fsm = create_fsm_builder()
		.add_state(out_of_game)
		.add_state(awaiting_discard)
		.add_state(discarded)
		.add_state(awaiting_other_player)
		.add_state(awaiting_naki_calls)
		.add_state(naki_called)
		.add_state(round_finished)
		.add_state(game_finished)
		.add_event(player_starts)
		.add_event(player_waiting)
		.add_event(draw_event)
		.add_event(own_draw_event)
		.build(out_of_game)

	fsm.trigger_event("player_starts", new Tile(5))
	// fsm.trigger_event("testing")

	ArenaMessageBus.register((data: ServerEvent) => {
		if (data.arena_message.arenaevent_type != ArenaEventType.ArenaBoardEvent) {
			return true
		}

		const board_event = data.arena_message.board_event
		switch (board_event.boardevent_type) {
			case BoardEventType.PotentialActionEvent:
				// TODO: Prompt the player to perform some action
				throw new Error("Not yet implemented")

			case BoardEventType.PlayerActionEvent:
				handle_player_action_event(board_event.action_data, board_event.from_player)
				break;
			case BoardEventType.GameSetupEvent:
				handle_game_setup_event(board_event.setup)
				break;
			case BoardEventType.GameEndEvent:

				break;
		}

		return true
	})

	return fsm
}


const three_canvas = ref<HTMLCanvasElement>()
const game_container = ref<HTMLDivElement>()
const is_fullscreen = ref(false)

const discard_required = ref(false)

let animation_id: number

let renderer: IRenderer
let selection_manager: ISelectionManager
let action_animator: IActionAnimator

let fsm: ReturnType<typeof make_fsm>

onMounted(() => {
	if (!three_canvas.value) return
	initialize_tiles()

	const manager = new ThreeJSRenderer(three_canvas.value)
	renderer = manager
	selection_manager = manager
	action_animator = manager
	animate(0, 0)

	window.addEventListener('resize', on_window_resize)
	window.addEventListener('click', on_click)

	fsm = make_fsm()

	// Watch updates to the tile set
	watch(game_data.trigger_action, (action) => {
		if (action === undefined) {
			return
		}

		switch (action.tag) {
			case "newset":
				action_animator.clear_tiles()
				const uuids = action.set.map(tile => action_animator.add_tile(tile))
				game_data.tiles = []
				for (let i = 0; i < uuids.length; i++) {
					game_data.tiles.push({
						id: uuids[i],
						tile: action.set[i]
					})
				}
				break;
			case "add":
				const uuid = action_animator.add_tile(action.add, 0, action.location)
				game_data.tiles.splice(action.location, 0, {
					id: uuid,
					tile: action.add,
				})
				break;
			case "remove":
				game_data.tiles.splice(action.location, 1)
				action_animator.remove_tile(action.id)
				break;
		}
	})

	ArenaMessageBus.register(message_handler)
	ArenaMessageBus.register(debug_message_printer)
})

function on_click(event: MouseEvent) {
	const selection = selection_manager.get_selection()
	if (selection === null) {
		return
	}

	if (!discard_required.value) {
		return
	}

	action_animator.remove_tile(selection.id)
	discard_required.value = false
}

function on_window_resize() {
	if (!three_canvas.value || !game_container.value) return

	const width = game_container.value.clientWidth
	const height = game_container.value.clientHeight

	three_canvas.value.width = width
	three_canvas.value.height = height

	renderer.window_resize(width, height)
}

function animate(t: number, dt: number) {
	renderer.animate_frame(dt)
	animation_id = requestAnimationFrame(new_t => { animate(new_t, new_t - t) })
}

onUnmounted(() => {
	if (animation_id) {
		cancelAnimationFrame(animation_id)
	}
	window.removeEventListener('resize', on_window_resize)
	window.removeEventListener('click', on_click)

	renderer.stop()
})

function debug_message_printer(event: ServerEvent) {
	const msg = `ServerEvent message: ${ServerEventType[event.serverevent_type]}\n\t`

	const debug_board_event = function (board_event: BoardEvent, msg: string) {
		msg += `BoardEvent message: ${BoardEventType[board_event.boardevent_type]}\n\t\t`
		switch (board_event.boardevent_type) {
			case BoardEventType.PlayerActionEvent:
				msg += `PlayerActionEvent message: ${ActionType[board_event.action_data.action_type]} from player ${board_event.from_player}`
				break
			case BoardEventType.PotentialActionEvent:
				msg += `PotentialActionEvent message: ${board_event.actions.map((action) => ActionType[action.action_type])}`
				break
			case BoardEventType.GameSetupEvent:
				msg += `GameSetupEvent message: ${board_event.setup.map((setup) => SetupType[setup.setup_type])}`
				break
			case BoardEventType.GameEndEvent:
				msg += `GameEndEvent message: ${board_event.result}`

		}
	}

	switch (event.arena_message.arenaevent_type) {
		case ArenaEventType.ArenaBoardEvent:
			debug_board_event(event.arena_message.board_event, msg)
	}

	return true
}

function message_handler(event: ServerEvent) {
	switch (event.arena_message.arenaevent_type) {
		case ArenaEventType.ArenaBoardEvent:
			handle_board_event(event.arena_message.board_event)
	}
	return true
}

function handle_board_event(new_event: BoardEvent) {
	if (new_event === undefined) {
		throw new Error('Event is undefined')
	}

	switch (new_event.boardevent_type) {
		case BoardEventType.PlayerActionEvent:
			handle_player_action_event(new_event)
			break
		case BoardEventType.PotentialActionEvent:
			handle_potential_action_event(new_event.actions)
			break
		case BoardEventType.GameSetupEvent:
			handle_game_setup_event(new_event.setup)
			break
		case BoardEventType.GameEndEvent:
			// Handle game end event
			break
		default:
			throw new Error("Unexpected")
	}
}

function handle_player_action_event(action: Action, from_player: number) {

	switch (action.action_type) {
		case ActionType.Tsumo:
			throw new Error("Win")
		case ActionType.Ron:
			throw new Error("Win")
		case ActionType.Riichi:
			break;
		case ActionType.Toss:
			break;
		case ActionType.Skip:
			break;
		case ActionType.Pon:
			break;
		case ActionType.Kan:
			break;
		case ActionType.Chii:
			break;
		case ActionType.Draw:
			fsm.trigger_event("draw_event", new Tile(action.drawn_tile), from_player)
			break;
	}
}

function handle_potential_action_event(actions: Action[]) {
	for (let action of actions) {
		switch (action.action_type) {
			case ActionType.Tsumo:
				break;
			case ActionType.Ron:
				break;
			case ActionType.Riichi:
				break;
			case ActionType.Toss:
				discard_required.value = true
				break;
			case ActionType.Skip:
				break;
			case ActionType.Pon:
				break;
			case ActionType.Kan:
				break;
			case ActionType.Chii:
				break;
			case ActionType.Draw:
				const tile = new Tile(action.drawn_tile)
				fsm.trigger_event("own_draw_event", tile)
				break;
		}
	}
}

function handle_game_setup_event(setups: Setup[]) {
	for (let setup of setups) {
		switch (setup.setup_type) {
			case SetupType.INITIAL_TILES:
				action_animator.clear_tiles()
				Tile.from(setup.data)
					.sort(Tile.sort)
					.forEach((tile) => {
						action_animator.add_tile(tile)
					})
				break
			case SetupType.DORA:
				action_animator.add_dora(new Tile(setup.data))
				break
			case SetupType.STARTING_POINTS:
				scoreboard_state.value.scoreboard_values = setup.data
				break
			case SetupType.PLAYER_NUMBER:
				scoreboard_state.value.player_idx = setup.data
				break
			case SetupType.PLAYER_ORDER:
				throw new Error("Should be removed")
			case SetupType.ROUND_WIND:
				scoreboard_state.value.round_wind = setup.data
				break
			case SetupType.ROUND_NUMBER:
				scoreboard_state.value.round_number = setup.data
				break
		}
	}
}
</script>

<template>
	<ScoreBoard v-if="in_game" :scoreboard_values="scoreboard_state.scoreboard_values"
		:player_to_order_map="scoreboard_state.player_to_order_map" :round_wind="scoreboard_state.round_wind"
		:round_number="scoreboard_state.round_number" :players="scoreboard_state.players"
		:player_idx="scoreboard_state.player_idx">
	</ScoreBoard>
	<div ref="game_container" class="game-view-container" :class="{ fullscreen: is_fullscreen }">
		<canvas ref="three_canvas" class="game-canvas"></canvas>
		<div class="game-ui"></div>
	</div>
</template>

<style scoped>
.game-view-container {
	position: relative;
	width: 100%;
	height: 100%;
	border-radius: 12px;
	overflow: hidden;
	cursor: pointer;
	transition: all 0.3s ease;
}

.game-canvas {
	width: 100%;
	height: 100%;
	display: block;
}

.game-ui {
	position: absolute;
	top: 16px;
	left: 16px;
	z-index: 10;
}

.game-info {
	background: rgba(0, 0, 0, 0.7);
	backdrop-filter: blur(10px);
	padding: 12px 16px;
	border-radius: 8px;
	color: white;
	font-size: 14px;
}

.game-info h3 {
	margin: 0 0 4px 0;
	font-size: 16px;
	font-weight: 600;
}

.game-info p {
	margin: 0;
	opacity: 0.8;
	font-size: 12px;
}
</style>
