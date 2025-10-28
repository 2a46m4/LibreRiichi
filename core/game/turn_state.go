package game

import "fmt"

type TurnType uint8

const (
	CHOOSING_DISCARD TurnType = iota
	POST_ACTIONS_PENDING
	GAME_ENDED
)

type TransitionType uint8

const (
	DRAW_TRANSITION TransitionType = iota
	TOSS_TRANSITION
	PLAY_POST_ACTIONS_TRANSITION
	WIN_OR_EXHAUST_TRANSITION
	RESET_TRANSITION
)

type BadStateTransition struct{}

func (BadStateTransition) Error() string {
	return "Bad state transition"
}

type TurnState struct {
	TurnNumber uint8
	TurnType   TurnType
}

func (coord *TurnState) transition(intoState TurnType, validStates ...TurnType) error {
	for _, valid := range validStates {
		if coord.TurnType == valid {
			coord.TurnType = intoState
			return nil	
		} 
	}
	return BadStateTransition{}
}

// Returns the next action (TurnType), or error if not a valid move
func (coord *TurnState) Transition(action TransitionType, fromPlayer uint8) (TurnType, error) {
	switch action {
	case DRAW_TRANSITION:
		coord.transition(CHOOSING_DISCARD, POST_ACTIONS_PENDING)
	case PLAY_POST_ACTIONS_TRANSITION:
		coord.transition(CHOOSING_DISCARD, POST_ACTIONS_PENDING)
	case RESET_TRANSITION:
		coord.transition(CHOOSING_DISCARD, GAME_ENDED)	
	case TOSS_TRANSITION:
		coord.transition(POST_ACTIONS_PENDING, CHOOSING_DISCARD)
	case WIN_OR_EXHAUST_TRANSITION:
	default:
		panic(fmt.Sprintf("unexpected game.TransitionType: %#v", action))
	}
}

func (coord TurnState) GetPlayer() uint8 {
	return coord.TurnNumber % 4
}
