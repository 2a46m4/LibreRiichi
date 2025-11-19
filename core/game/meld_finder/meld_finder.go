package meldfinder

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type MeldType uint8

const (
	Shuntsu MeldType = iota // Sequence
	Koutsu // Triplet
)

type WinningCombination struct {
	Melds [4]Meld
	Pair Tile
}

type Meld struct {
	Type MeldType
	FirstTile Tile
}

// Finds all combinations of winning Mahjong hands
func FindMelds(tiles []Tile, alreadyCompleted int) (ret []WinningCombination) {

	tileSet := NewMultiSet(tiles...)
	pairs := findPairs(tileSet)
	// Iterate over all possible combinations of pairs
	for _, pair := range pairs {
		tileSet.Remove(pair, pair)
		// Try to find 4 melds given that pair
		melds := findMelds(tileSet, alreadyCompleted)
		for _, meld := range melds {
			combo := WinningCombination{
				Melds: [4]Meld{},
				Pair:  pair,
			}
			copy(combo.Melds[:], meld)
			ret = append(ret, combo)
		}
	}

	return ret
}

// Find all lists of melds in the set that are at least of length 4 - alreadyCompleted
func findMelds(tileSet MultiSet[Tile], alreadyCompleted int) (result [][]Meld) {

	// Make a copy, since we don't want to modify the map while
	// iterating over it
	tileCopy := NewMultiSet(tileSet.ToSlice()...)
	removeAndFindMelds := func (tile1, tile2, tile3 Tile) [][]Meld {
		tileCopy.Remove(tile1, tile2, tile3)
		meldsList := findMelds(tileCopy, alreadyCompleted + 1)
		tileCopy.Add(tile1, tile2, tile3)
		return meldsList
	}

	for tile := range tileSet {
		// Shuntsu
		if tile.IsNumberTile() {
			count1 := tileCopy.Count(tile)
			count2 := tileCopy.Count(tile+1)
			count3 := tileCopy.Count(tile+2)

			if count1 >= 1 && count2 >= 1 && count3 >= 1 {
				currentMeld := Meld{
					Type:      Shuntsu,
					FirstTile: tile,
				}

				meldsList := removeAndFindMelds(tile, tile+1, tile+2)
				for _, melds := range meldsList {
					melds = append(melds, currentMeld)
					result = append(result, melds)
				}
			}
		}

		// Koutsu
		count := tileCopy.Count(tile)
		if count >= 3 {
			currentMeld := Meld{
				Type:      Koutsu,
				FirstTile: tile,
			}

			meldsList := removeAndFindMelds(tile, tile, tile)
			for _, melds := range meldsList {
				melds = append(melds, currentMeld)
				result = append(result, melds)
			}
		}
		
	}

	return result
}

// Finds all possible pairs in a set of tiles
func findPairs(tileSet MultiSet[Tile]) (ret []Tile) {
	for tile, count := range tileSet {
		if count >= 2 {
			ret = append(ret, tile)
		}
	}

	return ret
}

