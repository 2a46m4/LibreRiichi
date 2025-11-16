package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
)

type MahjongGame struct {
	scoringState ScoringState
	gameState    GameState
	mahjongRound MahjongRound
	ordering     Ordering

	roundWind   Wind
	roundNumber uint8
}

func NewMahjongGame() *MahjongGame {
	return &MahjongGame{
		gameState: *InitGameState(),
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

func (game *MahjongGame) SendRoundSetup() (sendInfos []MessageSendInfo) {
	for arenaIdx := uint8(0); arenaIdx < 4; arenaIdx++ {
		gameIdx := game.Ordering.GameIdx(arenaIdx)
		initialTiles := game.TileState.Hands[gameIdx].ClosedHand.GetHand()
		setup := []Setup{
			{
				Type: INITIAL_TILES,
				Data: initialTiles,
			},
		}
		sendInfos = append(sendInfos, MessageSendInfo{
			Events: []BoardEvent{
				GameSetupEvent{Setup: setup},
			},
			SendTo: arenaIdx,
		})
	}

	return sendInfos
}

func (game *MahjongGame) StartGame() (messages []MessageSendInfo, err error) {
	err = game.gameState.Transition("start-game", game, &messages)
	return messages, err
}

func (game *MahjongGame) StartRound() ([]MessageSendInfo, error) {
	if err := game.GameState.Transition(
		"start-round",
	); err != nil {
		return nil, err
	}

	game.TileState = CreateNewRound()

	setup := game.SendRoundSetup()
	newTurnInfo, err := game.handleNewTurn(0)
	if err != nil {
		return nil, err
	}

	return append(setup, newTurnInfo...), nil
}

func (game *MahjongGame) HandleEvent(action Action, arenaIdx uint8) ([]MessageSendInfo, error) {
	gameIdx := game.Ordering.ArenaToGame[arenaIdx]
	nextTurnInfo, err := game.handleNewTurn(gameIdx)
	if err != nil {
		return nil, err
	}

	return nextTurnInfo, nil
}

func (game *MahjongGame) RoundEnd() error {
	game.RoundState.IncrementRound()
	return nil
}

func (game *MahjongGame) RoundEnded() bool {
	// Do the increment post-round
	return false
}

func (game *MahjongGame) GameEnd() error {
	return nil
}

func (game *MahjongGame) GameEnded() bool {
	return false
}
