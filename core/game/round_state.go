package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/looplab/fsm"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	"codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	util "codeberg.org/ijnakashiar/LibreRiichi/core/util"
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
			"before_discard-tile": roundState.discardTileTest,
			"discard-tile":     roundState.discardTile,
			"call-naki":        roundState.callNaki,
			"before_no-naki":          roundState.noNakiTest,
			"no-naki":          roundState.noNaki,
			"round-draw":       roundState.roundDraw,
			"round-win":        roundState.roundWin,
			"before_event":     roundState.infoTransition,
		},
	)
	roundState.context = context.Background()
	roundState.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	return roundState
}

func (roundState *RoundState) GenerateGraphs() string {
    return fsm.Visualize(roundState.RoundFSM)
}

func (roundState *RoundState) infoTransition(context context.Context, event *fsm.Event) {
	roundState.log.Info("Transitioning:", "from", event.Src, "to", event.Dst, "event", event.Event)
}

func getRoundSetup(tileState TileData) (sendInfos []MessageSendInfo) {
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

// Handles an event by dispatching it to the right handler in roundState and returns an error if there is an invalid transition
func (roundState *RoundState) HandleEvent(action Action, gameIdx uint8, extraInfo ...any) (msg []MessageSendInfo, err error) {
    args := append([]any{action, gameIdx}, extraInfo...)

    roundState.log.Info("Handling event:", "action", fmt.Sprintf("%#v", action))
    var call string
    switch action.(type) {
    case Chii: call = "call-naki"
    case Kan: call = "call-naki"
    case Pon: call = "call-naki"
    case Ron: call = "call-naki"
    case Skip:
	call = "no-naki"
    case Riichi:
	call = "discard-tile"
    case Toss:
	call = "discard-tile"
    case Tsumo:
	call = "discard-tile"
    case Draw:
	roundState.log.Error("Wrong action: %#v", action)
    default:
	roundState.log.Error("unexpected core.Action: %#v", action)
	panic(fmt.Sprintf("unexpected core.Action: %#v", action))
    }

    err = roundState.Transition(call, args...)
    if err != nil {
	panic(err)
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

	messages := getRoundSetup(round.data.tileData)

	err := event.FSM.Event(context, "pre-draw")
	err = event.FSM.Event(context, "draw-tile",
		round.data.turnData.GetExpectedDrawPlayer(),
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

	if playerIdx != round.data.turnData.TurnNumber {
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
	action := round.data.tileData.Draw(playerIdx)

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
	playerHand := &round.data.tileData.Hands[playerIdx]

	// Check for Ankan, Riichi, Tsumo potential options
	if playerHand.InRiichi {
		tile, err := playerHand.ClosedHand.Last()
		if err != nil {
			panic("Bad state")
		}

		potentialActions.Actions = append(potentialActions.Actions, Toss{
			TileToToss: tile,
		})
	} else {
		potentialActions.Actions = append(potentialActions.Actions, Toss{
			TileToToss: tile.Invalid, // Meaning all tiles
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

	if playerIdx != round.data.turnData.TurnNumber {
		event.Cancel()
		return
	}

	hand := &round.data.tileData.Hands[playerIdx]
	lastTile, err := hand.TileJustReceived()
	if err != nil {
		event.Cancel(err)
		return
	}

	// Riichi must toss the last tile
	if hand.InRiichi && (lastTile != action.TileToToss) {
	    event.Cancel(errors.New("If the hand is in riichi, it must toss the last tile"))
		return
	}

	if !round.data.tileData.Hands[playerIdx].TestDiscard(action.TileToToss) {
		event.Cancel(errors.New("TestDiscard failed"))
		return
	}
}

func (roundState *RoundState) discardTile(context context.Context, event *fsm.Event) {
    roundState.log.Info("Discarding tiles")
    action := event.Args[0].(Toss)
    playerIdx := event.Args[1].(uint8)
    round := event.Args[2].(*MahjongRound)

    tossData := round.data.tileData.Discard(playerIdx, action.TileToToss)
    tossEvent := PlayerActionEvent{
	Action:     tossData,
	FromPlayer: playerIdx,
    }

    res := []MessageSendInfo{}
    waitingNakiCalls := false
    for i := range uint8(4) {
	// Check for any calls
	info := round.data.CheckNaki(playerIdx, i)
	for _, event := range info.Events{
	    roundState.appendAwaitingMessages(AwaitAction{
	    	PotentialActions: event.(PotentialActionEvent).Actions,
	    	SentTo:          i,
	    })
	}
	
	if len(info.Events) != 0 {
	    waitingNakiCalls = true
	    info.Events = append(info.Events, tossEvent)
	    res = append(res, info)
	} else {
	    res = append(res, MessageSendInfo{
	    	Events: []BoardEvent{tossEvent},
	    	SendTo: 0,
	    })
	}
    }

    // Transition directly to no naki state if no naki, otherwise
    // store in metadata and transition only when there is no outgoing
    // naki calls left
    if !waitingNakiCalls {
	err := roundState.RoundFSM.Event(context, "no-naki")
	if err != nil {
	    roundState.log.Info("Failed to transition to no-naki:", "err", err)
	}
    }
    

    roundState.setReturn(res)
}

func (roundState *RoundState) callNaki(context context.Context, event *fsm.Event) {
    roundState.log.Info("CallNaki called")
    panic("TODO")
    // TODO: Handle naki (call) logic
    // Process player making a call (chi, pon, kan)
}

// Two possible ways to trigger no naki:
//
// 1. Skip (Can fail the transition if there are more naki calls
// waiting for a return message)
// 
// 2. There are no naki calls possible. Then this transition accepts two calls:
//    - Skip: the action itself
//    - uint8: the player that performed the skip
// 
// This function should be called when the player calls either skip or
// when there are no naki following a discard.
func (roundState *RoundState) noNakiTest(context context.Context, event *fsm.Event) {
    msgs, ok := roundState.getAwaitingMessages()
    if !ok {
	roundState.log.Warn("awaiting messages not found")
    }
    if len(msgs) == 0 {
	return
    }
    
    skipAction := event.Args[0].(Skip)
    playerIdx := event.Args[1].(uint8)
    action := skipAction.ActionToSkip
    for msgI, msg := range msgs {
	if msg.SentTo != playerIdx {
	    continue
	}

	for i, potentialActions := range msg.PotentialActions {
	    if reflect.TypeOf(potentialActions) == reflect.TypeOf(action) {
		util.Remove(&msg.PotentialActions, i)

		if len(msg.PotentialActions) == 0 {
		    util.Remove(&msgs, msgI)
		}
	    }
	}
    }

    roundState.setAwaitingMessages(msgs)

    if len(msgs) >= 0 {
	event.Cancel()
    }
}

// Checks whether or not the game should end
func (roundState *RoundState) noNaki(context context.Context, event *fsm.Event) {
    roundState.log.Info("NoNaki called")

    roundState.RoundFSM.Event()
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

func (roundState *RoundState) appendAwaitingMessages(info ...AwaitAction) {
    msgsRaw, ok := roundState.RoundFSM.Metadata("awaiting")
    if !ok {
	roundState.RoundFSM.SetMetadata("awaiting", info)
    }
    msgs := msgsRaw.([]AwaitAction)
    roundState.RoundFSM.SetMetadata("awaiting", append(msgs, info...))
}

func (roundState *RoundState) setAwaitingMessages(info []AwaitAction) {
    roundState.RoundFSM.SetMetadata("awaiting", info)
}

func (roundState *RoundState) getAwaitingMessages() (info []AwaitAction, ok bool) {
    infoRaw, ok := roundState.RoundFSM.Metadata("awaiting")
    return infoRaw.([]AwaitAction), ok
}
