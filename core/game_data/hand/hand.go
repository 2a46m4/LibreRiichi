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

func (hand Hand) FullHand() bool {
	numOpenTriplets := uint8(0)
	numOpenTriplets += hand.Kans.count
	numOpenTriplets += hand.Chiis.count
	numOpenTriplets += hand.Pons.count
	// 13 since ClosedHand.index is 0-based
	return (hand.ClosedHand.index + numOpenTriplets*3) == 13
}

func (hand Hand) TestDraw() bool {
	return !hand.FullHand()
}

func (hand *Hand) Draw(tile Tile) {
	hand.ClosedHand.Add(tile)
}

func (hand Hand) TestDiscard(tile Tile) bool {
	return hand.FullHand() && hand.ClosedHand.HasTile(tile)
}

func (hand *Hand) Discard(tile Tile) {
	hand.ClosedHand.RemoveTile(tile)
}

func (hand Hand) TestPon(tile Tile) bool {
	return !hand.InRiichi && !hand.FullHand() && (hand.ClosedHand.HasTileN(tile) == 2)
}

func (hand *Hand) Pon(tile Tile) {
	hand.ClosedHand.RemoveTile(tile, tile)
	hand.Pons.Add(tile)
}

func (hand Hand) TestDaiminKan(tile, draw Tile) bool {
	return !hand.InRiichi && !hand.FullHand() && (hand.ClosedHand.HasTileN(tile) == 3)
}

func (hand *Hand) DaiminKan(tile, draw Tile) {
	hand.ClosedHand.RemoveTile(tile, tile, tile)
	hand.ClosedHand.Add(draw)
	hand.Kans.Add(tile, DAIMINKAN)
}

func (hand Hand) TestShouminkan(tile, draw Tile) bool {
	return !hand.InRiichi && !hand.FullHand() && hand.Pons.Has(tile)
}

func (hand *Hand) ShouminKan(tile, draw Tile) {
	hand.Pons.Remove(tile)
	hand.ClosedHand.Add(draw)
	hand.Kans.Add(tile, SHOUMINKAN)
}

func (hand Hand) TestAnKan(tile, draw Tile) bool {
	return !hand.InRiichi && hand.FullHand() && (hand.ClosedHand.HasTileN(tile) == 4)
}

func (hand *Hand) AnKan(tile, draw Tile) {
	hand.ClosedHand.RemoveTile(tile, tile, tile, tile)
	hand.ClosedHand.Add(draw)
	hand.Kans.Add(tile, ANKAN)
}

func (hand Hand) TestChii(firstTile Tile, tiles [2]Tile) bool {
	return !hand.InRiichi && !hand.FullHand() && hand.ClosedHand.HasTile(tiles[0]) && hand.ClosedHand.HasTile(tiles[1])
}

// firstTile is the first tile of the sequence and can be in the
// closed hand or from the discard pile.
// 
// tiles are tiles in the closed hand that need to be removed.
func (hand *Hand) Chii(firstTile Tile, tiles [2]Tile) {
	hand.ClosedHand.RemoveTile(tiles[:]...)
	hand.Chiis.Add(firstTile)
}

func (hand Hand) TestRiichi(tile Tile) {
	
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
