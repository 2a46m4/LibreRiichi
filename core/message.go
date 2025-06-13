package core

import (
	"encoding/json"
	"fmt"
)

type ArenaMessageType uint8

const (
	// Messages that are sent from game (server) to player (client)
	PlayerJoinedEventType ArenaMessageType = iota
	GameStartedEventType
	SetupEventType
	PlayerActionEventType
	PotentialActionEventType

	// Messages that are sent from player (client) to game (server)
	StartGameActionType
	PlayerActionType
	QuitActionType
)

// ArenaMessage are messages that are sent between clients and server
type ArenaMessage struct {
	MessageType ArenaMessageType `json:"message_type"`
	Data        json.RawMessage  `json:"data"`
}

type ServerArenaHandler interface {
	HandleStartGameAction(StartGameActionTypeData) error
	HandlePlayerAction(PlayerActionTypeData) error
	HandleQuitAction(QuitActionTypeData) error
}

type ClientArenaHandler interface {
	HandlePlayerJoinedEvent(PlayerJoinedEventTypeData) error
	HandleGameStartedEvent(GameStartedEventTypeData) error
	HandleSetupEvent(SetupEventTypeData) error
	HandlePlayerActionEvent(PlayerActionEventTypeData) error
	HandlePotentialActionEvent(PotentialActionEventTypeData) error
}

type PlayerJoinedEventTypeData struct{}

type GameStartedEventTypeData struct{}

type SetupEventTypeData struct {
	Setup Setup `json:"setup"`
}

type StartGameActionTypeData struct{}

type PlayerActionEventTypeData struct {
	Action     ActionType      `json:"action_type"`
	FromPlayer uint8           `json:"from_player"`
	Data       json.RawMessage `json:"data"`
}

type PotentialActionEventTypeData struct {
	Action ActionType      `json:"action_type"`
	Data   json.RawMessage `json:"data"`
}

type PlayerActionTypeData struct {
	Action ActionType      `json:"action_type"`
	Data   json.RawMessage `json:"data"`
}

type QuitActionTypeData struct{}

type PlayerJoinedEventData struct {
	Client Client `json:"client"`
}

func ServerArenaDecodeAndDispatch(handler ServerArenaHandler, rawData []byte) error {
	var raw struct {
		MessageType ArenaMessageType `json:"message_type"`
		Data        json.RawMessage  `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	switch raw.MessageType {
	case PlayerActionType:
		message := PlayerActionTypeData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandlePlayerAction(message)
	case QuitActionType:
		message := QuitActionTypeData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleQuitAction(message)
	case StartGameActionType:
		message := StartGameActionTypeData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleStartGameAction(message)
	default:
		return fmt.Errorf("unexpected core.ArenaMessageType: %#v", raw.MessageType)
	}
}

func ClientArenaDecodeAndDispatch(handler ClientArenaHandler, rawData []byte) error {
	var raw struct {
		MessageType ArenaMessageType `json:"message_type"`
		Data        json.RawMessage  `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	switch raw.MessageType {
	case PlayerJoinedEventType:
		message := PlayerJoinedEventTypeData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandlePlayerJoinedEvent(message)
	case GameStartedEventType:
		message := GameStartedEventTypeData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleGameStartedEvent(message)
	case SetupEventType:
		message := SetupEventTypeData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleSetupEvent(message)
	case PlayerActionEventType:
		message := PlayerActionEventTypeData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandlePlayerActionEvent(message)
	case PotentialActionEventType:
		message := PotentialActionEventTypeData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandlePotentialActionEvent(message)
	default:
		return fmt.Errorf("unexpected core.ArenaMessageType: %#v", raw.MessageType)
	}
}
