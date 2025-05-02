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

func StartNewGame() error {
	return nil
}

func RespondToAction(action PlayerAction) ([]ActionResult, bool) {
	return nil, false
}

func GetGameResults() (GameResult, error) {
	return GameResult{}, nil
}

func (MahjongGame) GetMaxPlayers() int {
	return 4
}
