package core

//go:generate go run ../generate_message.go -- BoardEvent

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game/game_result"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type BoardEvent interface {
	BoardEventWrapper()
}

type PlayerActionEvent struct {
	Action     Action `json:"action_data"` // wrap
	FromPlayer uint8  `json:"from_player"`
}

type PotentialActionEvent struct {
	Actions []Action `json:"actions"` // wrap
}

type GameSetupEvent struct {
	Setup []Setup `json:"setup"`
}

type GameEndEvent struct {
	GameResult GameResult `json:"result"`
}
