package game

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

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
	log *slog.Logger

	scoring  Scoring
	turnData TurnData
	tileData TileData

	outgoingRequests util.Set[OutgoingNakiRequest]

	// It's important to serialize accesses to these channels,
	// otherwise the ordering will be wrong
	Inbox chan any
	// The roundstate replies through this channel
	Reply chan any

	gameStarted  bool
	roundStarted bool
}

func InitRoundState() RoundState {
	return RoundState{
		log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})),
		scoring:          InitScoring(25000),
		turnData:         InitTurnData(),
		tileData:         CreateNewRound(),
		outgoingRequests: util.NewSet[OutgoingNakiRequest](),
		gameStarted:      false,
		roundStarted:     false,
		Inbox:            make(chan any),
		Reply:            make(chan any),
	}
}

// Transition to a await toss state
//
// The player either draws the tile and can discard any tile in their
// closed hand, or must discard the most recently tossed tile if they are in Riichi
func (roundState *RoundState) drawTile() []MessageSendInfo {
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

	return ret
}

func (roundState *RoundState) discardTile(action Action, playerIdx uint8) (info []MessageSendInfo, err error, shouldContinue bool) {
	switch action := action.(type) {
	case Toss:
		info, err = roundState.handleToss(action, playerIdx)
		shouldContinue = true
	case Riichi:
		info, err = roundState.handleRiichi(action, playerIdx)
		shouldContinue = true
	case Tsumo:
		info, err = roundState.handleTsumo(action, playerIdx)
		if err != nil {
			shouldContinue = true
		} else {
			shouldContinue = false
		}
		
	default:
		return nil, errors.New("Wrong action"), true
	}
	return info, err, shouldContinue
}

func (roundState *RoundState) handleTsumo(action Tsumo, playerIdx uint8) ([]MessageSendInfo, error) {
	panic("unimplemented")
}

func (roundState *RoundState) handleRiichi(action Riichi, playerIdx uint8) ([]MessageSendInfo, error) {
	panic("unimplemented")
}

func (roundState *RoundState) handleToss(action Toss, playerIdx uint8) ([]MessageSendInfo, error) {
	// Perform checks to see if its possible
	if playerIdx != roundState.turnData.CurrentPlayer {
		return nil, errors.New(fmt.Sprint("Unexpected turn number: Got ", playerIdx, " but expected ", roundState.turnData.CurrentPlayer))
	}

	hand := &roundState.tileData.Hands[playerIdx]
	lastTile, err := hand.TileJustReceived()
	if err != nil {
		return nil, err
	}

	// Riichi must toss the last tile
	if hand.InRiichi && (lastTile != action.TileToToss) {
		return nil, errors.New("If the hand is in riichi, it must toss the last tile")
	}

	err = roundState.tileData.Hands[playerIdx].TestDiscard(action.TileToToss)
	if err != nil {
		return nil, errors.New("TestDiscard failed")
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
		}
	}

	if !waitingNakiCalls {
		for i := range uint8(4) {
			res = append(res, MessageSendInfo{
				Events: []BoardEvent{NoNakiEvent{}},
				SendTo: i,
			})
		}

		res = append(res, MessageSendInfo{
			Events: []BoardEvent{tossEvent},
			SendTo: 0,
		})
	}

	// Now we need to wait for naki returns, or skip if there
	// aren't any.
	return res, nil
}

func (roundState *RoundState) handleNaki(action Action, playerIdx uint8) ([]MessageSendInfo, error, bool) {
	// TODO: Handle naki (call) logic
	// Process player making a call (chi, pon, kan)

	switch action := action.(type) {
	case Chii:
	case Pon:
	case Kan:
	case Skip:
		roundState.noNaki(action, playerIdx)
	}

	return nil, nil, true
}

func (roundState *RoundState) noNaki(skipAction Skip, playerIdx uint8) {
	if isOutgoingRequest(OutgoingNakiRequest{
		Action: skipAction,
		Player: playerIdx,
	}, roundState) {
		roundState.outgoingRequests.Remove(OutgoingNakiRequest{
			Action: skipAction.ActionToSkip,
			Player: playerIdx,
		})
	}
}

func (roundState *RoundState) roundDraw() {
	// TODO: Handle round draw logic
	// Process exhaustive draw scenario

	panic("TODO")
}

func (roundState *RoundState) roundWin() {
	// TODO: Handle round win logic
	// Process player winning the round (tsumo/ron)
	roundState.turnData.NextRound(false)
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

func (round *RoundState) gameLoop() {
	// We only loop through the round. Once the rounds are
	// finished, the game should be finished as well.
	msg := <- round.Inbox
	switch msg.(type) {
	case startGame:
		round.Reply <- nil
		round.gameStarted = true
	default:
		round.handleOtherReplies(msg)
	}

	for msg := range round.Inbox { // Round loop
		switch msg.(type) {
		case startRound:
			round.roundStarted = true
		default:
			round.handleOtherReplies(msg)
			continue
		}

		shouldContinueRound := round.roundLoop()

		// Handle any post round events here
		// Wait for round cleanup call
		round.roundStarted = false
		ROUND_CLEANUP:
			for cleanup := range round.Inbox {
				switch cleanup := cleanup.(type) {
				case roundEndCleanup:
					// TODO, send replies back about the round results
					break ROUND_CLEANUP
				default:
					round.handleOtherReplies(cleanup)
					continue
				}	
			}

		if !shouldContinueRound {
			break 
		}
	}

	// Game ended, compute ending statistics
	// TODO
	// Wait for game end check and cleanup message
	round.gameStarted = false
	GAME_CLEANUP:
		for msg = range round.Inbox {
			switch msg.(type) {
			case gameEndCleanup:
				// TODO: Game end cleanup
				break GAME_CLEANUP
			default:
				round.handleOtherReplies(msg)
			}
		}
}

// TODO: We can reply to queries that don't modify the round state in here
func (round *RoundState) handleOtherReplies(msg any) {
	switch msg.(type) {
	case hasGameStarted:
		round.Reply <- round.gameStarted
	case hasRoundStarted:
		round.Reply <- round.roundStarted
	case shouldContinueRound: // Is this right?
		round.Reply <- round.roundStarted
	case shouldGameEnd:
		round.Reply <- round.gameStarted
	default:
		round.Reply <- errors.New("Wrong state")
	}
}

// Returns whether or not we should continue to the next round (game end)
// Does not handle any round teardown
func (roundState *RoundState) roundLoop() bool {
	// The first thing we do is reply, the start round command is waiting for a response
	setup := getRoundSetup(roundState)
	ROUND_LOOP:
	for {
		// Draw tile, appending setup if needed
		if setup != nil {
			roundState.Reply <- append(setup, roundState.drawTile()...)
			setup = nil
		} else {
			roundState.Reply <- roundState.drawTile()
		}

		// Wait for riichi, toss, or tsumo calls
		for ret := range roundState.Inbox {
			switch ret := ret.(type) {
			case handleEvent:
				msg, err, shouldContinue := roundState.discardTile(ret.event, ret.from)
				if err != nil {
					roundState.Reply <- err
					break
				}
				roundState.Reply <- msg 
				if !shouldContinue {
					break ROUND_LOOP
				}
			default:
				roundState.handleOtherReplies(ret)
			}
		}

		// Check any waiting naki calls
		for !roundState.outgoingRequests.Empty() {
			msg := <-roundState.Inbox
			switch msg := msg.(type) {
			case handleEvent:
				ret, err, shouldContinue := roundState.handleNaki(msg.event, msg.from)
				if err != nil {
					roundState.Reply <- err
					break
				}
				roundState.Reply <- ret
				if !shouldContinue {
					break ROUND_LOOP
				}
			default:
				roundState.handleOtherReplies(msg)
			}
		}

		// Check if the round should still continue
		// This happens when there are no more draws left.
		// Otherwise, draw a new tile
		if roundState.tileData.LiveWall.End() {
			// Handle draws
			break ROUND_LOOP
		}
	}
	return false
}

