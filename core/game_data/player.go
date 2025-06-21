package core

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
	"errors"
	"slices"
)

// TODO: Separate Kans into the different types

// ==================== TYPES ====================

type Player struct {
	Hand
	Discards []Tile // For furiten

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

func (player Player) idxOfTile(tile Tile) (int, error) {
	for idx, handTile := range player.ClosedHand {
		if tile == handTile {
			return idx, nil
		}
	}
	return 0, errors.New("Not found")
}

// Finds the tiles that results in a winning hand. Returns an empty list if the hand is not in Tenpai.
func (player Player) checkWaitingTiles() []Tile {
	// TODO: Memoize it later?
	return nil
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
	player.Discards = nil

	player.HandOpen = false
}

// This function essentially keeps track of the player turn. If it's
// the player's turn, there should be an extra tile in the hand.
func (player Player) ExtraTileInHand() bool {
	numOpenTriplets := 0
	numOpenTriplets += len(player.Kans)
	numOpenTriplets += len(player.Chiis)
	numOpenTriplets += len(player.Pons)
	return (len(player.ClosedHand) + numOpenTriplets*3) == 14
}

// Alias
func (player Player) IsPlayerTurn() bool {
	return player.ExtraTileInHand()
}

func (player *Player) Draw(drawn Tile) error {
	if player.ExtraTileInHand() {
		return TooManyTilesErr{}
	}

	player.ClosedHand = append(player.ClosedHand, drawn)
	return nil
}

func (player *Player) Toss(discarded Tile) error {
	if !player.ExtraTileInHand() {
		return TooLittleTilesErr{}
	}

	for i := range player.ClosedHand {
		if player.ClosedHand[i] == discarded {
			player.ClosedHand[i] = Last(player.ClosedHand)
			Pop(&player.ClosedHand)
			player.Discards = append(player.Discards, discarded)
			return nil
		}
	}
	return errors.New("Tile not found")
}

func (player Player) TestChii(tossedTile Tile, tilesInHand [2]Tile) error {
	if player.ExtraTileInHand() {
		return TooManyTilesErr{}
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
		last := Pop(&player.ClosedHand)
		if index >= len(player.ClosedHand) {
			continue
		} else {
			player.ClosedHand[index] = last
		}
	}

	player.Chiis = append(player.Chiis, min(onTile, chiiSequence[0], chiiSequence[1]))
	player.HandOpen = true

	return nil
}

func (player Player) TestAnkan(onTile Tile) error {
	if !player.ExtraTileInHand() {
		return TooLittleTilesErr{}
	}

	if player.countNumInClosedHand(onTile) == 4 {
		return nil
	}
	return errors.New("Not enough tiles to kan")
}

func (player *Player) Ankan(onTile Tile) error {
	if err := player.TestAnkan(onTile); err != nil {
		return err
	}
	for range 4 {
		tileIdx, err := player.idxOfTile(onTile)
		if err != nil {
			panic(err)
		}

		Swap(player.ClosedHand, uint(tileIdx), uint(len(player.ClosedHand)-1))
		Pop(&player.ClosedHand)
	}

	player.Kans = append(player.Kans, onTile)
	return nil
}

func (player Player) TestDaiminkan(onTile Tile) error {
	if player.ExtraTileInHand() {
		return TooManyTilesErr{}
	}

	if player.countNumInClosedHand(onTile) == 3 {
		return nil
	}
	return errors.New("Not enough tiles")
}

func (player *Player) Daiminkan(onTile Tile) error {
	if err := player.TestDaiminkan(onTile); err != nil {
		return err
	}
	for range 3 {
		tileIdx, err := player.idxOfTile(onTile)
		if err != nil {
			panic(err)
		}

		Swap(player.ClosedHand, uint(tileIdx), uint(len(player.ClosedHand)-1))
		Pop(&player.ClosedHand)
	}

	player.Kans = append(player.Kans, onTile)
	player.HandOpen = true
	return nil
}

func (player *Player) TestShouminkan(onTile Tile) error {
	if player.ExtraTileInHand() {
		return TooManyTilesErr{}
	}

	if slices.Contains(player.Pons, onTile) {
		return nil
	}
	return errors.New("Does not have a pon")
}

func (player *Player) Shouminkan(onTile Tile) error {
	if err := player.TestShouminkan(onTile); err != nil {
		return err
	}
	for idx, tile := range player.Pons {
		if tile == onTile {
			Swap(player.Pons, uint(idx), uint(len(player.Pons)-1))
			Pop(&player.Pons)
		}
	}
	player.Kans = append(player.Kans, onTile)
	player.HandOpen = true
	return nil
}

func (player Player) TestPon(onTile Tile) error {
	if player.ExtraTileInHand() {
		return TooManyTilesErr{}
	}

	if player.countNumInClosedHand(onTile) < 2 {
		return errors.New("Cannot pon: Not enough tiles")
	}
	return nil
}

func (player *Player) Pon(onTile Tile) error {
	if err := player.TestPon(onTile); err != nil {
		return err
	}
	for range 2 {
		tileIdx, err := player.idxOfTile(onTile)
		if err != nil {
			panic(err)
		}

		Swap(player.ClosedHand, uint(tileIdx), uint(len(player.ClosedHand)-1))
		Pop(&player.ClosedHand)
	}

	player.Pons = append(player.Pons, onTile)
	player.HandOpen = true
	return nil
}

func (player Player) TestRon(onTile Tile) error {
	// The player needs to have a hand with the correct tile, and waits cannot be in the discard pile
	if player.ExtraTileInHand() {
		return TooManyTilesErr{}
	}

	waitingTiles := player.checkWaitingTiles()
	if idx := slices.Index(waitingTiles, onTile); idx == -1 {
		return errors.New("Tile is not part of waiting tiles")
	}

	for _, waitingTile := range waitingTiles {
		if slices.Index(player.Discards, waitingTile) != -1 {
			return errors.New("Hand in furiten, cannot discard")
		}
	}

	yakus := GetYaku(player.Hand, onTile)
	if yakus == NO_YAKU {
		return errors.New("No yaku")
	}

	return nil
}

// Returns the game result or an error
func (player *Player) Ron(onTile Tile) (WinResult, error) {

	if err := player.TestRon(onTile); err != nil {
		return WinResult{}, err
	}

	yakus := GetYaku(player.Hand, onTile)
	if yakus == NO_YAKU {
		panic("Logic error")
	}

	return WinResult{
		Yakus:       yakus,
		WinningHand: player.Hand,
		WinningTile: onTile,
		WonByRon:    true,
	}, nil
}

func (player Player) TestTsumo() {

}

func (player Player) Tsumo(tsumoTile Tile) (WinResult, error) {
	return WinResult{}, nil
}

func (player Player) GetRiichiDiscards() []Tile {
	return nil
}

func (player Player) TestRiichi(onTile Tile) error {
	return nil
}

func (player Player) Riichi(onTile Tile) error {
	return nil
}
