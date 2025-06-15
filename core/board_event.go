package core

import (
	"encoding/json"
	"fmt"
)

type BoardEventType uint8

const (
	// An action that a player performed, affecting the board state
	PlayerActionEventType BoardEventType = iota
	// A potential action available to the player
	PotentialActionEventType
	// A setup event
	GameSetupEventType
	// A game end event
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

func BoardEventDecodeAndDispatch(handler BoardEventHandler, rawData []byte) error {
	var raw struct {
		MessageType BoardEventType  `json:"message_type"`
		Data        json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	switch raw.MessageType {
	case GameEndEventType:
		data := GameEndEventData{}
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return err
		}
		return handler.HandleGameEndEventType(data)
	case GameSetupEventType:
		data := GameSetupEventData{}
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return err
		}
		return handler.HandleGameSetupEventType(data)
	case PlayerActionEventType:
		data := PlayerActionEventData{}
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return err
		}
		return handler.HandlePlayerActionEventType(data)
	case PotentialActionEventType:
		data := PotentialActionEventData{}
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return err
		}
		return handler.HandlePotentialActionEventType(data)
	default:
		panic(fmt.Sprintf("unexpected core.BoardEventType: %#v", raw.MessageType))
	}
}
