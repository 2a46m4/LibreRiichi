package core

type MahjongGame struct {
	Players []Player

	Tiles         [136]Tile
	LiveWallIndex int
	DeadWall      []Tile
}

func (MahjongGame) JoinArena(PlayerIdx int) error {
	return nil
}

// Returns nil when a new game can be started.
func (MahjongGame) StartNewGame() ([][]Setup, error) {
	return [][]Setup{}, nil
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
