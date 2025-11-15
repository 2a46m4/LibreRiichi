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

func InitGameState() *GameState {
	gameState := GameState{
		FSM:     &fsm.FSM{},
		context: nil,
	}

	gameState.FSM = fsm.NewFSM(
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
				Dst: "in-round",
			},
			fsm.EventDesc{
				Name: "in-round-event",
				Src: []string{
					"in-round",
				},
				Dst: "in-round",
			},
			fsm.EventDesc{
				Name: "round-end",
				Src: []string{
					"in-round",
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
			"start-game": gameState.HandleStartGame,
		},
	)
	return &gameState
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

func (gameState *GameState) HandleStartGame(context context.Context, event *fsm.Event) {

}
