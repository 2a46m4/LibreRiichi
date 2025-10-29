package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type GameActionHandler struct {
	game MahjongGame
}

// HandleChii implements core.ActionHandler.
func (g GameActionHandler) HandleChii(Chii, E) (T, error) {
	panic("unimplemented")
}

// HandleDraw implements core.ActionHandler.
func (g GameActionHandler) HandleDraw(Draw, E) (T, error) {
	panic("unimplemented")
}

// HandleKan implements core.ActionHandler.
func (g GameActionHandler) HandleKan(Kan, E) (T, error) {
	panic("unimplemented")
}

// HandlePon implements core.ActionHandler.
func (g GameActionHandler) HandlePon(Pon, E) (T, error) {
	panic("unimplemented")
}

// HandleRiichi implements core.ActionHandler.
func (g GameActionHandler) HandleRiichi(Riichi, E) (T, error) {
	panic("unimplemented")
}

// HandleRon implements core.ActionHandler.
func (g GameActionHandler) HandleRon(Ron, E) (T, error) {
	panic("unimplemented")
}

// HandleSkip implements core.ActionHandler.
func (g GameActionHandler) HandleSkip(Skip, E) (T, error) {
	panic("unimplemented")
}

// HandleToss implements core.ActionHandler.
func (g GameActionHandler) HandleToss(Toss, E) (T, error) {
	panic("unimplemented")
}

// HandleTsumo implements core.ActionHandler.
func (g GameActionHandler) HandleTsumo(Tsumo, E) (T, error) {
	panic("unimplemented")
}

func x() {
	var x ActionHandler[uint8, MessageSendInfo] = GameActionHandler{}
}
