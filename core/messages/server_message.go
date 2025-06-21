package core

import (
	"encoding/json"
	"fmt"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/errors"

	"github.com/google/uuid"
)

type MessageType uint8

const (
	// Messages that are sent from server to client in response to an event
	ServerArenaEventType MessageType = iota

	// Messages sent in response to an action
	GenericRepsonseType
	ListArenasResponseType

	// Messages that are sent from client to server
	InitialMessageActionType
	JoinArenaActionType
	ServerArenaActionType
	ListArenasActionType
	CreateArenaActionType
)

type Message struct {
	MessageType  MessageType `json:"message_type"`
	MessageIndex uint        `json:"message_index"`
	Data         any         `json:"data"`
}

// ==================== EVENTS ====================

type ServerArenaMessageEventData struct {
	ArenaMessage ArenaMessage `json:"arena_message"`
}

// ==================== RESPONSES ====================

type GenericResponseData struct {
	Success    bool   `json:"success"`
	FailReason string `json:"fail_reason"`
}

type ListArenasResponseData struct {
	ArenaList []uuid.UUID `json:"arena_list"`
}

// ==================== ACTIONS ====================

type ServerActionHandler[Input, Return any] interface {
	HandleInitialMessage(InitialMessageActionData, ...Input) (Return, error)
	HandleJoinArena(JoinArenaActionData, ...Input) (Return, error)
	HandleServerArena(ServerArenaActionData, ...Input) (Return, error)
	HandleListArenas(ListArenasActionData, ...Input) (Return, error)
	HandleCreateArena(CreateArenaActionData, ...Input) (Return, error)
}

type InitialMessageActionData struct {
	Name string `json:"name"`
}

type JoinArenaActionData struct {
	ArenaName string `json:"arena_name"`
}

type ServerArenaActionData struct {
	ArenaMessage ArenaMessage `json:"arena_message"`
}

type ListArenasActionData struct{}

type CreateArenaActionData struct {
	ArenaName string `json:"arena_name"`
}

// ==================== DECODE AND DISPATCH ====================

func (msg *Message) UnmarshalJSON(rawData []byte) error {
	var raw struct {
		MessageType MessageType     `json:"message_type"`
		Data        json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	msg.MessageType = raw.MessageType

	switch raw.MessageType {
	case InitialMessageActionType:
		initialMessageReturn := InitialMessageActionData{}
		err := json.Unmarshal(raw.Data, &initialMessageReturn)
		if err != nil {
			return err
		}
		msg.Data = initialMessageReturn
	case ServerArenaActionType:
		serverArenaAction := ServerArenaActionData{}
		err := json.Unmarshal(raw.Data, &serverArenaAction)
		if err != nil {
			return err
		}
		msg.Data = serverArenaAction
	case JoinArenaActionType:
		data := JoinArenaActionData{}
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
		msg.Data = data
	case CreateArenaActionType:
		data := CreateArenaActionData{}
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
		msg.Data = data
	case ListArenasActionType:
		data := ListArenasActionData{}
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
		msg.Data = data
	default:
		return fmt.Errorf("unexpected web.MessageType: %#v", raw.MessageType)
	}

	return nil
}

func ServerActionDispatch[Input, Return any](handler ServerActionHandler[Input, Return], message Message, input ...Input) (ret Return, err error) {
	switch message.MessageType {
	case InitialMessageActionType:
		initialMessageReturn, ok := message.Data.(InitialMessageActionData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleInitialMessage(initialMessageReturn, input...)
	case ServerArenaActionType:
		serverArenaAction, ok := message.Data.(ServerArenaActionData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleServerArena(serverArenaAction, input...)
	case ListArenasActionType:
		data, ok := message.Data.(ListArenasActionData)
		if !ok {
			return ret, BadMessage{}
		}
		handler.HandleListArenas(data, input...)
	case CreateArenaActionType:
		data, ok := message.Data.(CreateArenaActionData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleCreateArena(data, input...)
	case JoinArenaActionType:
		data, ok := message.Data.(JoinArenaActionData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleJoinArena(data, input...)
	default:
	}
	return ret, fmt.Errorf("unexpected web.MessageType: %#v", message.MessageType)
}
