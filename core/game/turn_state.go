package game

type TurnType uint8

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

type BadStateTransition struct{}

func (BadStateTransition) Error() string {
	return "Bad state transition"
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
	TurnNumber uint8
	TurnType   TurnType
}

func InitGameState() GameState {
	return GameState{
		TurnNumber: 0,
		TurnType:   OUT_OF_GAME,
	}
}

// Returns the next action TurnType
func (coord *GameState) Transition(action TransitionType, fromPlayer uint8) error {
	newState := ACTION_TRANSITION_MTX[I(coord.TurnType)|I(action)]
	if newState != INVALID {
		coord.TurnType = newState
		if action == DRAW_TRANSITION {
			coord.TurnNumber += 1
		}
		return nil
	} else {
		return BadStateTransition{}
	}
}

func (coord GameState) GetPlayer() uint8 {
	return coord.TurnNumber % 4
}
