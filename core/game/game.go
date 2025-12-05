package game

import (
	"context"
	"errors"
	"log/slog"
	"os"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	"log"
	"github.com/looplab/fsm"
)

// Essentially a thin wrapper over game state and changes the ordering
type MahjongGame struct {
	gameState    GameState
	ordering     Ordering
}

func NewMahjongGame() *MahjongGame {
	return &MahjongGame{
		gameState:  *InitGameState(),
	}
}

func (game *MahjongGame) StartGame() (messages []MessageSendInfo, err error) {
	err = game.gameState.Transition("start-game")
	sendInfo, _ := game.gameState.GetReturn()
	ChangeToArenaIdx(sendInfo, game.ordering)
	return sendInfo, err
}

func (game *MahjongGame) StartRound() (msgs []MessageSendInfo, err error) {
	err = game.gameState.Transition("start-round", true)
	sendInfo, _ := game.gameState.GetReturn()
	ChangeToArenaIdx(sendInfo, game.ordering)
	return sendInfo, err
}

func (game *MahjongGame) ContinueRound() (msgs []MessageSendInfo, err error) {
    err = game.gameState.Transition("start-round", false)
    sendInfo, _ := game.gameState.GetReturn()
    ChangeToArenaIdx(sendInfo, game.ordering)
    return nil, nil
}


func (game *MahjongGame) HandleEvent(action Action, arenaIdx uint8) (msgs []MessageSendInfo, err error) {
    gameIdx := game.ordering.GameIdx(arenaIdx)
    err = game.gameState.Transition("handle-event", action, gameIdx)
    _, isNoTransition := err.(fsm.NoTransitionError)
    if !isNoTransition {
	log.Println("Error occurred: ", err)
    } else {
	err = nil
    }
    sendInfo, _ := game.gameState.GetReturn()
    ChangeToArenaIdx(sendInfo, game.ordering)
    return sendInfo, err
}

func (game *MahjongGame) HasRoundEnded() bool {
	return game.gameState.Current() == "in-game"
}

func (game *MahjongGame) RoundEndCleanup() (msgs []MessageSendInfo, err error) {
    // Do the increment post-round
    // TODO
    // game.mahjongRound.roundState


    return nil, nil
}

func (game *MahjongGame) HasGameEnded() bool {
	return game.gameState.Current() == "finished-game"
}

// TODO: Send game results
func (game *MahjongGame) GameEndCleanup() (msgs []MessageSendInfo, err error) {
	return nil, nil
}



// Stores the state of the current game (in game, in round, etc.)
type GameState struct {
    *fsm.FSM
    // We can probably use this for timeout events when waiting
    // for the user to return some input
    context context.Context
    *slog.Logger

    roundState RoundState
}

func InitGameState() *GameState {
	gameState := GameState{
		FSM:     &fsm.FSM{},
		context: context.Background(),
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})),
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
				Name: "handle-event",
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
				Dst: "finished-game",
			},
		},
		fsm.Callbacks{
			"before_start-game":   gameState.CheckStartGamePossible,
			"start-game":          gameState.HandleStartGame,
			"before_start-round":  gameState.CheckStartRoundPossible,
			"start-round":         gameState.HandleStartRound,
			"before_handle-event": gameState.CheckHandleEventPossible,
			"handle-event":        gameState.HandleEvent,
			"before_round-end":    gameState.BeforeRoundEnd,
			"round-end":           gameState.RoundEnd,
			"before_event":        gameState.generalTransition,
		},
	)
	return &gameState
}

func (gameState *GameState) generalTransition(context context.Context, event *fsm.Event) {
	gameState.Info("Transitioning:", "from", event.Src, "to", event.Dst, "event", event.Event)
}

func (gameState *GameState) Transition(event string, arguments ...any) error {
	return gameState.FSM.Event(gameState.context, event, arguments...)
}

func (gameState *GameState) CheckStartGamePossible(context context.Context, event *fsm.Event) {
    // TODO: Do some checks that starting the game is possible
    // e.g. Check if player is admin
}

func (gameState *GameState) HandleStartGame(context context.Context, event *fsm.Event) {
    isFirstRound := event.Args[1].(bool)

    gameState.Info("Handling start game")
    if isFirstRound {
	gameState.roundState = InitRoundState()
    } 
    info := gameState.roundState.StartRound()

    gameState.SetMetadata("return", info)
}

func (gameState *GameState) CheckStartRoundPossible(context context.Context, event *fsm.Event) {
    // TODO: Do some checks
    // TODO: Potentially combine everything into one big FSM

	// Transition the round, because we can signal a failure here and cancel the transition
	round := event.Args[0].(*MahjongRound)
	isFirstRound := event.Args[1].(bool)
	err := round.roundState.Transition("start-round", &round.data, isFirstRound)
	if err != nil {
		event.Cancel(err)
		return
	}
}

func (gameState *GameState) HandleStartRound(context context.Context, event *fsm.Event) {
	round := event.Args[0].(*MahjongRound).roundState
	if round.RoundFSM.Current() != "waiting-discard" {
		panic("Should have already transitioned")
	} else {
		data, hasData := round.GetReturn()
		if hasData {
			gameState.setReturn(data)
		}
	}
}

func (gameState *GameState) CheckHandleEventPossible(context context.Context, event *fsm.Event) {
	
}

func (gameState *GameState) HandleEvent(context context.Context, event *fsm.Event) {
    action := event.Args[1].(Action)
    gameIdx := event.Args[2].(uint8)
    msgInfo, err := gameState.roundState.HandleEvent(action, gameIdx)
    if err != nil {
	panic("Unable to continue")
    }

    if gameState.roundState.RoundEnded() {
	gameState.setReturn(msgInfo)
	gameState.Transition("round-end")
    }
}

func (gameState *GameState) BeforeRoundEnd(context context.Context, event *fsm.Event) {
	round := event.Args[0].(*MahjongRound)

	if !round.roundState.RoundEnded() {
		event.Cancel(errors.New("Round is still going"))
	}
}

func (gameState *GameState) RoundEnd(context context.Context, event *fsm.Event) {
	// Compute some ending results, etc.

}

func (gameState *GameState) GetReturn() (data []MessageSendInfo, ok bool) {
    dataRaw, ok := gameState.Metadata("return")
    gameState.DeleteMetadata("return")
    return dataRaw.([]MessageSendInfo), ok
}

func (gameState *GameState) setReturn(data []MessageSendInfo) {
    gameState.SetMetadata("return", data)
}

func (gameState *GameState) GenerateGraphs() string {
    return fsm.Visualize(gameState.FSM)
}

// Modifies the original array
func convertToArenaIdx(msgs []MessageSendInfo, ordering Ordering) []MessageSendInfo {
	for i, msg := range msgs {
		msgs[i] = MessageSendInfo{
			Events: msg.Events,
			SendTo: ordering.ArenaIdx(msg.SendTo),
		}
	}
	return msgs
}
