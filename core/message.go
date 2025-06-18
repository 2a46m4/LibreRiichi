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
}

type ClientArenaHandler interface {
	HandlePlayerJoinedEvent(PlayerJoinedEventData) error
	HandlePlayerQuitEvent(PlayerQuitEventData) error
	HandleGameStartedEvent(GameStartedEventData) error
	HandleArenaBoardEvent(ArenaBoardEventData) error
}
type PlayerJoinedEventData struct {
	Client Client `json:"client"`
}

type GameStartedEventData struct{}

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
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
		msg.Data = data
	case GameStartedEventType:
		data := GameStartedEventData{}
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
		msg.Data = data
	case PlayerActionType:
		data := PlayerActionData{}
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
		msg.Data = data
	case PlayerJoinedEventType:
		data := PlayerJoinedEventData{}
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
		msg.Data = data
	case PlayerQuitActionType:
		data := PlayerQuitActionData{}
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
		msg.Data = data
	case PlayerQuitEventType:
		data := PlayerQuitEventData{}
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
		msg.Data = data
	case StartGameActionType:
		data := StartGameActionData{}
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
		msg.Data = data
	default:
		panic(fmt.Sprintf("unexpected core.ArenaMessageType: %#v", raw.MessageType))
	}

	return nil
}

func ServerArenaDispatch(handler ServerArenaHandler, msg ArenaMessage) error {
	switch msg.MessageType {
	case PlayerActionType:
		message, ok := msg.Data.(PlayerActionData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandlePlayerAction(message)
	case PlayerQuitActionType:
		message, ok := msg.Data.(PlayerQuitActionData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandlePlayerQuitAction(message)
	case StartGameActionType:
		message, ok := msg.Data.(StartGameActionData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandleStartGameAction(message)
	default:
		return fmt.Errorf("unexpected core.ArenaMessageType: %#v", msg.MessageType)
	}
}

func ClientArenaDispatch(handler ClientArenaHandler, msg ArenaMessage) error {
	switch msg.MessageType {
	case PlayerJoinedEventType:
		message, ok := msg.Data.(PlayerJoinedEventData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandlePlayerJoinedEvent(message)
	case PlayerQuitEventType:
		message, ok := msg.Data.(PlayerQuitEventData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandlePlayerQuitEvent(message)
	case GameStartedEventType:
		message, ok := msg.Data.(GameStartedEventData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandleGameStartedEvent(message)
	case ArenaBoardEventType:
		message, ok := msg.Data.(ArenaBoardEventData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandleArenaBoardEvent(message)
	default:
		return fmt.Errorf("unexpected core.ArenaMessageType: %#v", msg.MessageType)
	}
}
