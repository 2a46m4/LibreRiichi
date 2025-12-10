package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/looplab/fsm"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	"codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
	yaku "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/yaku"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	util "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type OutgoingNakiRequest struct {
	Action Action
	Player uint8
}

// Manages the state of the round and validates that turn transitions are correct
type RoundState struct {
	RoundFSM *fsm.FSM
	context  context.Context
	log      *slog.Logger

	scoring  Scoring
	turnData TurnData
	tileData TileData

	outgoingRequests util.Set[OutgoingNakiRequest]
}

func InitRoundState() RoundState {
	roundState := RoundState{
		RoundFSM: &fsm.FSM{},
		context:  context.Background(),
		log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})),
		scoring:          InitScoring(25000),
		turnData:         InitTurnData(),
		tileData:         CreateNewRound(),
		outgoingRequests: util.NewSet[OutgoingNakiRequest](),
	}

	// TODO: Decouple the receiver functions and move this in the constructor
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
			"start-round":         roundState.startRound,
			"draw-tile":           roundState.drawTile,
			"before_discard-tile": roundState.discardTile,
			"call-naki":           roundState.callNaki,
			"before_no-naki":      roundState.noNaki,
			"round-draw":          roundState.roundDraw,
			"round-win":           roundState.roundWin,
			"before_event": func(ctx context.Context, event *fsm.Event) {
				roundState.log.Info("Transitioning:",
					"from", event.Src,
					"to", event.Dst,
					"event", event.Event)
			},
		},
	)

	return roundState
}

func (roundState *RoundState) StartRound() ([]MessageSendInfo, error) {
	err := roundState.Transition("start-round")
	ret, _ := roundState.GetReturn()
	return ret, err
}

// Handles an event by dispatching it to the right handler in roundState and returns an error if there is an invalid transition
func (roundState *RoundState) HandleEvent(action Action, gameIdx uint8) (msg []MessageSendInfo, err error) {
	roundState.log.Info("Handling event:", "action", fmt.Sprintf("%#v", action), "playeridx", gameIdx)
	var call string
	switch action.(type) {
	case Chii:
		call = "call-naki"
	case Kan:
		call = "call-naki"
	case Pon:
		call = "call-naki"
	case Ron:
		call = "call-naki"
	case Skip:
		call = "no-naki"
	case Riichi:
		call = "discard-tile"
	case Toss:
		call = "discard-tile"
	case Tsumo:
		call = "discard-tile"
	case Draw:
		roundState.log.Error("Wrong action: %#v", "action", action)
	default:
		roundState.log.Error("unexpected core.Action: %#v", "action", action)
		panic(fmt.Sprintf("unexpected core.Action: %#v", action))
	}

	err = roundState.Transition(call, action, gameIdx)
	if err != nil {
		panic(fmt.Sprintf("Error: %s\nAction:%#v\nFromPlayerGameIdx:%#v", err.Error(), action, gameIdx))
	}

	ret, _ := roundState.GetReturn()
	return ret, err
}

func (roundState *RoundState) Transition(event string, args ...any) error {
	return roundState.RoundFSM.Event(roundState.context, event, args...)
}

func (roundState *RoundState) startRound(context context.Context, event *fsm.Event) {
	messages := getRoundSetup(roundState)

	err := event.FSM.Event(context, "draw-tile")
	if err != nil {
		panic("Started round but couldn't draw tile")
	}
	tileMsgs, ok := roundState.GetReturn()
	if !ok {
		panic("Getting return failed")
	}

	roundState.setReturn(append(messages, tileMsgs...))
}

// Transition to a await toss state
//
// The player either draws the tile and can discard any tile in their
// closed hand, or must discard the most recently tossed tile if they are in Riichi
func (roundState *RoundState) drawTile(context context.Context, event *fsm.Event) {
	playerIdx := roundState.turnData.CurrentPlayer
	action := roundState.tileData.Draw(playerIdx)

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
	playerHand := &roundState.tileData.Hands[playerIdx]

	// If the player is in Riichi, we can only discard the most recently obtained tile
	if playerHand.InRiichi {
		tile := playerHand.ClosedHand.Last()

		potentialActions.Actions = append(potentialActions.Actions, Toss{
			TileToToss: tile,
		})
	} else {
		potentialActions.Actions = append(potentialActions.Actions, Toss{
			TileToToss: tile.Invalid, // Meaning all tiles
		})
	}

	// Check for Ankan, Riichi, Tsumo potential options
	if CanAnkan(Kan{TileToKan: action.DrawnTile}, playerHand) {
		potentialActions.Actions = append(potentialActions.Actions, Kan{
			TileToKan: action.DrawnTile,
		})
	}

	riichiTargets := GetRiichiTargets(playerHand)
	for _, target := range riichiTargets {
		potentialActions.Actions = append(potentialActions.Actions, target)
	}

	yaku.CheckYakuAndScore(playerHand, yaku.YakuContext{
		IsSelfDrawn:           true,
		IsIppatsu:             roundState.turnData.IppatsuPossible(playerIdx),
		IsLastLiveTile:        roundState.tileData.End(),
		IsDeadWallCall:        false, // TODO
		IsFromOpponentKanCall: false,
		IsDoubleRiichi:        roundState.turnData.IsDoubleRiichi(playerIdx),
		IsTenhou:              roundState.turnData.TotalTurns == 0,
		IsChiihou:             roundState.turnData.TotalTurns < 4,
		HandInRiichi:          playerHand.InRiichi,
		RoundWind:             roundState.turnData.RoundWind,
		PlayerWind:            roundState.turnData.GetPlayerWind(playerIdx),
		IsDealer:              roundState.turnData.CurrentDealer == playerIdx,
		WinningTileIdx:        playerHand.ClosedHand.Length() - 1,
	})

	if len(potentialActions.Actions) > 0 {
		ret[playerIdx].Events = append(ret[playerIdx].Events, potentialActions)
	}

	roundState.setReturn(ret)
}

// Arguments:
//   - Toss: Action performed
//   - uint8: The player index
func (roundState *RoundState) discardTile(context context.Context, event *fsm.Event) {
	roundState.log.Info("Checking if discard tile can succeed")
	action := event.Args[0].(Toss)
	playerIdx := event.Args[1].(uint8)

	{ // Perform checks to see if its possible
		if playerIdx != roundState.turnData.CurrentPlayer {
			event.Cancel(errors.New(fmt.Sprint("Unexpected turn number: Got ", playerIdx, " but expected ", roundState.turnData.CurrentPlayer)))
			return
		}

		hand := &roundState.tileData.Hands[playerIdx]
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

		err = roundState.tileData.Hands[playerIdx].TestDiscard(action.TileToToss) 
		if err != nil {
			event.Cancel(errors.New("TestDiscard failed"))
			panic(action)
			return
		}

	}

	tossData := roundState.tileData.Discard(playerIdx, action.TileToToss)
	tossEvent := PlayerActionEvent{
		Action:     tossData,
		FromPlayer: playerIdx,
	}

	res := []MessageSendInfo{}
	waitingNakiCalls := false
	for i := range uint8(4) {
		// Check for kan, pon, chii, or ron

		if i == playerIdx {
			continue
		}

		actions := PotentialActionEvent{}
		hand := &roundState.tileData.Hands[i]

		ponResult := CheckPon(hand, action.TileToToss)
		if ponResult {
			actions.Actions = append(actions.Actions, Pon{
				TileToPon: action.TileToToss,
			})
		}

		kanResult := CheckKan(hand, action.TileToToss)
		if kanResult {
			actions.Actions = append(actions.Actions, Kan{
				TileToKan: action.TileToToss,
			})
		}

		chiiResult := GetChiiTargets(playerIdx, action.TileToToss, i, hand)
		for _, res := range chiiResult {
			actions.Actions = append(actions.Actions, res)
		}

		ronResult := CanRon(playerIdx, action.TileToToss, i, hand, yaku.YakuContext{
			IsSelfDrawn:           false,
			IsIppatsu:             false,
			IsLastLiveTile:        false,
			IsDeadWallCall:        false,
			IsFromOpponentKanCall: false,
			IsDoubleRiichi:        false,
			IsTenhou:              false,
			IsChiihou:             false,
			HandInRiichi:          false,
			RoundWind:             0,
			PlayerWind:            0,
			IsDealer:              false,
			WinningTileIdx:        i,
		})
		if ronResult != nil {
			actions.Actions = append(actions.Actions, ronResult)
		}

		info := MessageSendInfo{
			Events: []BoardEvent{actions},
			SendTo: i,
		}

		for _, event := range info.Events {
			for _, action := range event.(PotentialActionEvent).Actions {
				roundState.outgoingRequests.Add(OutgoingNakiRequest{
					Action: action,
					Player: i,
				})
			}
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

// Arguments
func (roundState *RoundState) callNaki(context context.Context, event *fsm.Event) {
	roundState.log.Info("CallNaki called")

	// TODO: Handle naki (call) logic
	// Process player making a call (chi, pon, kan)
}

// Two possible ways to trigger no naki event:
//
//  1. Skip (Can fail the transition if there are more naki calls
//     waiting for a return message)
//  2. There are no naki calls possible.
//
// Arguments:
//   - Skip: the action itself (Can be defaulted in case 2)
//   - uint8: the player that performed the skip
func (roundState *RoundState) noNaki(context context.Context, event *fsm.Event) {
	skipAction := event.Args[0].(Skip)
	playerIdx := event.Args[1].(uint8)
	if isOutgoingRequest(OutgoingNakiRequest{
		Action: skipAction,
		Player: playerIdx,
	}, roundState) {
		roundState.outgoingRequests.Remove(OutgoingNakiRequest{
			Action: skipAction.ActionToSkip,
			Player: playerIdx,
		})
	}

	if len(roundState.outgoingRequests) >= 0 {
		event.Cancel()
	}

	// Checks whether or not the game should end.
	// This happens when there are no more draws left.
	// Otherwise, draw a new tile
	if roundState.tileData.LiveWall.End() {
		err := roundState.RoundFSM.Event(context, "round-draw")
		util.PanicIf(err)
		return
	}

	err := roundState.RoundFSM.Event(context, "draw-tile")
	util.PanicIf(err)
}

func (roundState *RoundState) roundDraw(context context.Context, event *fsm.Event) {
	// TODO: Handle round draw logic
	// Process exhaustive draw scenario

	panic("TODO")
}

func (roundState *RoundState) roundWin(context context.Context, event *fsm.Event) {
	// TODO: Handle round win logic
	// Process player winning the round (tsumo/ron)
	roundState.turnData.NextRound(false)
}

func (roundState *RoundState) RoundEnded() bool {
	return roundState.RoundFSM.Is("out-of-round")
}

func (roundState *RoundState) GetReturn() (ret []MessageSendInfo, ok bool) {
	data, ok := roundState.RoundFSM.Metadata("return")
	roundState.RoundFSM.DeleteMetadata("return")
	return data.([]MessageSendInfo), ok
}

func (roundState *RoundState) setReturn(data []MessageSendInfo) {
	roundState.RoundFSM.SetMetadata("return", data)
}

// Returns true if the request was in the outgoing set
func isOutgoingRequest(request OutgoingNakiRequest, roundState *RoundState) bool {
	actionSkip, isSkip := request.Action.(Skip)
	if isSkip {
		request = OutgoingNakiRequest{
			Action: actionSkip,
			Player: request.Player,
		}
	}

	return roundState.outgoingRequests.In(request)
}

func getRoundSetup(roundData *RoundState) (sendInfos []MessageSendInfo) {

	// Create setup data for each player
	for gameIdx := range uint8(4) {
		setup := []Setup{
			{
				Type: DORA,
				Data: roundData.tileData.DeadWall.dora.getLastDoraTile(),
			},
			{
				Type: PLAYER_NUMBER, // TODO: Don't need to send twice
				Data: gameIdx,
			},
			{
				Type: ROUND_NUMBER,
				Data: roundData.turnData.GetRoundNumber(),
			},
			{
				Type: ROUND_WIND,
				Data: roundData.turnData.RoundWind,
			},
			{
				Type: STARTING_POINTS,
				Data: [4]uint32{
					roundData.scoring.Points[0],
					roundData.scoring.Points[1],
					roundData.scoring.Points[2],
					roundData.scoring.Points[3],
				},
			},
			{
				Type: INITIAL_TILES,
				Data: roundData.tileData.Hands[gameIdx].ClosedHand.GetHand(),
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

func (roundState *RoundState) GenerateGraphs() string {
	return fsm.Visualize(roundState.RoundFSM)
}
