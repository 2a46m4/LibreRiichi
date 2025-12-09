package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
)

type LiveWall struct {
	tiles [70]Tile
	// Index represents
	index uint8
}

func (wall *LiveWall) GetLiveTile() (tile Tile) {
	tile = wall.tiles[wall.index]
	wall.index += 1
	return tile
}

func (wall LiveWall) End() bool {
	return wall.index == 70
}
