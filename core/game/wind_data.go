package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

// Game to Wind offset
// Stores the offset that player 0 is to the East wind
// e.g.
// Player 0 is East wind  → offset 0
// Player 0 is South wind → offset 1
type WindData uint8

func (windData WindData) GetPlayerWind(gameIdx uint8) Wind {
	return Wind((gameIdx + uint8(windData)) % 4)
}

func (windData *WindData) IncrementWind() {
	*windData = (*windData + 1) % 4
}

func (windData WindData) GetPlayerWindArena(ordering Ordering, arenaIdx uint8) Wind {
	return Wind((ordering.GameIdx(arenaIdx) + uint8(windData)) % 4)
}
