package game

import (
	"context"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	"github.com/looplab/fsm"
)

type RoundState struct {
	RoundWind   Wind
	RoundNumber uint8

	RoundFSM *fsm.FSM
	context  context.Context
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
			"start-round":  roundState.StartRound,
			"draw-tile":    roundState.DrawTile,
			"discard-tile": roundState.DiscardTile,
			"call-naki":    roundState.CallNaki,
			"no-naki":      roundState.NoNaki,
			"round-draw":   roundState.RoundDraw,
			"round-win":    roundState.RoundWin,
		},
	)
	roundState.context = context.Background()

	return roundState
}

func (roundState *RoundState) IncrementRound() {
	roundState.RoundNumber += 1
	if roundState.RoundNumber == 4 {
		roundState.RoundWind += 1
	}
}

func (roundState *RoundState) Transition(event string, args ...any) error {
	err := roundState.RoundFSM.Event(roundState.context, event, args...)
	return err
}

func (roundState *RoundState) StartRound(context context.Context, event *fsm.Event) {
	// TODO: Initialize round start logic
	// Set up initial round state, determine dealer, etc.
}

func (roundState *RoundState) DrawTile(context context.Context, event *fsm.Event) {
	// TODO: Handle tile drawing logic
	// Process player drawing a tile from the wall
}

func (roundState *RoundState) DiscardTile(context context.Context, event *fsm.Event) {
	// TODO: Handle tile discard logic
	// Process player discarding a tile
}

func (roundState *RoundState) CallNaki(context context.Context, event *fsm.Event) {
	// TODO: Handle naki (call) logic
	// Process player making a call (chi, pon, kan)
}

func (roundState *RoundState) NoNaki(context context.Context, event *fsm.Event) {
	// TODO: Handle no-naki logic
	// Process when no players make a call
}

func (roundState *RoundState) RoundDraw(context context.Context, event *fsm.Event) {
	// TODO: Handle round draw logic
	// Process exhaustive draw scenario
}

func (roundState *RoundState) RoundWin(context context.Context, event *fsm.Event) {
	// TODO: Handle round win logic
	// Process player winning the round (tsumo/ron)
}
