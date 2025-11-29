package hand

import (
	"errors"
	"fmt"
	"slices"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
)

type Hand struct {
	ClosedHand ClosedHand
	OpenMelds  OpenMelds
	InRiichi   bool
	WaitingFor []Tile // Only used when the hand is in riichi
}

func (hand Hand) Open() bool {
	return !hand.Closed()
}

func (hand Hand) Closed() bool {
	return hand.OpenMelds.IsEmpty()
}

func (hand Hand) InTenpai() bool {
	return false
}

func (hand Hand) FullHand() bool {
	// 13 since ClosedHand.index is 0-based
	return (hand.ClosedHand.index + hand.OpenMelds.count*3) == 13
}

// Returns the tile that the player just received
func (hand Hand) TileJustReceived() (tile Tile, err error) {
	if hand.FullHand() {
		return hand.ClosedHand.hand[hand.ClosedHand.index-1], nil
	} else {
		return tile, errors.New(fmt.Sprintf("Not full hand, don't have an extra tile: %#v %#v", hand))
	}
}

func (hand Hand) OpenMeldCount() uint8 {
	return hand.OpenMelds.count
}

func (hand Hand) TestDraw() bool {
	return !hand.FullHand()
}

func (hand *Hand) Draw(tile Tile) {
	hand.ClosedHand.Add(tile)
}

func (hand Hand) TestDiscard(tile Tile) bool {
	if !hand.FullHand() {
		return false
	}

	if hand.InRiichi {
		handTile, _ := hand.ClosedHand.Last()
		return handTile == tile
	} else {
		return hand.ClosedHand.HasTile(tile)
	}

}

func (hand *Hand) Discard(tile Tile) {
	hand.ClosedHand.RemoveTile(tile)
}

func (hand Hand) TestPon(tile Tile) bool {
	return !hand.InRiichi && !hand.FullHand() && (hand.ClosedHand.HasTileN(tile) == 2)
}

func (hand *Hand) Pon(tile Tile) {
	hand.ClosedHand.RemoveTile(tile, tile)
	hand.OpenMelds.Add(OpenMeld{
		Type:      PON_MELD,
		FirstTile: tile,
	})
}

func (hand Hand) TestDaiminKan(tile, draw Tile) bool {
	return !hand.InRiichi && !hand.FullHand() && (hand.ClosedHand.HasTileN(tile) == 3)
}

func (hand *Hand) DaiminKan(tile Tile) {
	hand.ClosedHand.RemoveTile(tile, tile, tile)
	hand.OpenMelds.Add(OpenMeld{
		Type:      DAIMINKAN_MELD,
		FirstTile: tile,
	})
}

func (hand Hand) TestShouminkan(tile Tile) bool {
	return !hand.InRiichi && !hand.FullHand() && hand.OpenMelds.Has(OpenMeld{
		Type:      PON_MELD,
		FirstTile: tile,
	})
}

func (hand *Hand) ShouminKan(tile Tile) {
	hand.OpenMelds.Remove(OpenMeld{
		Type:      PON_MELD,
		FirstTile: tile,
	})
	hand.OpenMelds.Add(OpenMeld{
		Type:      SHOUMINKAN_MELD,
		FirstTile: tile,
	})
}

func (hand Hand) TestAnKan(tile Tile) bool {
	return !hand.InRiichi && hand.FullHand() && (hand.ClosedHand.HasTileN(tile) == 4)
}

func (hand *Hand) AnKan(tile Tile) {
	hand.ClosedHand.RemoveTile(tile, tile, tile, tile)
	hand.OpenMelds.Add(OpenMeld{
		Type:      ANKAN_MELD,
		FirstTile: tile,
	})
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
	hand.OpenMelds.Add(OpenMeld{
		Type:      CHII_MELD,
		FirstTile: firstTile,
	})
}

// Also need to test that the hand has a yaku
func (hand Hand) TestRiichi(tile Tile) bool {
	return !hand.InRiichi && hand.FullHand() && hand.ClosedHand.HasTile(tile)
}

func (hand *Hand) Riichi(tile Tile) {
	// Calculate the tiles that the user is waiting for

	hand.InRiichi = true
	hand.ClosedHand.RemoveTile(tile)
}

// TODO: Also need to test that the hand has a yaku
func (hand Hand) TestRon(tile Tile) bool {
	return !hand.FullHand() && slices.Contains(hand.WaitingFor, tile)
}

func (hand *Hand) Ron(tile Tile) {
	hand.ClosedHand.Add(tile)
}

// Also need to test that the hand has a yaku
func (hand Hand) TestTsumo(tile Tile) bool {
	return hand.FullHand() && slices.Contains(hand.WaitingFor, tile)
}

func (hand *Hand) Tsumo(tile Tile) {
	hand.ClosedHand.Add(tile)
}
