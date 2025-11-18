package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type DiscardPile struct {
	Discards [26]Tile
	index uint8
}

func (pile *DiscardPile) Add(tile Tile) {
	pile.Discards[pile.index] = tile
	pile.index += 1
}

func (pile *DiscardPile) Remove() (tile Tile) {
	tile = pile.Discards[pile.index]
	pile.index -= 1
	return tile
}

func (pile *DiscardPile) Reset() {
	pile.index = 0
}

func (pile *DiscardPile) IsEmpty() bool {
	return pile.index == 0
}

func (pile DiscardPile) Last() Tile {
	if pile.index == 0 {
		panic("Bad state: No tiles have been discarded yet")
	} else {
		return pile.Discards[pile.index - 1]
	}
}
