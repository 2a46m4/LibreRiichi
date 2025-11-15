package game

import (
	"context"
	"fmt"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	"github.com/looplab/fsm"
)

type GameState struct {
	*fsm.FSM
	context context.Context
}

func InitGameState() GameState {
	return GameState{
		fsm.NewFSM(
			"out-of-game",
			fsm.Events{
				fsm.EventDesc{
					Name: "start-game",
					Src: []string{
						"out-of-game",
					},
					Dst: "in-game",
				},
				fsm.EventDesc{
					Name: "start-round",
					Src: []string{
						"in-game",
					},
					Dst: "in-round,pre-draw",
				},
				fsm.EventDesc{
					Name: "draw-tile",
					Src: []string{
						"in-round,pre-draw",
						"in-round,naki-called",
					},
					Dst: "in-round,waiting-discard",
				},
				fsm.EventDesc{
					Name: "discard-tile",
					Src: []string{
						"in-round,waiting-discard",
					},
					Dst: "in-round,waiting-naki",
				},
				fsm.EventDesc{
					Name: "call-naki",
					Src: []string{
						"in-round,waiting-naki",
					},
					Dst: "in-round,waiting-discard",
				},
				fsm.EventDesc{
					Name: "no-naki",
					Src: []string{
						"in-round,waiting-naki",
					},
					Dst: "in-round,pre-draw",
				},
				fsm.EventDesc{
					Name: "round-draw",
					Src: []string{
						"in-round,naki-finished",
					},
					Dst: "round-ended",
				},
				fsm.EventDesc{
					Name: "round-win",
					Src: []string{
						"in-round,waiting-discard",
						"in-round,in-round,waiting-naki",
					},
					Dst: "in-game",
				},
				fsm.EventDesc{
					Name: "game-end",
					Src: []string{
						"in-game",
					},
					Dst: "out-of-game",
				},
			},
			fsm.Callbacks{
				"start-game": func(context context.Context, event *fsm.Event) {
					tileState := event.Args[0].(*TileState)
					fmt.Println(tileState)
				},
			},
		),
		context.Background(),
	}
}

func (gameState *GameState) Transition(event string, arguments ...any) error {
	return gameState.FSM.Event(gameState.context, event, arguments...)
}

func (gameState *GameState) HandleEvent(
	event Action,
	gameIdx uint8,
	tileState *TileState,
	roundState *RoundState,
	windState *WindState,
	turnState *TurnState,
) {

}

func (gameState *GameState) HandleStartGame() {

}
