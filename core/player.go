package core

import (
	"errors"
	"slices"
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

// ==================== PRIVATE FUNCTIONS ====================

func (player Player) countNumInClosedHand(tile Tile) int {
	count := 0
	for _, handTile := range player.ClosedHand {
		if tile == handTile {
			count++
		}
	}
	return count
}

// ==================== PUBLIC FUNCTIONS ====================

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
			player.ClosedHand[i] = Last(player.ClosedHand)
			_, player.ClosedHand = Pop(player.ClosedHand)
			return nil
		}
	}
	return errors.New("Tile not found")
}

func (player Player) TestChii(tossedTile Tile, tilesInHand [2]Tile) error {
	if player.ExtraTileInHand() {
		return errors.New("Player should not have extra tile in hand")
	}

	tiles := [3]Tile{tossedTile, tilesInHand[0], tilesInHand[1]}
	slices.Sort(tiles[:])
	if tiles[1]-tiles[0] != 1 || tiles[2]-tiles[1] != 1 {
		return errors.New("Tiles are not in a sequence")
	}

	exist := 0
	for _, seq := range tilesInHand {
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

	err := player.TestChii(onTile, chiiSequence)
	if err != nil {
		return err
	}

	indices := make([]int, 0, 2)
	for _, seq := range chiiSequence {
		for i, tile := range player.ClosedHand {
			if tile == seq {
				indices = append(indices, i)
			}
		}
	}
	if len(indices) != 2 {
		return errors.New("Does not have all the tiles")
	}

	// Pop the larger index first
	if indices[1] > indices[0] {
		Swap(indices, 0, 1)
	}
	for _, index := range indices {
		var last Tile
		last, player.ClosedHand = Pop(player.ClosedHand)
		if index >= len(player.ClosedHand) {
			continue
		} else {
			player.ClosedHand[index] = last
		}
	}

	player.Chiis = append(player.Chiis, min(onTile, chiiSequence[0], chiiSequence[1]))

	return nil
}

func (player Player) TestAnkan(onTile Tile) error {
	if player.countNumInClosedHand(onTile) == 4 {
		return nil
	} 
	return errors.New("Not enough tiles")
}

func (player *Player) Ankan(onTile Tile) error {
	return nil
}

func (player Player) TestDaiminkan(onTile Tile) error {
	if player.countNumInClosedHand(onTile) == 3 {
		return nil
	} 
	return errors.New("Not enough tiles")
}

func (player *Player) Daiminkan(onTile Tile) error {
	return nil
}

func (player *Player) TestShouminkan(onTile Tile) error {
	return nil
}

func (player *Player) Shouminkan(onTile Tile) error {
	return nil
}

func (player Player) TestPon(onTile Tile) error {
	return nil
}

func (player *Player) Pon(onTile Tile) error {
	return nil
}

func (player Player) TestRon(onTile Tile) error {
	return nil
}

func (player Player) Ron(onTile Tile) error {
	return nil
}

func (player Player) Riichi(onTile Tile) error {
	return nil
}
