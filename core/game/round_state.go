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

// Manages the state of the round and validates that turn transitions are correct
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
			"before_event": roundState.infoTransition,
		},
	)
	roundState.context = context.Background()
	roundState.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	return roundState
}

func (roundState *RoundState) infoTransition(context context.Context, event *fsm.Event) {
	roundState.log.Info("Transitioning: ", "from", event.Src, "to", event.Dst, "event", event.Event)
}


func getRoundSetup(tileState TileState) (sendInfos []MessageSendInfo) {
	for gameIdx := range uint8(4) {
		initialTiles := tileState.Hands[gameIdx].ClosedHand.GetHand()
		sendInfos = append(sendInfos, MessageSendInfo{
			Events: []BoardEvent{
				GameSetupEvent{Setup: []Setup{
					{
						Type: INITIAL_TILES,
						Data: initialTiles,
					},
				}},
			},
			SendTo: gameIdx,
		})
	}

	return sendInfos
}

// Handles an event and returns an error if there is an invalid transition
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
	err = event.FSM.Event(context, "draw-tile",
		round.data.turnState.GetExpectedDrawPlayer(),
		round,
	)
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
	playerIdx := event.Args[0].(uint8)
	round := event.Args[1].(*MahjongRound)

	if playerIdx != round.data.turnState.TurnNumber {
		event.Cancel()
		return
	}
}

// Transition to a await toss state
//
// The player either draws the tile and can discard any tile in their
// closed hand, or must discard the most recently tossed tile if they are in Riichi
func (roundState *RoundState) drawTile(context context.Context, event *fsm.Event) {
	playerIdx := event.Args[0].(uint8)
	round := event.Args[1].(*MahjongRound)
	action := round.data.tileState.Draw(playerIdx)

	ret := make([]MessageSendInfo, 0, 4)
	for i := range 4 {
		ret = append(ret, MessageSendInfo{
			Events: []BoardEvent{
				PlayerActionEvent{
					Action:     action,
					FromPlayer: playerIdx,
				},
			},
			SendTo: uint8(i),
		})
	}

	potentialActions := PotentialActionEvent{}
	playerHand := &round.data.tileState.Hands[playerIdx]

	// Check for Ankan, Riichi, Tsumo potential options
	if playerHand.InRiichi {
		tile, err := playerHand.ClosedHand.Last()
		if err != nil {
			panic("Bad state")
		}

		potentialActions.Actions = append(potentialActions.Actions, Toss{
			TileToToss: tile,
		})
	}

	if playerHand.TestAnKan(action.DrawnTile) {
		potentialActions.Actions = append(potentialActions.Actions, Kan{
			TileToKan: action.DrawnTile,
		})
	}

	if playerHand.TestRiichi(action.DrawnTile) {
		potentialActions.Actions = append(potentialActions.Actions, Riichi{
			TileToRiichi: action.DrawnTile,
		})
	}

	if len(potentialActions.Actions) > 0 {
		ret[playerIdx].Events = append(ret[playerIdx].Events, potentialActions)
	}

	roundState.setReturn(ret)
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

	hand := &round.data.tileState.Hands[playerIdx]
	lastTile, err := hand.TileJustReceived()
	if err != nil {
		event.Cancel()
		return
	}

	// Riichi must toss the last tile
	if hand.InRiichi && (lastTile != action.TileToToss) {
		event.Cancel()
		return
	}

	if !round.data.tileState.Hands[playerIdx].TestDiscard(action.TileToToss) {
		event.Cancel()
		return
	}
}

func (roundState *RoundState) discardTile(context context.Context, event *fsm.Event) {
	action := event.Args[0].(Toss)
	playerIdx := event.Args[1].(uint8)
	round := event.Args[2].(*MahjongRound)
	round.data.tileState.Discard(playerIdx, action.TileToToss)

	// Check for any calls
	// round.data
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
