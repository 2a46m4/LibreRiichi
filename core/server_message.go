package core

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type MessageType uint8

const (
	// Messages that are sent from server to client
	InitialMessageEventType MessageType = iota
	JoinArenaEventType
	ServerArenaEventType

	// Messages that are sent from client to server
	InitialMessageActionType
	JoinArenaActionType
	ServerArenaActionType
)

type Message struct {
	MessageType MessageType     `json:"message_type"`
	Data        json.RawMessage `json:"data"`
}

// The server will accept these message types
type ServerHandler interface {
	HandleInitialMessageAction(InitialMessageActionData) error
	HandleJoinArenaAction(JoinArenaActionData) error
	HandleServerArenaAction(ServerArenaActionData) error
}

type ClientHandler interface {
	HandleInitialMessageEvent(InitialMessageEventData) error
	HandleJoinArenaEvent() error
	HandleServerArenaMessageEvent(ServerArenaMessageEventData) error
}

type InitialMessageEventData struct{}

type ServerArenaMessageEventData struct {
	ArenaMessage ArenaMessage
}

type InitialMessageActionData struct {
	Name string `json:"name"`
}

type ServerArenaActionData struct {
	ArenaMessage ArenaMessage
}

type JoinArenaActionData struct {
	ArenaName uuid.UUID
}

// This function decodes the message type and then dispatches the
// correct handler based on the message type
func ServerDecodeAndDispatch(handler ServerHandler, rawData []byte) error {
	var raw struct {
		MessageType MessageType     `json:"message_type"`
		Data        json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	switch raw.MessageType {
	case InitialMessageActionType:
		initialMessageReturn := InitialMessageActionData{}
		err := json.Unmarshal(raw.Data, &initialMessageReturn)
		if err != nil {
			return err
		}
		return handler.HandleInitialMessageAction(initialMessageReturn)
	case ServerArenaActionType:
		serverArenaAction := ServerArenaActionData{}
		err := json.Unmarshal(raw.Data, &serverArenaAction)
		if err != nil {
			return err
		}
		return handler.HandleServerArenaAction(serverArenaAction)
	default:
		return fmt.Errorf("unexpected web.MessageType: %#v", raw.MessageType)
	}
}

func ClientDecodeAndDispatch(handler ClientHandler, rawData []byte) error {
	var raw struct {
		MessageType MessageType     `json:"message_type"`
		Data        json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	switch raw.MessageType {
	case ServerArenaEventType:
		arenaMessage := ServerArenaMessageEventData{}
		err := json.Unmarshal(raw.Data, &arenaMessage)
		if err != nil {
			return err
		}
		return handler.HandleServerArenaMessageEvent(arenaMessage)
	case InitialMessageEventType:
		initialMessage := InitialMessageEventData{}
		err := json.Unmarshal(raw.Data, &initialMessage)
		if err != nil {
			return err
		}
		return handler.HandleInitialMessageEvent(initialMessage)
	default:
		return fmt.Errorf("unexpected web.MessageType: %#v", raw.MessageType)
	}
}
