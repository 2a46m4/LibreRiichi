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

func (msg *BoardEvent) UnmarshalJSON(rawData []byte) error {
	var raw struct {
		MessageType BoardEventType  `json:"event_type"`
		Data        json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	msg.EventType = raw.MessageType
	switch raw.MessageType {
	case GameEndEventType:
		data := GameEndEventData{}
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return err
		}
		msg.Data = data
	case GameSetupEventType:
		data := GameSetupEventData{}
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return err
		}
		msg.Data = data
	case PlayerActionEventType:
		data := PlayerActionEventData{}
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return err
		}
		msg.Data = data
	case PotentialActionEventType:
		data := PotentialActionEventData{}
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return err
		}
		msg.Data = data
	default:
		return fmt.Errorf("unexpected core.BoardEventType: %#v", raw.MessageType)
	}
	return nil
}

func BoardEventDispatch(handler BoardEventHandler, event BoardEvent) error {
	switch event.EventType {
	case GameEndEventType:
		message, ok := event.Data.(GameEndEventData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandleGameEndEventType(message)
	case GameSetupEventType:
		message, ok := event.Data.(GameSetupEventData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandleGameSetupEventType(message)
	case PlayerActionEventType:
		message, ok := event.Data.(PlayerActionEventData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandlePlayerActionEventType(message)
	case PotentialActionEventType:
		message, ok := event.Data.(PotentialActionEventData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandlePotentialActionEventType(message)
	default:
		return fmt.Errorf("unexpected core.BoardEventType: %#v", event.EventType)
	}
}
