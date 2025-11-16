package game

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/looplab/fsm"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
)

type RoundState struct {
	RoundFSM *fsm.FSM
	context  context.Context

	log *slog.Logger
}

func InitRoundState() *RoundState {
	roundState := &RoundState{}
	roundState.RoundFSM = fsm.NewFSM(
		"out-of-round",
		fsm.Events{
			fsm.EventDesc{
				Name: "start-round",
				Src: []string{
					"out-of-round",
				},
				Dst: "pre-draw",
			},
			fsm.EventDesc{
				Name: "draw-tile",
				Src: []string{
					"pre-draw",
				},
				Dst: "waiting-discard",
			},
			fsm.EventDesc{
				Name: "discard-tile",
				Src: []string{
					"waiting-discard",
				},
				Dst: "waiting-naki",
			},
			fsm.EventDesc{
				Name: "call-naki",
				Src: []string{
					"waiting-naki",
				},
				Dst: "waiting-discard",
			},
			fsm.EventDesc{
				Name: "no-naki",
				Src: []string{
					"waiting-naki",
				},
				Dst: "pre-draw",
			},
			fsm.EventDesc{
				Name: "round-draw",
				Src: []string{
					"naki-finished",
				},
				Dst: "out-of-round",
			},
			fsm.EventDesc{
				Name: "round-win",
				Src: []string{
					"waiting-discard",
					"waiting-naki",
				},
				Dst: "out-of-round",
			},
		},
		fsm.Callbacks{
			"start-round":  roundState.startRound,
			"draw-tile":    roundState.drawTile,
			"discard-tile": roundState.discardTile,
			"call-naki":    roundState.callNaki,
			"no-naki":      roundState.noNaki,
			"round-draw":   roundState.roundDraw,
			"round-win":    roundState.roundWin,
		},
	)
	roundState.context = context.Background()
	roundState.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	return roundState
}

func (game *RoundState) SendRoundSetup(ordering Ordering, tileState TileState) (sendInfos []MessageSendInfo) {
	for arenaIdx := uint8(0); arenaIdx < 4; arenaIdx++ {
		gameIdx := ordering.GameIdx(arenaIdx)
		initialTiles := tileState.Hands[gameIdx].ClosedHand.GetHand()
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

func (roundState *RoundState) HandleEvent(action Action, gameIdx uint8) (msg []MessageSendInfo, err error) {

	switch action.(type) {
	case Chii:
	case Kan:
	case Pon:
	case Ron:
		err = roundState.Transition("call-naki", action, gameIdx)
	case Skip:
		err = roundState.Transition("no-naki", action, gameIdx)
	case Riichi:
	case Toss:
	case Tsumo:
		err = roundState.Transition("discard-tile", action, gameIdx)
	case Draw:
		roundState.log.Error("Wrong action: %#v", action)
	default:
		roundState.log.Error("unexpected core.Action: %#v", action)
		panic(fmt.Sprintf("unexpected core.Action: %#v", action))
	}

	ret, _ := roundState.GetReturn()
	return ret.([]MessageSendInfo), err
}


func (roundState *RoundState) Transition(event string, args ...any) error {
	return roundState.RoundFSM.Event(roundState.context, event, args...)
}

func (roundState *RoundState) startRound(context context.Context, event *fsm.Event) {
	messages := []MessageSendInfo{}
	round := event.Args[0].(*MahjongRound)
	isFirstRound := event.Args[1].(bool)

	if isFirstRound {
		round.data = InitMahjongRoundData()
	} else {
		round.data.IncrementRound()		
	}

	err := event.FSM.Event(context, "draw-tile", round)
	if err != nil {
		panic("Started round but couldn't draw tile")
	}
	tileMsgs, ok := roundState.GetReturn()
	if !ok {
		panic("Getting return failed")
	}

	roundState.setReturn(append(messages, tileMsgs.([]MessageSendInfo)...))
}

func (roundState *RoundState) drawTile(context context.Context, event *fsm.Event) {
	// TODO: Handle tile drawing logic
	// Process player drawing a tile from the wall
}

func (roundState *RoundState) discardTile(context context.Context, event *fsm.Event) {
	// TODO: Handle tile discard logic
	// Process player discarding a tile
}

func (roundState *RoundState) callNaki(context context.Context, event *fsm.Event) {


	// TODO: Handle naki (call) logic
	// Process player making a call (chi, pon, kan)
}

func (roundState *RoundState) noNaki(context context.Context, event *fsm.Event) {
	// TODO: Handle no-naki logic
	// Process when no players make a call
}

func (roundState *RoundState) roundDraw(context context.Context, event *fsm.Event) {
	// TODO: Handle round draw logic
	// Process exhaustive draw scenario
}

func (roundState *RoundState) roundWin(context context.Context, event *fsm.Event) {
	// TODO: Handle round win logic
	// Process player winning the round (tsumo/ron)
}

func (roundState *RoundState) RoundEnded() bool {
	return roundState.RoundFSM.Is("out-of-round")
}

func (roundState *RoundState) GetReturn() (ret any, ok bool) {
	ret, ok = roundState.RoundFSM.Metadata("return")
	roundState.RoundFSM.DeleteMetadata("return")
	return ret, ok
}

func (roundState *RoundState) setReturn(data any) {
	roundState.RoundFSM.SetMetadata("return", data)
}
