package game

import (
	"context"
	"errors"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	"github.com/looplab/fsm"
)

// Stores the state of the current game (in game, in round, etc.)
type GameState struct {
	*fsm.FSM
	context context.Context
}

func InitGameState() *GameState {
	gameState := GameState{
		FSM:     &fsm.FSM{},
		context: nil,
	}

	gameState.FSM = fsm.NewFSM(
		"out-of-game",
		fsm.Events{
			fsm.EventDesc{
				Name: "start-game",
				Src: []string{
					"out-of-game",
				},
				Dst: "in-game",
			},
			fsm.EventDesc{
				Name: "start-round",
				Src: []string{
					"in-game",
				},
				Dst: "in-round",
			},
			fsm.EventDesc{
				Name: "in-round-event",
				Src: []string{
					"in-round",
				},
				Dst: "in-round",
			},
			fsm.EventDesc{
				Name: "round-end",
				Src: []string{
					"in-round",
				},
				Dst: "in-game",
			},
			fsm.EventDesc{
				Name: "game-end",
				Src: []string{
					"in-game",
				},
				Dst: "out-of-game",
			},
		},
		fsm.Callbacks{
			"before_start-game": gameState.CheckStartGamePossible,
			"start-game":        gameState.HandleStartGame,
			"start-round":       gameState.HandleStartRound,
			"handle-event":      gameState.HandleEvent,
		},
	)
	return &gameState
}

func (gameState *GameState) Transition(event string, arguments ...any) error {
	return gameState.FSM.Event(gameState.context, event, arguments...)
}

func getGameSetup(roundData MahjongRoundData,
	ordering Ordering, scoringState ScoringState) (sendInfos []MessageSendInfo) {

	// Create setup data for each player
	for arenaIdx := uint8(0); arenaIdx < 4; arenaIdx++ {
		gameIdx := ordering.GameIdx(arenaIdx)

		setup := []Setup{
			{
				Type: DORA,
				Data: roundData.tileState.DeadWall.getLastDoraTile(),
			},
			{
				Type: PLAYER_NUMBER,
				Data: ordering.GameIdx(arenaIdx),
			},
			{
				Type: ROUND_NUMBER,
				Data: uint8(0), // First round
			},
			{
				Type: ROUND_WIND,
				Data: roundData.windState.GetPlayerWind(gameIdx), // Get player's seat wind
			},
			{
				Type: STARTING_POINTS,
				Data: [4]uint32{
					scoringState.Points[0],
					scoringState.Points[1],
					scoringState.Points[2],
					scoringState.Points[3],
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

func (gameState *GameState) CheckStartGamePossible(context context.Context, event *fsm.Event) {
	// TODO: Do some checks that starting the game is possible
	event.Cancel(errors.New("Hello"))
}

func (gameState *GameState) HandleStartGame(context context.Context, event *fsm.Event) {
	game := event.Args[0].(*MahjongGame)

	game.scoringState = InitScoring(25000)
	game.mahjongRound = InitMahjongRound(0)
	game.ordering = InitRandomOrdering()

	setup := getGameSetup(game.mahjongRound.data, game.ordering, game.scoringState)
	gameState.SetMetadata("return", setup)
}

func (gameState *GameState) CheckStartRoundPossible(context context.Context, event *fsm.Event) {
	// TODO: Do some checks

	// Actually transition here, because we can signal a failure here
	round := event.Args[0].(*MahjongRound)
	err := round.roundState.Transition("start-round", round)
	if err != nil {
		event.Cancel(err)
		return
	}
}

func (gameState *GameState) HandleStartRound(context context.Context, event *fsm.Event) {
	round := event.Args[0].(*MahjongRound).roundState
	if round.RoundFSM.Current() != "pre-draw" {
		panic("Should have already transitioned")
	} else {
		data, hasData := round.GetReturn()
		if hasData {
			gameState.setReturn(data)
		}
	}
}

func (gameState *GameState) HandleEvent(context context.Context, event *fsm.Event) {
	round := event.Args[0].(*MahjongRound)
	ordering := event.Args[1].(*Ordering)
	action := event.Args[2].(Action)
	arenaIdx := event.Args[3].(uint8)
	returnValues := event.Args[4].(*[]MessageSendInfo)
	var doContinue bool

	round.roundState.Transition("handle-event")

	nextTurnInfo, err := roundState.handleNewTurn(ordering.GameIdx(arenaIdx))
	if err != nil {
		return nil, err
	}

	returnValues = nextTurnInfo
}

func (gameState *GameState) GetReturn() (data any, ok bool) {
	data, ok = gameState.Metadata("return")
	gameState.DeleteMetadata("return")
	return data, ok
}

func (gameState *GameState) setReturn(data any) {
	gameState.SetMetadata("return", data)
}
