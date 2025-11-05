package core

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type Hand struct {
	ClosedHand ClosedHand
	Kans       KanCalls
	Pons       OpenCalls
	Chiis      OpenCalls // Chiis are the start of the sequence
	InRiichi   bool
}

func (hand Hand) Open() bool {
	return hand.Chiis.IsEmpty() && hand.Kans.IsEmpty() && hand.Pons.IsEmpty()
}

func (hand Hand) Draw(tile Tile) {
	hand.ClosedHand.Add(tile)
}

func (hand Hand) Discard(tile Tile) {
	hand.ClosedHand.RemoveTile(tile)
}

func (hand Hand) FullHand() bool {
	numOpenTriplets := uint8(0)
	numOpenTriplets += hand.Kans.count
	numOpenTriplets += hand.Chiis.count
	numOpenTriplets += hand.Pons.count
	// 13 since ClosedHand.index is 0-based
	return (hand.ClosedHand.index + numOpenTriplets*3) == 13
}


func (hand *Hand) Pon(tile, discard Tile) {
	hand.ClosedHand.RemoveTile(tile, tile, discard)
	hand.Pons.Add(tile)
}

func (hand *Hand) DaiminKan(tile, draw Tile) {
	hand.ClosedHand.RemoveTile(tile, tile, tile)
	hand.ClosedHand.Add(draw)
	hand.Kans.Add(tile, DAIMINKAN)
}

func (hand *Hand) ShouminKan(tile, draw Tile) {
	hand.Pons.Remove(tile)
	hand.ClosedHand.Add(draw)
	hand.Kans.Add(tile, SHOUMINKAN)
}

func (hand *Hand) AnKan(tile, draw Tile) {
	hand.ClosedHand.RemoveTile(tile, tile, tile, tile)
	hand.ClosedHand.Add(draw)
	hand.Kans.Add(tile, ANKAN)
}

func (hand *Hand) Chii(firstTile Tile, tiles [2]Tile, discard Tile) {
	hand.ClosedHand.RemoveTile(tiles[:]...)
	hand.Chiis.Add(firstTile)
}

func (hand *Hand) Riichi(tile Tile) {
	hand.InRiichi = true
	hand.ClosedHand.RemoveTile(tile)
}

func (hand *Hand) Ron(tile Tile) {
	hand.ClosedHand.Add(tile)
}

func (hand *Hand) Tsumo(tile Tile) {
	hand.ClosedHand.Add(tile)
}
