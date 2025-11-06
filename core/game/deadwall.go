package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type DoraDeadWall struct {
	tiles          [5]Tile
	numberRevealed uint8
}

func (d *DoraDeadWall) Reset(tiles []Tile) {
	for i, t := range tiles {
		d.tiles[i] = t
	}
	d.numberRevealed = 0
}

func (d *DoraDeadWall) revealDora() {
	d.numberRevealed += 1
}

func (d *DoraDeadWall) getLastDoraTile() Tile {
	return d.tiles[d.numberRevealed-1]
}

type UradoraDeadWall DoraDeadWall

type KanDeadWall struct {
	tiles       [4]Tile
	numberDrawn uint8
}

func (d *KanDeadWall) Reset(tiles []Tile) {
	for i, t := range tiles {
		d.tiles[i] = t
	}
	d.numberDrawn = 0
}

func (d *KanDeadWall) drawKan() (tile Tile) {
	tile = d.tiles[d.numberDrawn]
	d.numberDrawn += 1
	return tile
}

type DeadWall struct {
	UradoraDeadWall
	DoraDeadWall
	KanDeadWall
}

func (d *DeadWall) Reset(tiles []Tile) {
	// First 4 tiles go to KanDeadWall
	d.KanDeadWall.Reset(tiles[0:4])
	// Next 5 tiles go to DoraDeadWall
	d.DoraDeadWall.Reset(tiles[4:9])
	// Last 5 tiles go to UradoraDeadWall
	d.UradoraDeadWall.Reset(tiles[9:14])
}
