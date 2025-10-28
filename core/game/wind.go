package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

// Game to Wind offset
// Stores the offset that player 0 is to the East wind
// e.g.
// Player 0 is East wind  → offset 0
// Player 0 is South wind → offset 1
type WindState uint8

func (windState WindState) GetPlayerWind(gameIdx uint8) Wind {
	return Wind((gameIdx + uint8(windState)) % 4)
}

func (windState *WindState) IncrementWind() {
	*windState = (*windState + 1) % 4
}

func (windState WindState) GetPlayerWindArena(ordering Ordering, arenaIdx uint8) Wind {
	return Wind((ordering.GameIdx(arenaIdx) + uint8(windState)) % 4)
}
