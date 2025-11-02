package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type RoundState struct {
	RoundWind Wind
	RoundNumber uint8
}

func (roundState* RoundState) IncrementRound() {
	roundState.RoundNumber += 1
	if roundState.RoundNumber == 4 {
		roundState.RoundWind += 1
	}
}
