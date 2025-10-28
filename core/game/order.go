package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type Ordering struct {
	// Maps Arena Index → Game Index
	ArenaToGame [4]uint8
	// Maps Game Index → Arena Index
	GameToArena [4]uint8
}

func (ordering *Ordering) InitRandom() {
	PermuteArray(ordering.ArenaToGame[:])
	for arenaIdx, gameIdx := range ordering.ArenaToGame {
		ordering.GameToArena[gameIdx] = uint8(arenaIdx)
	}
}

func (ordering Ordering) GameIdx(arenaIdx uint8) uint8 {
	return ordering.ArenaToGame[arenaIdx]
}

func (ordering Ordering) ArenaIdx(gameIdx uint8) uint8 {
	return ordering.GameToArena[gameIdx]
}
