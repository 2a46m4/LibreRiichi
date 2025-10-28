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

	if roundState.RoundWind == North {
		roundState.RoundWind = East
	} else {
		roundState.RoundWind += 1
	}
}
