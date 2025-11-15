package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
)

type MahjongGame struct {
	ScoringState
	GameState
	MahjongRoundState
	Ordering
}

type MahjongRoundState struct {
	TileState
	RoundState
	WindState
	TurnState
}

func NewMahjongGame() *MahjongGame {
	return &MahjongGame{
		GameState: InitGameState(),
	}
}

func (game *MahjongGame) handleNewTurn(playerIdx uint8) ([]MessageSendInfo, error) {
	err := game.GameState.Transition("draw-tile")
	if err != nil {
		return nil, err
	}
	err = game.TurnState.Try(Draw{}, playerIdx)
	if err != nil {
		return nil, err
	}

	draw := game.TileState.Draw(playerIdx)
	_ = draw // Use the variable to avoid "declared and not used" error

	// Draw, checking that we still have moves
	// Check riichi, tsumo, kan
	return nil, nil
}

func (game *MahjongGame) SendGameSetup() (sendInfos []MessageSendInfo) {

	// Create setup data for each player
	for arenaIdx := uint8(0); arenaIdx < 4; arenaIdx++ {
		gameIdx := game.Ordering.GameIdx(arenaIdx)

		setup := []Setup{
			{
				Type: DORA,
				Data: game.DeadWall.DoraDeadWall.getLastDoraTile(),
			},
			{
				Type: PLAYER_NUMBER,
				Data: game.Ordering.GameIdx(arenaIdx),
			},
			{
				Type: ROUND_NUMBER,
				Data: uint8(0), // First round
			},
			{
				Type: ROUND_WIND,
				Data: game.WindState.GetPlayerWind(gameIdx), // Get player's seat wind
			},
			{
				Type: STARTING_POINTS,
				Data: [4]uint32{
					game.ScoringState.Points[0],
					game.ScoringState.Points[1],
					game.ScoringState.Points[2],
					game.ScoringState.Points[3],
				},
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

func (game *MahjongGame) StartGame() ([]MessageSendInfo, error) {
	if err := game.GameState.Transition("start-game"); err != nil {
		return nil, err
	}

	game.ScoringState = InitScoring(25000)
	game.RoundState = *InitRoundState()
	game.WindState = 0
	game.Ordering = InitRandomOrdering()
	setup := game.SendGameSetup()
	return setup, nil
}

func (game *MahjongGame) StartRound() ([]MessageSendInfo, error) {
	if err := game.GameState.Transition("start-round"); err != nil {
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
