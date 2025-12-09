package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type TurnData struct {
	// Game index of the current dealer, i.e. the East player
	// This always starts at 0.
	CurrentDealer uint8
	// CurrentPlayer is the game index of the currently active player
	// that is either discarding or has just discarded.
	CurrentPlayer uint8
	// TotalTurns is the number of draws that have elapsed since
	// the start
	TotalTurns uint8

	// How many turns it has been since the player riichi'd. If the
	// player hasn't riichi'd, then this value is -1.
	TurnsSinceRiichi [4]int8

	// The current round wind
	RoundWind Wind
}

func InitTurnData() TurnData {
	return TurnData{
		CurrentDealer:    0,
		CurrentPlayer:    3, // For the first draw
		TurnsSinceRiichi: [4]int8{-1, -1, -1, -1},
		TotalTurns:       0,
		RoundWind:        East,
	}
}

// Returns the 1-indexed round number of the current round wind
func (turn TurnData) GetRoundNumber() uint8 {
	return turn.CurrentDealer + 1
}

func (turn *TurnData) NextRound(newDealer bool) {
	turn.TotalTurns = 0
	if newDealer {
		turn.CurrentDealer = (turn.CurrentDealer + 1) % 4
	}
	if turn.CurrentDealer == 0 {
		turn.RoundWind += 1
	}
	turn.CurrentPlayer = turn.CurrentDealer
	turn.TurnsSinceRiichi = [4]int8{-1, -1. - 1, -1}
}

func (turn *TurnData) TrackRiichi(playerIdx uint8) {
	turn.TurnsSinceRiichi[playerIdx] = 0
}

// Increments to the next player and returns it
func (turn *TurnData) NextPlayer() uint8 {
	turn.CurrentPlayer = (turn.CurrentPlayer + 1) % 4
	turn.TotalTurns++
	for i := range turn.TurnsSinceRiichi {
		if turn.TurnsSinceRiichi[i] != -1 {
			turn.TurnsSinceRiichi[i]++
		}
	}
	return turn.CurrentPlayer
}

func (turn *TurnData) IppatsuPossible(idx uint8) bool {
	return turn.TurnsSinceRiichi[idx] >= 0 && turn.TurnsSinceRiichi[idx] <= 4
}

func (turn *TurnData) IsDoubleRiichi(idx uint8) bool {
	return turn.TurnsSinceRiichi[idx] == int8(turn.TotalTurns)
}

func (turn *TurnData) GetPlayerWind(idx uint8) Wind {
	return East + (Wind(4+idx-turn.CurrentDealer) % 4)
}
