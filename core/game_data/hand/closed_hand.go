package core

import (
	"errors"
	"slices"

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
	return slices.Contains(closed.hand[:closed.index], tile)
}

func (closed *ClosedHand) HasTileN(tile Tile) (count int) {
	return Count(closed.hand[:closed.index], tile)
}

func (closed *ClosedHand) Last() (tile Tile, err error) {
	if closed.index == 0 {
		return tile, errors.New("out of bounds")
	}

	return closed.hand[closed.index-1], nil
}
