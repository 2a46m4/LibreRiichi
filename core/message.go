package core

import (
	"encoding/json"
	"fmt"
)

type ArenaMessageType uint8

const (
	// Messages that are sent from game (server) to player (client)
	PlayerJoinedEventType ArenaMessageType = iota
	PlayerQuitEventType
	GameStartedEventType
	ArenaBoardEventType

	// Messages that are sent from player (client) to game (server)
	StartGameActionType
	PlayerActionType
	PlayerJoinActionType
	PlayerQuitActionType
)

// ArenaMessage are messages that are sent between clients and server
// Should only indicate things that change the arena, not the game
type ArenaMessage struct {
	MessageType ArenaMessageType `json:"message_type"`
	Data        any              `json:"data"`
}

type ServerArenaHandler interface {
	HandleStartGameAction(StartGameActionData) error
	HandlePlayerAction(PlayerActionData) error
	HandlePlayerQuitAction(PlayerQuitActionData) error
	HandlePlayerJoinAction(PlayerJoinActionData) error
}

type ClientArenaHandler interface {
	HandlePlayerJoinedEvent(PlayerJoinedEventData) error
	HandlePlayerQuitAction(PlayerQuitEventData) error
	HandleGameStartedEvent(GameStartedEventData) error
	HandleArenaBoardEvent(ArenaBoardEventData) error
}
type PlayerJoinedEventData struct {
	Client Client `json:"client"`
}

type GameStartedEventData struct{}

type PlayerJoinActionData struct{}

type StartGameActionData struct{}

type ArenaBoardEventData struct {
	BoardEvent // For handling generic games, this should be replaced
}

type PlayerActionData struct {
	ActionData // For handling generic games, this should be replaced
}

type PlayerQuitActionData struct{}

type PlayerQuitEventData struct{}

func (msg *ArenaMessage) UnmarshalJSON(rawData []byte) error {
	var raw struct {
		MessageType ArenaMessageType `json:"message_type"`
		Data        json.RawMessage  `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	msg.MessageType = raw.MessageType

	switch raw.MessageType {
	case ArenaBoardEventType:
		data := ArenaBoardEventData{}
		err := json.Unmarshal(raw.Data, &initialMessageEvent)
		if err != nil {
			return err
		}
		msg.Data = InitialMessageEvent
	case GameStartedEventType:
	case PlayerActionType:
	case PlayerJoinActionType:
	case PlayerJoinedEventType:
	case PlayerQuitActionType:
	case PlayerQuitEventType:
	case StartGameActionType:
	default:
		panic(fmt.Sprintf("unexpected core.ArenaMessageType: %#v", raw.MessageType))
	}

	return nil
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
		message := PlayerActionData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandlePlayerAction(message)
	case PlayerJoinActionType:
		message := PlayerJoinActionData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandlePlayerJoinAction(message)
	case PlayerQuitActionType:
		message := PlayerQuitActionData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandlePlayerQuitAction(message)
	case StartGameActionType:
		message := StartGameActionData{}
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
		message := PlayerJoinedEventData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandlePlayerJoinedEvent(message)
	case PlayerQuitEventType:
		message := PlayerJoinedEventData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandlePlayerJoinedEvent(message)
	case GameStartedEventType:
		message := GameStartedEventData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleGameStartedEvent(message)
	case ArenaBoardEventType:
		message := ArenaBoardEventData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleArenaBoardEvent(message)
	default:
		return fmt.Errorf("unexpected core.ArenaMessageType: %#v", raw.MessageType)
	}
}
