package core

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type ClosedHand struct {
	hand  [14]Tile
	index uint8
}

func (closed *ClosedHand) Add(tiles ...Tile) {
	for _, tile := range tiles {
		closed.hand[closed.index] = tile
		closed.index += 1
	}
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

func (closed *ClosedHand) GetHand() []Tile {
	return closed.hand[:closed.index]
}

func (closed *ClosedHand) HasTile(tile Tile) bool {
	for i := uint8(0); i < closed.index; i++ {
		if closed.hand[i] == tile {
			return true
		}
	}
	return false
}
