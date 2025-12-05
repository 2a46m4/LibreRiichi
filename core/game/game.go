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
    ordering     Ordering
    // We can probably use this for timeout events when waiting
    // for the user to return some input
    context context.Context
    *slog.Logger
    state *fsm.FSM
    roundState RoundState
}

func NewMahjongGame() *MahjongGame {
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,
	Level:     slog.LevelDebug,
    }))
    return &MahjongGame{
	state: fsm.NewFSM(
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
		"before_event": func(ctx context.Context, e *fsm.Event) {
		    logger.Info("Transitioning:", "from", e.Src, "to", e.Dst, "event", e.Event)
		},
	    },
	),
	context: context.Background(),
	Logger: logger,
	roundState: InitRoundState(),
    }
}

func (game *MahjongGame) StartGameAndRound() (messages []MessageSendInfo, err error) {
    if game.state.Cannot("start-game") {
	return nil, errors.New("Can't start game")
    }

    err = game.state.Event(game.context, "start-game")
    if err != nil {
	return nil, err
    }

    err = game.state.Event(game.context, "start-round")
    if err != nil {
	return nil, err
    }
    sendInfo, err := game.roundState.StartRound()
    ChangeToArenaIdx(sendInfo, game.ordering)
    return sendInfo, err
}

func (game *MahjongGame) ContinueRound() (msgs []MessageSendInfo, err error) {
    if game.state.Cannot("start-round") {
	return msgs, errors.New("Can't continue round")
    }

    err = game.state.Event(game.context, "start-round")
    if err != nil {
	return msgs, err
    }

    msgs, _ = game.roundState.StartRound()
    ChangeToArenaIdx(msgs, game.ordering)
    return msgs, nil
}


func (game *MahjongGame) HandleEvent(action Action, arenaIdx uint8) (msgs []MessageSendInfo, err error) {
    if game.state.Cannot("handle-event") {
	return msgs, errors.New("Cannot handle event")
    }

    gameIdx := game.ordering.GameIdx(arenaIdx)
    msgs, err = game.roundState.HandleEvent(action, gameIdx)
    _, isNoTransition := err.(fsm.NoTransitionError)
    if !isNoTransition {
	log.Println("Error occurred: ", err)
    } else {
	err = nil
    }
    ChangeToArenaIdx(msgs, game.ordering)

    if game.roundState.RoundEnded() {
	game.state.Event(game.context, "round-end")
    }

    return msgs, err
}

func (game *MahjongGame) IsInGame() bool {
    return game.state.Current() == "in-game"
}

func (game *MahjongGame) HasRoundEnded() bool {
	return game.state.Current() != "in-round"
}

func (game *MahjongGame) RoundEndCleanup() (msgs []MessageSendInfo, err error) {
    // Do the increment post-round
    // TODO
    // game.mahjongRound.roundState


    return nil, nil
}

func (game *MahjongGame) HasGameEnded() bool {
	return game.state.Current() == "finished-game"
}

// TODO: Send game results
func (game *MahjongGame) GameEndCleanup() (msgs []MessageSendInfo, err error) {
	return nil, nil
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
