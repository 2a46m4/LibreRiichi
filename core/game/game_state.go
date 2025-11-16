package game

import (
	"context"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	"github.com/looplab/fsm"
)

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
			"start-game": gameState.HandleStartGame,
			"start-round": gameState.HandleStartRound,
			"": gameState.HandleStartRound,
		},
	)
	return &gameState
}

func (gameState *GameState) Transition(event string, arguments ...any) error {
	return gameState.FSM.Event(gameState.context, event, arguments...)
}

func (gameState *GameState) HandleStartGame(context context.Context, event *fsm.Event) {
	game := event.Args[0].(*MahjongGame)
	returnValues := event.Args[1].(*[]MessageSendInfo)

	game.ScoringState = InitScoring(25000)
	game.RoundState = *InitRoundState()
	game.WindState = 0
	game.Ordering = InitRandomOrdering()

	*returnValues = game.SendGameSetup()
}

func (gameState *GameState) HandleStartRound(context context.Context, event *fsm.Event) {
	roundState := event.Args[0].(*MahjongRoundState)
	returnValues := event.Args[2].(*[]MessageSendInfo)
	event.Err = roundState.Transition("start-round", roundState, returnValues)
}

func (gameState *GameState) HandleEvent(context context.Context, event *fsm.Event) {
	roundState := event.Args[0].(*MahjongRoundState)
	ordering := event.Args[1].(*Ordering)
	action := event.Args[2].(Action)
	arenaIdx := event.Args[3].(uint8)
	returnValues := event.Args[4].(*[]MessageSendInfo)

	nextTurnInfo, err := roundState.handleNewTurn(ordering.GameIdx(arenaIdx))
	if err != nil {
		return nil, err
	}

	returnValues = nextTurnInfo
}
