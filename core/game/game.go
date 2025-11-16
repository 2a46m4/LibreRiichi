package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
)

type MahjongGame struct {
	scoring Scoring
	gameState    GameState
	mahjongRound MahjongRound
	ordering     Ordering

	roundWind   Wind
	roundNumber uint8

	firstRound bool
}

func NewMahjongGame() *MahjongGame {
	return &MahjongGame{
		gameState: *InitGameState(),
		firstRound: true,
	}
}

func (game *MahjongGame) handleNewTurn(playerIdx uint8) ([]MessageSendInfo, error) {
	err := game.gameState.Transition("draw-tile")
	if err != nil {
		return nil, err
	}
	err = game.turnState.Try(Draw{}, playerIdx)
	if err != nil {
		return nil, err
	}

	draw := game.TileState.Draw(playerIdx)
	_ = draw // Use the variable to avoid "declared and not used" error

	// Draw, checking that we still have moves
	// Check riichi, tsumo, kan
	return nil, nil
}

func (game *MahjongGame) StartGame() (messages []MessageSendInfo, err error) {
	err = game.gameState.Transition("start-game", game)
	sendInfo, _ := game.gameState.GetReturn()
	return sendInfo.([]MessageSendInfo), err
}

func (game *MahjongGame) StartRound() (msgs []MessageSendInfo, err error) {
	err = game.gameState.Transition("start-round", &game.mahjongRound, game.firstRound)
	sendInfo, _ := game.gameState.GetReturn()
	return sendInfo.([]MessageSendInfo), err
}

func (game *MahjongGame) HandleEvent(action Action, arenaIdx uint8) (msgs []MessageSendInfo, err error) {
	gameIdx := game.ordering.ArenaToGame[arenaIdx]
	err = game.gameState.Transition("handle-event", game.mahjongRound, game.ordering, action, arenaIdx)

	if err != nil {
		return nil, err
	}

	return nextTurnInfo, nil
}

func (game *MahjongGame) RoundEnd() error {
	// Do the increment post-round
	game.mahjongRound.data.IncrementRound()
		
	if game.firstRound {
		game.firstRound = false
	}

	return nil
}

func (game *MahjongGame) RoundEnded() bool {
	return game.gameState.Current() == "round-end"
}

func (game *MahjongGame) GameEnd() error {
	return nil
}

func (game *MahjongGame) GameEnded() bool {
	return game.gameState.Current() == "game-end"
}
