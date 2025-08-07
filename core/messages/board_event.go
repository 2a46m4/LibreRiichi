package core

//go:generate go run ../generate_message.go -- BoardEvent

import . "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"


type BoardEvent interface {
	BoardEventWrapper()
}

type PlayerActionEventData struct {
	Action Action     `json:"action_data"` // wrap
	FromPlayer uint8 `json:"from_player"`
}

type PotentialActionEventData struct {
	Action Action `json:"action_data"` // wrap
}

type GameSetupEventData struct {
	Setup []Setup `json:"setup"`
}

type GameEndEventData struct {
	GameResult GameResult `json:"result"`
}
