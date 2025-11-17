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
					"no-naki-state",
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
				Dst: "no-naki-state",
			},
			fsm.EventDesc{
				Name: "round-draw",
				Src: []string{
					"no-naki-state",
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
			"start-round":      roundState.startRound,
			"before_draw-tile": roundState.drawTileTest,
			"draw-tile":        roundState.drawTile,
			"discard-tile":     roundState.discardTile,
			"call-naki":        roundState.callNaki,
			"no-naki":          roundState.noNaki,
			"round-draw":       roundState.roundDraw,
			"round-win":        roundState.roundWin,
		},
	)
	roundState.context = context.Background()
	roundState.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	return roundState
}

func getRoundSetup(tileState TileState) (sendInfos []MessageSendInfo) {
	for gameIdx := uint8(0); gameIdx < 4; gameIdx++ {
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
			SendTo: gameIdx,
		})
	}

	return sendInfos
}

func (roundState *RoundState) HandleEvent(action Action, gameIdx uint8, extraInfo ...any) (msg []MessageSendInfo, err error) {
	args := append([]any{action, gameIdx}, extraInfo...)

	switch action.(type) {
	case Chii:
	case Kan:
	case Pon:
	case Ron:
		err = roundState.Transition("call-naki", args)
	case Skip:
		err = roundState.Transition("no-naki", args)
	case Riichi:
	case Toss:
	case Tsumo:
		err = roundState.Transition("discard-tile", args)
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
	round := event.Args[0].(*MahjongRound)
	isFirstRound := event.Args[1].(bool)

	if isFirstRound {
		round.data = InitMahjongRoundData()
	} else {
		round.data.IncrementRound()
	}

	messages := getRoundSetup(round.data.tileState)

	err := event.FSM.Event(context, "pre-draw")
	err = event.FSM.Event(context, "draw-tile", round)
	if err != nil {
		panic("Started round but couldn't draw tile")
	}
	tileMsgs, ok := roundState.GetReturn()
	if !ok {
		panic("Getting return failed")
	}

	roundState.setReturn(append(messages, tileMsgs.([]MessageSendInfo)...))
}

// The FSM should guarantee that we are in the correct state so we
// only need to check that the person drawing the tile is correct
func (roundState *RoundState) drawTileTest(context context.Context, event *fsm.Event) {
	roundState.log.Info("Checking if draw tile can succeed")
	playerIdx := event.Args[1].(uint8)
	round := event.Args[2].(*MahjongRound)

	if playerIdx != round.data.turnState.TurnNumber {
		event.Cancel()
		return
	}
}

func (roundState *RoundState) drawTile(context context.Context, event *fsm.Event) {
	playerIdx := event.Args[1].(uint8)
	round := event.Args[2].(*MahjongRound)
	action := round.data.tileState.Draw(playerIdx)
	roundState.setReturn(action)
}

func (roundState *RoundState) discardTileTest(context context.Context, event *fsm.Event) {
	roundState.log.Info("Checking if discard tile can succeed")
	action := event.Args[0].(Toss)
	playerIdx := event.Args[1].(uint8)
	round := event.Args[2].(*MahjongRound)

	if playerIdx != round.data.turnState.TurnNumber {
		event.Cancel()
		return
	}

	if round.data.tileState.Hands[playerIdx]

	round.data.tileState.Discard()

	// if action.TileToToss
}

func (roundState *RoundState) discardTile(context context.Context, event *fsm.Event) {
	
}

func (roundState *RoundState) callNaki(context context.Context, event *fsm.Event) {

	// TODO: Handle naki (call) logic
	// Process player making a call (chi, pon, kan)
}

func (roundState *RoundState) noNaki(context context.Context, event *fsm.Event) {
	// TODO: After players make a discard and there are no naki calls
	// left, transition to here which should transition directly to another draw
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
