package core

import (
	"slices"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type OpenMeldType uint8

const (
	ANKAN_MELD OpenMeldType = iota
	DAIMINKAN_MELD
	SHOUMINKAN_MELD
	PON_MELD
	CHII_MELD
)

type OpenMeld struct {
	Type OpenMeldType
	FirstTile Tile
}

func (meld OpenMeld) Eq(other OpenMeld) bool {
	return meld.FirstTile == other.FirstTile && meld.Type == other.Type
}

type OpenMelds struct {
	melds [4]OpenMeld
	count uint8
}

func (call *OpenMelds) Add(meld OpenMeld) {
	call.melds[call.count] = meld
	call.count += 1
}

func (call OpenMelds) IsEmpty() bool {
	return call.count == 0
}

func (call *OpenMelds) Remove(melds ...OpenMeld) {
	for _, meld := range melds {
		found := false
		for idx, ourMelds := range call.melds {
			if meld.Eq(ourMelds) {
				Swap(call.melds[:], uint(idx), uint(call.count-1))
				call.count -= 1
				found = true
				break
			}
		}

		if !found {
			panic("Couldn't find tile")
		}
	}
}

func (call OpenMelds) Has(meld OpenMeld) bool {
	return slices.ContainsFunc(call.melds[:call.count], meld.Eq)
}

