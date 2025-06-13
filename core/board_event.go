package core

import "encoding/json"

type BoardEventType uint8

const (
	PlayerActionEventType BoardEventType = iota
	PotentialActionEventType
	GameSetupEventType
	GameEndEventType
)

type BoardEvent struct {
	EventType BoardEventType  `json:"event_type"`
	Data      json.RawMessage `json:"data"`
}

type PlayerActionEventData struct {
	ActionData
	FromPlayer uint8 `json:"from_player"`
}

type PotentialActionEventData struct {
	ActionData
}

type GameSetupEventData struct {
	Setup []Setup `json:"setup"`
}

type GameEndEventData struct {
	GameResult GameResult `json:"result"`
}
