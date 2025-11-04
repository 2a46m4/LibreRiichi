package game

import (
	"testing"
)

func TestTurnState(t *testing.T) {
	turn := TurnState{
		TurnNumber: 0,
		TotalTurns: 0,
	}

	// Test initial state
	if turn.GetCurrentPlayer() != 0 {
		t.Errorf("Expected initial player to be 0, got %d", turn.GetCurrentPlayer())
	}
	if turn.GetTotalTurns() != 0 {
		t.Errorf("Expected initial total turns to be 0, got %d", turn.GetTotalTurns())
	}

	// Test PlayerDraw
	err := turn.PlayerDraw(0)
	if err != nil {
		t.Errorf("PlayerDraw returned error: %v", err)
	}
	if turn.GetCurrentPlayer() != 0 {
		t.Errorf("Expected current player to be 0 after draw, got %d", turn.GetCurrentPlayer())
	}
	if turn.GetTotalTurns() != 1 {
		t.Errorf("Expected total turns to be 1 after draw, got %d", turn.GetTotalTurns())
	}

	// Test NextTurn
	turn.NextTurn()
	if turn.GetCurrentPlayer() != 1 {
		t.Errorf("Expected current player to be 1 after NextTurn, got %d", turn.GetCurrentPlayer())
	}
	if turn.GetTotalTurns() != 1 {
		t.Errorf("Expected total turns to remain 1 after NextTurn, got %d", turn.GetTotalTurns())
	}

	// Test PlayerDraw on different player
	err = turn.PlayerDraw(1)
	if err != nil {
		t.Errorf("PlayerDraw returned error: %v", err)
	}
	if turn.GetCurrentPlayer() != 1 {
		t.Errorf("Expected current player to be 1 after draw, got %d", turn.GetCurrentPlayer())
	}
	if turn.GetTotalTurns() != 2 {
		t.Errorf("Expected total turns to be 2 after second draw, got %d", turn.GetTotalTurns())
	}

	// Test PlayerPon
	err = turn.PlayerPon(3)
	if err != nil {
		t.Errorf("PlayerPon returned error: %v", err)
	}
	if turn.GetCurrentPlayer() != 3 {
		t.Errorf("Expected current player to be 3 after pon, got %d", turn.GetCurrentPlayer())
	}
	if turn.GetTotalTurns() != 2 {
		t.Errorf("Expected total turns to remain 2 after pon, got %d", turn.GetTotalTurns())
	}

	// Test PlayerKan
	err = turn.PlayerKan(2)
	if err != nil {
		t.Errorf("PlayerKan returned error: %v", err)
	}
	if turn.GetCurrentPlayer() != 2 {
		t.Errorf("Expected current player to be 2 after kan, got %d", turn.GetCurrentPlayer())
	}
	if turn.GetTotalTurns() != 2 {
		t.Errorf("Expected total turns to remain 2 after kan, got %d", turn.GetTotalTurns())
	}

	// Test PlayerChii
	err = turn.PlayerChii(1)
	if err != nil {
		t.Errorf("PlayerChii returned error: %v", err)
	}
	if turn.GetCurrentPlayer() != 1 {
		t.Errorf("Expected current player to be 1 after chii, got %d", turn.GetCurrentPlayer())
	}
	if turn.GetTotalTurns() != 2 {
		t.Errorf("Expected total turns to remain 2 after chii, got %d", turn.GetTotalTurns())
	}

	// Test wrap-around with NextTurn
	turn.TurnNumber = 3
	turn.NextTurn()
	if turn.GetCurrentPlayer() != 0 {
		t.Errorf("Expected current player to wrap to 0, got %d", turn.GetCurrentPlayer())
	}

	// Test invalid player index
	err = turn.PlayerDraw(4)
	if err == nil {
		t.Error("Expected error for invalid player index 4, but got none")
	}
	err = turn.PlayerPon(5)
	if err == nil {
		t.Error("Expected error for invalid player index 5, but got none")
	}
}

func TestTurnStateIntegration(t *testing.T) {
	// Test a typical game sequence
	game := &MahjongGame{} // This will have embedded TurnState with zero values

	// Simulate Player 0's turn
	err := game.PlayerDraw(0)
	if err != nil {
		t.Errorf("PlayerDraw failed: %v", err)
	}
	if game.GetCurrentPlayer() != 0 {
		t.Errorf("Expected player 0, got %d", game.GetCurrentPlayer())
	}
	if game.GetTotalTurns() != 1 {
		t.Errorf("Expected 1 total turn, got %d", game.GetTotalTurns())
	}

	// Player 0 calls pon on player 3's discard
	err = game.PlayerPon(0)
	if err != nil {
		t.Errorf("PlayerPon failed: %v", err)
	}
	if game.GetCurrentPlayer() != 0 {
		t.Errorf("Expected player 0 after pon, got %d", game.GetCurrentPlayer())
	}
	if game.GetTotalTurns() != 1 {
		t.Errorf("Expected total turns to remain 1 after pon, got %d", game.GetTotalTurns())
	}
}