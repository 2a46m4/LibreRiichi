package core

import (
	"errors"

	core "codeberg.org/ijnakashiar/LibreRiichi/core"
)

type MahjongGame struct {
	Players     []Player
	PlayerOrder []uint8

	Tiles         [136]Tile
	LiveWallIndex int
	DeadWall      []Tile
}

func (game *MahjongGame) setupGame() {
	// Randomize the player index
	core.PermuteArray[uint8](game.PlayerOrder)

	for idx := range game.PlayerOrder {
		player := game.Players[idx]

	}
}

func (MahjongGame) JoinArena(PlayerIdx int) error {
	return nil
}

// Returns nil when a new game can be started, otherwise an error
func (game MahjongGame) StartNewGame() ([][]Setup, error) {
	if len(game.Players) != 4 {
		return nil, errors.New("Not enough players")
	}

	game.setupGame()

	return nil, nil
}

// Returns the updated game state and things to notify when a player action is taken
func RespondToAction(action PlayerAction) ([]ActionResult, bool) {
	return nil, false
}

// Return the game results
func GetGameResults() (GameResult, error) {
	return GameResult{}, nil
}

func (MahjongGame) GetMaxPlayers() int {
	return 4
}
