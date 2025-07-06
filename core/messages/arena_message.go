package core

import (
	"encoding/json"
	"fmt"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/errors"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	"github.com/google/uuid"
)

type ArenaMessageType uint8

const (
	// Messages that are sent from game (server) to player (client)
	PlayerJoinedEventType ArenaMessageType = iota
	PlayerQuitEventType
	GameStartedEventType
	ArenaBoardEventType
	ListPlayersResponseType

	// Messages that are sent from player (client) to game (server)
	StartGameActionType
	PlayerActionType
	PlayerQuitActionType
	ListPlayersActionType
)

// ArenaMessage are messages that are sent between clients and server
// Should only indicate things that change the arena, not the game
type ArenaMessage struct {
	MessageType ArenaMessageType `json:"message_type"`
	Data        any              `json:"data"`
}

// ==================== EVENTS ====================

type PlayerJoinedEventData struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

type PlayerQuitEventData struct {
	Name string `json:"name"`
}

type GameStartedEventData struct{}

type ArenaBoardEventData struct {
	BoardEvent // For handling generic games, this should be replaced
}

type ListPlayersResponseData struct {
	Names []string `json:"names"`
}

// ==================== ACTIONS ====================

type ArenaActionHandler[Input any] interface {
	HandleStartGameAction(StartGameActionData, Input) error
	HandlePlayerAction(PlayerActionData, Input) error
	HandlePlayerQuitAction(PlayerQuitActionData, Input) error
	HandleListPlayersAction(ListPlayersActionData, Input) error
}

type StartGameActionData struct{}

type PlayerQuitActionData struct{}

type PlayerActionData struct {
	ActionData // For handling generic games, this should be replaced
}

type ListPlayersActionData struct{}

// ==================== DECODING AND DISPATCH ====================

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
	case ListPlayersResponseType:
		data := ListPlayersResponseData{}
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
	case ListPlayersActionType:
		data := ListPlayersActionData{}
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected core.ArenaMessageType: %#v", raw.MessageType)
	}

	return nil
}

func ArenaActionDispatch[E any](handler ArenaActionHandler[E], msg ArenaMessage, input E) (err error) {
	switch msg.MessageType {
	case PlayerActionType:
		message, ok := msg.Data.(PlayerActionData)
		if !ok {
			return BadMessage{}
		}
		return handler.HandlePlayerAction(message, input)
	case PlayerQuitActionType:
		message, ok := msg.Data.(PlayerQuitActionData)
		if !ok {
			return BadMessage{}
		}
		return handler.HandlePlayerQuitAction(message, input)
	case StartGameActionType:
		message, ok := msg.Data.(StartGameActionData)
		if !ok {
			return BadMessage{}
		}
		return handler.HandleStartGameAction(message, input)
	case ListPlayersActionType:
		message, ok := msg.Data.(ListPlayersActionData)
		if !ok {
			return BadMessage{}
		}
		return handler.HandleListPlayersAction(message, input)
	default:
		return fmt.Errorf("unexpected core.ArenaMessageType: %#v", msg.MessageType)
	}
}
