package core

import (
	"errors"
	"slices"
	"sort"

	core "codeberg.org/ijnakashiar/LibreRiichi/core"
)

type Player struct {
	ClosedHand []Tile
	Kans       []Tile
	Pons       []Tile
	Chiis      []Tile // Chiis are the start of the sequence

	InRiichi bool
	HandOpen bool

	Points   uint32
	SeatWind Wind
}

func (player *Player) FreshHand(tiles []Tile) {
	if len(tiles) != 13 {
		panic("Not equal to 13")
	}

	player.ClosedHand = make([]Tile, 14)
	copy(player.ClosedHand, tiles)
	player.Kans = nil
	player.Pons = nil
	player.Chiis = nil

	player.InRiichi = false
	player.HandOpen = false
}

func (player *Player) ExtraTileInHand() bool {
	numOpenTriplets := 0
	numOpenTriplets += len(player.Kans)
	numOpenTriplets += len(player.Chiis)
	numOpenTriplets += len(player.Pons)
	return (len(player.ClosedHand) + numOpenTriplets*3) == 14
}

func (player *Player) Draw(drawn Tile) error {
	if player.ExtraTileInHand() {
		return errors.New("Too many tiles in hand")
	}

	player.ClosedHand = append(player.ClosedHand, drawn)
	return nil
}

func (player *Player) Toss(discarded Tile) error {
	if !player.ExtraTileInHand() {
		return errors.New("Too little tiles in hand")
	}

	for i := range player.ClosedHand {
		if player.ClosedHand[i] == discarded {
			player.ClosedHand[i] = core.Last(player.ClosedHand)
			_, player.ClosedHand = core.Pop(player.ClosedHand)
			return nil
		}
	}
	return errors.New("Tile not found")
}

func (player Player) GetPostTurnMoves(tossedTile Tile) []PlayerAction {
	// Iterate through all possible combinations of Chii
	tileNum := tossedTile.GetTileNumber()
	moves := make([]PlayerAction, 0)
	if tileNum < 7 { // 6, 7, 8
		tiles := [3]Tile{tossedTile, tossedTile + 1, tossedTile + 2}
		if player.TestChii(tiles) != nil {
			moves = append(moves, PlayerAction(
				ChiiAction{
					Action:          CHII,
					FromPlayer:      0,
					PotentialAction: true,
					Data: map[string]any{
						"tiles": tiles,
					},
				},
			))
		}
	}
	if tileNum >= 0 {

	}

	return nil
}

func (player Player) TestChii(tiles [3]Tile) error {
	if player.ExtraTileInHand() {
		return errors.New("Player should not have extra tile in hand")
	}

	if tiles[1]-tiles[0] != 1 || tiles[2]-tiles[1] != 1 {
		return errors.New("Tiles are not in a sequence")
	}

	exist := 0
	for _, seq := range tiles {
		for _, tile := range player.ClosedHand {
			if tile == seq {
				exist += 1
			}
		}
	}
	if exist != 3 {
		return errors.New("Non suitable tiles")
	}

	return nil
}

func (player *Player) Chii(onTile Tile, chiiSequence [2]Tile) error {
	tiles := [3]Tile{
		onTile, chiiSequence[0], chiiSequence[1],
	}
	slices.Sort(tiles[:])

	err := player.TestChii(tiles)
	if err != nil {
		return err
	}

	indices := make([]int, 0, 3)
	for _, seq := range tiles {
		for i, tile := range player.ClosedHand {
			if tile == seq {
				indices = append(indices, i)
			}
		}
	}
	if len(indices) != 3 {
		return errors.New("Does not have all the sequences")
	}

	sort.Sort(sort.Reverse(sort.IntSlice(indices)))
	for _, index := range indices {
		var last Tile
		last, player.ClosedHand = core.Pop(player.ClosedHand)
		if index >= len(player.ClosedHand) {
			continue
		} else {
			player.ClosedHand[index] = last
		}
	}

	player.Chiis = append(player.Chiis, tiles[0])

	return nil
}

func (player *Player) Pon(onTile Tile) error {

	return nil
}

func (player Player) Ron(onTile Tile) {

}

func (player Player) Riichi(onTile Tile) {

}
