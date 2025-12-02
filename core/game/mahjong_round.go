package game

import (
	"fmt"
	"slices"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

// Stores data of a round
type MahjongRound struct {
	roundState RoundState
	data       MahjongRoundData
}

func InitMahjongRound() MahjongRound {
	return MahjongRound{
		roundState: InitRoundState(),
		data:       InitMahjongRoundData(),
	}
}

func (round *MahjongRound) ContinueMahjongRound() {
	round.data.IncrementRound()
	round.roundState.Transition("start-round", round, false)
}

func CheckIfValidAction(action Action, gameIdx uint8, round *MahjongRound) bool {
    switch action := action.(type) {
    case Chii:
	if !slices.Contains(round.roundState.RoundFSM.AvailableTransitions(), "call-naki") {
	    return false
	}

    
	

    case Draw:
    case Kan:
    case Pon:
    case Riichi:
    case Ron:
    case Skip:
    case Toss:
    case Tsumo:
    default:
	panic(fmt.Sprintf("unexpected core.Action: %#v", action))
    }
}
