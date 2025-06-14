package core

import (
	"encoding/json"
	"fmt"
)

type BoardEventType uint8

const (
	PlayerActionEventType BoardEventType = iota
	PotentialActionEventType
	GameSetupEventType
	GameEndEventType
)

type BoardEvent struct {
	EventType BoardEventType `json:"event_type"`
	Data      any            `json:"data"`
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

type BoardEventHandler interface {
	HandlePlayerActionEventType(PlayerActionEventData) error
	HandlePotentialActionEventType(PotentialActionEventData) error
	HandleGameSetupEventType(GameSetupEventData) error
	HandleGameEndEventType(GameEndEventData) error
}

// TODO
func BoardEventDecodeAndDispatch(handler ClientArenaHandler, rawData []byte) error {
	var raw struct {
		MessageType ArenaMessageType `json:"message_type"`
		Data        json.RawMessage  `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	return nil
}
