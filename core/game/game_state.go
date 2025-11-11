package game

import (
	"fmt"

	"github.com/looplab/fsm"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=TurnType

type TurnType uint8

//go:generate go run golang.org/x/tools/cmd/stringer -type=TransitionType

type TransitionType uint8

const (
	INVALID TurnType = iota
	OUT_OF_GAME
	IN_GAME
	AWAITING_DISCARD
	DISCARDED
	AWAITING_NAKI
	NAKI_CALLED
	NAKI_FINISHED
)

const (
	_                                    = iota
	GAME_START_TRANSITION TransitionType = iota << 3
	DRAW_TRANSITION
	TOSS_TRANSITION
	COMPUTE_NAKI_TRANSITION
	NAKI_CALLED_TRANSITION
	AWAITING_DISCARD_TRANSITION
	NO_NAKI_TRANSITION
	GAME_FINISHED_TRANSITION

	TRANSITION_MASK  uint8 = 0b11111000
	TRANSITION_SHIFT       = 3
)

type BadStateTransition struct {
	from       TurnType
	transition TransitionType
}

func (err BadStateTransition) Error() string {
	return fmt.Sprintf("Bad state transition: attempted to transition from %v with %v", err.from, err.transition)
}

type I uint8

var ACTION_TRANSITION_MTX [128]TurnType = [128]TurnType{
	I(OUT_OF_GAME) | I(GAME_START_TRANSITION): IN_GAME,

	I(IN_GAME) | I(DRAW_TRANSITION): AWAITING_DISCARD,

	I(AWAITING_DISCARD) | I(TOSS_TRANSITION): DISCARDED,

	I(DISCARDED) | I(COMPUTE_NAKI_TRANSITION):  AWAITING_NAKI,
	I(DISCARDED) | I(NO_NAKI_TRANSITION):       NAKI_FINISHED,
	I(DISCARDED) | I(GAME_FINISHED_TRANSITION): OUT_OF_GAME,

	I(AWAITING_NAKI) | I(NAKI_CALLED_TRANSITION): NAKI_CALLED,

	I(NAKI_CALLED) | I(AWAITING_DISCARD_TRANSITION): AWAITING_DISCARD,

	I(NAKI_FINISHED) | I(DRAW_TRANSITION):          AWAITING_DISCARD,
	I(NAKI_FINISHED) | I(GAME_FINISHED_TRANSITION): OUT_OF_GAME,
}

type GameState struct {
	*fsm.FSM
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
					Name: "enter-round",
					Src: []string{
						"in-game",
					},
					Dst: "in-round,pre-draw",
				},
				fsm.EventDesc{
					Name: "draw-tile",
					Src: []string{
						"in-round,pre-draw",
					},
					Dst: "in-round,waiting-discard",
				}
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
					Dst: "in-round,naki-called",
				},
				fsm.EventDesc{
					Name: "no-naki",
					Src: []string{
						"in-round,waiting-naki",
					},
					Dst: "in-round,naki-finished",
				},
				
			},
			fsm.Callbacks{},
		),
	}
}

func (coord *GameState) Try(action TransitionType) error {
	return nil
	// copy := *coord
	// return copy.Transition(action)
}

// Returns the next action TurnType
// func (coord *GameState) Transition(action TransitionType) error {
// 	newState := ACTION_TRANSITION_MTX[I(coord.TurnType)|I(action)]
// 	if newState != INVALID {
// 		coord.TurnType = newState
// 		return nil
// 	} else {
// 		return BadStateTransition{
// 			from:       coord.TurnType,
// 			transition: action,
// 		}
// 	}
// }
