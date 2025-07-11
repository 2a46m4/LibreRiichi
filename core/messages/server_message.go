package core

import (
	"encoding/json"
	"fmt"
	"time"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/errors"
)

type MessageType uint8

const (
	// Messages that are sent from server to client in response to an event
	ServerArenaEventType MessageType = iota

	// Messages sent in response to an action
	GenericResponseType
	ListArenasResponseType
	ArenaInfoResponseType

	// Messages that are sent from client to server
	InitialMessageActionType
	JoinArenaActionType
	ServerArenaActionType
	ListArenasActionType
	CreateArenaActionType
	ArenaInfoActionType
)

type Message struct {
	MessageType  MessageType `json:"message_type"`
	MessageIndex uint        `json:"message_index"`
	Data         any         `json:"data"`
}

type AgentInfo struct {
	Name string `json:"name"`
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
	Success   bool     `json:"success"`
	ArenaList []string `json:"arena_list"`
}

type ArenaInfoResponseData struct {
	Success     bool        `json:"success"`
	Name        string      `json:"name"`
	Agents      []AgentInfo `json:"agents"`
	GameStarted bool        `json:"game_started"`
	DateCreated time.Time   `json:"date_created"`
}

// ==================== ACTIONS ====================

type ServerActionHandler[Return any] interface {
	HandleInitialMessage(InitialMessageActionData) (Return, error)
	HandleJoinArena(JoinArenaActionData) (Return, error)
	HandleServerArena(ServerArenaActionData) (Return, error)
	HandleListArenas(ListArenasActionData) (Return, error)
	HandleCreateArena(CreateArenaActionData) (Return, error)
	HandleGetArenaInfo(ArenaInfoActionData) (Return, error)
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

type ArenaInfoActionData struct{}

// ==================== DECODE AND DISPATCH ====================

func (msg *Message) UnmarshalJSON(rawData []byte) error {
	var raw struct {
		MessageType  MessageType     `json:"message_type"`
		MessageIndex uint            `json:"message_index"`
		Data         json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	msg.MessageType = raw.MessageType
	msg.MessageIndex = raw.MessageIndex

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
	case ArenaInfoActionType:
		data := ArenaInfoActionData{}
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return err
		}
		msg.Data = data
	default:
		return fmt.Errorf("unexpected web.MessageType: %#v during unmarshalling", raw.MessageType)
	}

	return nil
}

func ServerActionDispatch[Return any](handler ServerActionHandler[Return], message Message) (ret Return, err error) {
	switch message.MessageType {
	case InitialMessageActionType:
		initialMessageReturn, ok := message.Data.(InitialMessageActionData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleInitialMessage(initialMessageReturn)
	case ServerArenaActionType:
		serverArenaAction, ok := message.Data.(ServerArenaActionData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleServerArena(serverArenaAction)
	case ListArenasActionType:
		data, ok := message.Data.(ListArenasActionData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleListArenas(data)
	case CreateArenaActionType:
		data, ok := message.Data.(CreateArenaActionData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleCreateArena(data)
	case JoinArenaActionType:
		data, ok := message.Data.(JoinArenaActionData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleJoinArena(data)
	case ArenaInfoActionType:
		data, ok := message.Data.(ArenaInfoActionData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleGetArenaInfo(data)
	default:
	}

	return ret, fmt.Errorf("unexpected web.MessageType: %#v during dispatch", message.MessageType)
}
