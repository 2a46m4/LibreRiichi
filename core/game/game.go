package game

import (
	"errors"
	"fmt"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
)

type MahjongGame struct {
	ScoringState
	TileState
	RoundState
	WindState
	TurnState
	Ordering
	GameState
}

func NewMahjongGame() *MahjongGame {
	return &MahjongGame{
		GameState: InitGameState(),
	}
}

func (game *MahjongGame) handleNewTurn() error {
	err := game.GameState.Transition(DRAW_TRANSITION)
	if err != nil {
		return err
	}
	index := game.TurnState.PlayerDraw(0)

	// Draw, checking that we still have moves
	// Check riichi, tsumo, kan

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
	if game.GameState.TurnType != OUT_OF_GAME {
		return nil, errors.New("Game already started")
	}

	game.ScoringState = InitScoring(25000)
	game.RoundState = RoundState{}
	game.WindState = 0
	game.Ordering = InitRandomOrdering()
	setup := game.SendGameSetup()
	game.GameState.Transition(GAME_START_TRANSITION)
	return setup, nil
}

func (game *MahjongGame) StartRound() ([]MessageSendInfo, error) {
	if game.GameState.TurnType != OUT_OF_GAME {
		return nil, errors.New("Round already started")
	}

	game.TileState = CreateNewRound()

	setup := game.SendRoundSetup()
	firstArenaIdx := game.Ordering.GameToArena[0]

	game.handleNewTurn()
	game.TileState.Draw(0)

	return append(setup, MessageSendInfo{
		Events: []BoardEvent{},
		SendTo: 0,
	}), nil
}

func (game *MahjongGame) HandleEvent(action Action, arenaIdx uint8) ([]MessageSendInfo, error) {
	gameIdx := game.Ordering.ArenaToGame[arenaIdx]
	fmt.Println(gameIdx)

	return nil, nil
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
