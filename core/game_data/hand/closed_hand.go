package core

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type ClosedHand struct {
	hand  [14]Tile
	index uint8
}

func (closed *ClosedHand) Add(tile Tile) {
	closed.hand[closed.index] = tile
	closed.index += 1
}

func (closed *ClosedHand) RemoveTile(tiles ...Tile) {
	for _, tile := range tiles {
		found := false
		for handIdx, handTile := range closed.hand {
			if tile == handTile {
				Swap(closed.hand[:], uint(handIdx), uint(closed.index-1))
				closed.index -= 1
				found = true
				break
			}
		}

		if !found {
			panic("Couldn't find tile")
		}
	}
}
