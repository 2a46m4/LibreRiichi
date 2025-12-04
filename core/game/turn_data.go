package game

import (
    . "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type TurnData struct {
    // Game index of the current dealer, i.e. the East player
    CurrentDealer uint8
    // CurrentPlayer is the game index of the currently active player
    // that is either discarding or has just discarded.
    CurrentPlayer uint8
    // TotalTurns is the number of draws that have elapsed since
    // the start
    TotalTurns uint8

    // The current round wind
    RoundWind   Wind
}

func InitTurnData() TurnData {
    return TurnData{
    	CurrentDealer: 0,
    	CurrentPlayer: 3, // For the first draw
    	TotalTurns:  0,
    	RoundWind:   East,
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
}

// Increments to the next player and returns it
func (turn *TurnData) NextPlayer() uint8 {
    turn.CurrentPlayer = (turn.CurrentPlayer + 1) % 4
    turn.TotalTurns++
    return turn.CurrentPlayer
}

