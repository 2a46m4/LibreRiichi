package game

import (
	"errors"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type TurnState struct {
	// TurnNumber is the game index of the currently active player
	// that is either discarding or has just discarded.
	TurnNumber uint8
	// TotalTurns is the number of draws that have elapsed since
	// the start
	TotalTurns uint8
}

func (turn *TurnState) PlayerDraw(playerIdx uint8) error {
	expectedPlayer := (turn.TurnNumber + 1) % 4
	if playerIdx != expectedPlayer {
		return errors.New("invalid draw: not the next player's turn")
	}

	turn.TurnNumber = playerIdx
	turn.TotalTurns++

	return nil
}

func (turn *TurnState) PlayerPon(playerIdx uint8) error {
	if playerIdx == turn.TurnNumber {
		return errors.New("invalid pon: cannot call pon on own discard")
	}

	turn.TurnNumber = playerIdx
	return nil
}

func (turn *TurnState) PlayerKan(playerIdx uint8) error {
	if playerIdx == turn.TurnNumber {
		return errors.New("invalid kan: cannot call claimed kan on own discard")
	}

	turn.TurnNumber = playerIdx
	return nil
}

func (turn *TurnState) PlayerClosedKan(playerIdx uint8) error {
	if playerIdx != turn.TurnNumber {
		return errors.New("invalid closed kan: only current player can call closed kan")
	}
	return nil
}

func (turn *TurnState) PlayerChii(playerIdx uint8) error {
	expectedPlayer := (turn.TurnNumber + 1) % 4
	if playerIdx != expectedPlayer {
		return errors.New("invalid chii: only the next player can call chii")
	}
	turn.TurnNumber = playerIdx
	return nil
}

func (turn *TurnState) GetCurrentPlayer() uint8 {
	return turn.TurnNumber
}

func (turn *TurnState) GetTotalTurns() uint8 {
	return turn.TotalTurns
}

// ProcessAction dispatches an action to the appropriate turn state method
// Returns an error if the action is invalid for the current turn state
func (turn *TurnState) ProcessAction(action Action, playerIdx uint8) error {
	switch action.(type) {
	case Draw:
		return turn.PlayerDraw(playerIdx)
	case Pon:
		return turn.PlayerPon(playerIdx)
	case Kan:
		return turn.PlayerKan(playerIdx)
	case Chii:
		return turn.PlayerChii(playerIdx)
	default:
		return nil
	}
}

// ProcessKanAction handles kan actions with type distinction
// isClosedKan should be true for closed kans (from hand only), false for claimed kans
func (turn *TurnState) ProcessKanAction(action Kan, playerIdx uint8, isClosedKan bool) error {
	if isClosedKan {
		return turn.PlayerClosedKan(playerIdx)
	}
	return turn.PlayerKan(playerIdx)
}
