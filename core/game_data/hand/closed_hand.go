package hand

import (
	"slices"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
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

func (closed *ClosedHand) Pop(n uint8) {
	closed.index -= n
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

func (closed *ClosedHand) HasTile(tiles ...Tile) bool {
	res := true
	for _, tile := range tiles {
		res = res && slices.Contains(closed.hand[:closed.index], tile)
	}

	return res
}

func (closed *ClosedHand) HasTileN(tile Tile) (count int) {
	return Count(closed.hand[:closed.index], tile)
}

func (closed *ClosedHand) Last() Tile {
	return closed.hand[closed.index-1]
}

func (closed *ClosedHand) Length() uint8 {
	return closed.index
}

func (closed ClosedHand) UniqueTiles() ([]Tile, []uint8) {
	unique := make([]Tile, 0)
	count := make([]uint8, 0)
	for _, tile := range closed.GetHand() {
		if uniqueIdx := slices.Index(unique, tile); uniqueIdx != -1 {
			count[uniqueIdx] += 1
		} else {
			unique = append(unique, tile)
			count = append(count, 1)
		}
	}
	return unique, count
}

func (closed *ClosedHand) SortInplace() {
	slices.Sort(closed.GetHand())
}
