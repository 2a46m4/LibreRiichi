package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type ActionValidator struct {
	*TileState
}

func NewActionValidator(game *TileState) ActionValidator {
	return ActionValidator{game}
}

func (av *ActionValidator) ValidateDraw(playerIdx uint8, drawnTile Tile) error {
	

	return nil
}

func (av *ActionValidator) ValidateDiscard(playerIdx uint8, tileToToss Tile) error {
	// TODO: Validate if player can discard tile
	return true
}

func (av *ActionValidator) ValidateRiichi(playerIdx uint8, tileToRiichi Tile) error {
	// TODO: Validate if player can declare riichi
	return true
}

func (av *ActionValidator) ValidateTsumo(playerIdx uint8, tileToTsumo Tile) error {
	// TODO: Validate if player can win by tsumo
	return true
}

func (av *ActionValidator) ValidateRon(playerIdx uint8, tileToRon Tile) error {
	// TODO: Validate if player can win by ron
	return true
}

func (av *ActionValidator) ValidatePon(playerIdx uint8, tileToPon Tile, againstPlayer uint8) error {
	// TODO: Validate if player can call pon
	return true
}

func (av *ActionValidator) ValidateKan(playerIdx uint8, tileToKan Tile, againstPlayer uint8) error {
	// TODO: Validate if player can call kan (could be daiminkan or shouminkan)
	return true
}

func (av *ActionValidator) ValidateChii(playerIdx uint8, tileToChii Tile, tilesInHand [2]Tile) error {
	// TODO: Validate if player can call chii
	return true
}

func (av *ActionValidator) ValidateSkip(playerIdx uint8, actionToSkip Action) error {
	// TODO: Validate if player can skip action
	return true
}

func (av *ActionValidator) ValidateAction(action Action, playerIdx uint8, against uint8) error {
	switch act := action.(type) {
	case *Draw:
		return av.ValidateDraw(playerIdx, act.DrawnTile)
	case *Toss:
		return av.ValidateDiscard(playerIdx, act.TileToToss)
	case *Riichi:
		return av.ValidateRiichi(playerIdx, act.TileToRiichi)
	case *Tsumo:
		return av.ValidateTsumo(playerIdx, act.TileToTsumo)
	case *Ron:
		return av.ValidateRon(playerIdx, act.TileToRon)
	case *Pon:
		return av.ValidatePon(playerIdx, act.TileToPon, against) // TODO: determine againstPlayer
	case *Kan:
		return av.ValidateKan(playerIdx, act.TileToKan, against) // TODO: determine againstPlayer
	case *Chii:
		return av.ValidateChii(playerIdx, act.TileToChii, act.TilesInHand)
	case *Skip:
		return av.ValidateSkip(playerIdx, act.ActionToSkip)
	default:
		return false
	}
}


