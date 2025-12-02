package game

import (
	"fmt"
	util "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

// TODO
type GameIdx uint8
type ArenaIdx uint8

// The game index stays the same throughout rounds
//
// Round increments should increment the dealer index (if needed). So the next round after player 0 was dealer should be player
// All board events should send the game index
type Ordering struct {
	// Maps Arena Index → Game Index
	ArenaToGame [4]uint8
	// Maps Game Index → Arena Index
	GameToArena [4]uint8
}



// Create a random ordering
func InitRandomOrdering() (ordering Ordering) {
    for i := range uint8(4) {
	ordering.ArenaToGame[i] = i
    }
    util.PermuteArray(ordering.ArenaToGame[:])
    for arenaIdx, gameIdx := range ordering.ArenaToGame {
	ordering.GameToArena[gameIdx] = uint8(arenaIdx)
    }
    return ordering
}

func (ordering Ordering) GameIdx(arenaIdx uint8) uint8 {
	return ordering.ArenaToGame[arenaIdx]
}

func (ordering Ordering) ArenaIdx(gameIdx uint8) uint8 {
	return ordering.GameToArena[gameIdx]
}

// String returns a pretty-printed representation of the ordering mapping
func (ordering Ordering) String() string {
	result := "Ordering Mapping:\n"
	result += "Arena -> Game: "
	for arenaIdx, gameIdx := range ordering.ArenaToGame {
		result += fmt.Sprintf("%d -> %d", arenaIdx, gameIdx)
		if arenaIdx < 3 {
			result += ", "
		}
	}
	result += "\nGame -> Arena: "
	for gameIdx, arenaIdx := range ordering.GameToArena {
		result += fmt.Sprintf("%d -> %d", gameIdx, arenaIdx)
		if gameIdx < 3 {
			result += ", "
		}
	}
	return result
}
