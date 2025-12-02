package game

import (
	"errors"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type TurnData struct {
	// Game index of the current dealer
	CurrentDealer uint8
	// TurnNumber is the game index of the currently active player
	// that is either discarding or has just discarded.
	TurnNumber uint8
	// TotalTurns is the number of draws that have elapsed since
	// the start
	TotalTurns uint8
}

func InitTurnData() TurnData {
	return TurnData{
		CurrentDealer: 0,
		TurnNumber:    3, // For the first draw
		TotalTurns:    0,
	}
}

func (turn *TurnData) NextRound() {
	turn.TotalTurns = 0
	turn.TurnNumber = (turn.CurrentDealer + 1) % 4
	turn.CurrentDealer = (turn.CurrentDealer + 1) % 4
}

func (turn *TurnData) PlayerDraw() uint8 {
	expectedPlayer := (turn.TurnNumber + 1) % 4
	turn.TurnNumber = expectedPlayer
	turn.TotalTurns++
	return expectedPlayer
}

func (turn *TurnData) PlayerPon(playerIdx uint8) error {
	if playerIdx == turn.TurnNumber {
		return errors.New("invalid pon: cannot call pon on own discard")
	}

	turn.TurnNumber = playerIdx
	return nil
}

func (turn *TurnData) PlayerKan(playerIdx uint8) error {
	if playerIdx == turn.TurnNumber {
		return errors.New("invalid kan: cannot call claimed kan on own discard")
	}

	turn.TurnNumber = playerIdx
	return nil
}

func (turn *TurnData) PlayerClosedKan(playerIdx uint8) error {
	if playerIdx != turn.TurnNumber {
		return errors.New("invalid closed kan: only current player can call closed kan")
	}
	return nil
}

func (turn *TurnData) PlayerChii(playerIdx uint8) error {
	expectedPlayer := (turn.TurnNumber + 1) % 4
	if playerIdx != expectedPlayer {
		return errors.New("invalid chii: only the next player can call chii")
	}
	turn.TurnNumber = playerIdx
	return nil
}

func (turn *TurnData) GetCurrentPlayer() uint8 {
	return turn.TurnNumber
}

func (turn *TurnData) GetTotalTurns() uint8 {
	return turn.TotalTurns
}

// ProcessKanAction handles kan actions with type distinction
// isClosedKan should be true for closed kans (from hand only), false for claimed kans
func (turn *TurnData) ProcessKanAction(action Kan, playerIdx uint8, isClosedKan bool) error {
	if isClosedKan {
		return turn.PlayerClosedKan(playerIdx)
	}
	return turn.PlayerKan(playerIdx)
}
