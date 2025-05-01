package core

type MahjongGame struct {
	LiveWallIndex int
	Players       []Player

	Tiles [136]Tile
}

func (MahjongGame) JoinArena(PlayerIdx int) error {
	return nil
}

func StartNewGame() {

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
