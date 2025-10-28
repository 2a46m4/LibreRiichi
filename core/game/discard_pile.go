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
