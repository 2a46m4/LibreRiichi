package core

import (
	"slices"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type OpenCalls struct {
	tiles [4]Tile
	count uint8
}

func (call *OpenCalls) Add(tile Tile) {
	call.tiles[call.count] = tile
	call.count += 1
}

func (call OpenCalls) IsEmpty() bool {
	return call.count == 0
}

func (call *OpenCalls) Remove(tiles ...Tile) {
	for _, tile := range tiles {
		found := false
		for idx, callTile := range call.tiles {
			if tile == callTile {
				Swap(call.tiles[:], uint(idx), uint(call.count-1))
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

func (call OpenCalls) Has(tile Tile) bool {
	return slices.Contains(call.tiles[:call.count], tile)
}

type KanType uint8

const (
	ANKAN KanType = iota
	DAIMINKAN
	SHOUMINKAN
)

type KanCalls struct {
	OpenCalls
	kanType [4]KanType
}

func (call KanCalls) Add(tile Tile, kanType KanType) {
	call.tiles[call.count] = tile
	call.kanType[call.count] = kanType
	call.count += 1
}
