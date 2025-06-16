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

type BadTypeError struct{}

func (BadTypeError) Error() string {
	return "Bad type"
}

type Message struct {
	MessageType MessageType `json:"message_type"`
	Data        any         `json:"data"`
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
	case ServerArenaEventType:
		arenaEventMessage := ServerArenaMessageEventData{}
		err := json.Unmarshal(raw.Data, &arenaEventMessage)
		if err != nil {
			return err
		}
		msg.Data = arenaEventMessage
	case InitialMessageEventType:
		initialMessageEvent := InitialMessageEventData{}
		err := json.Unmarshal(raw.Data, &initialMessageEvent)
		if err != nil {
			return err
		}
		msg.Data = initialMessageEvent
	default:
		return fmt.Errorf("unexpected web.MessageType: %#v", raw.MessageType)
	}

	return nil
}

// This function decodes the message type and then dispatches the
// correct handler based on the message type
func ServerDispatch(handler ServerHandler, message Message) error {
	switch message.MessageType {
	case InitialMessageActionType:
		initialMessageReturn, ok := message.Data.(InitialMessageActionData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandleInitialMessageAction(initialMessageReturn)
	case ServerArenaActionType:
		serverArenaAction, ok := message.Data.(ServerArenaActionData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandleServerArenaAction(serverArenaAction)
	default:
		return fmt.Errorf("unexpected web.MessageType: %#v", message.MessageType)
	}
}

func ClientDispatch(handler ClientHandler, message Message) error {
	switch message.MessageType {
	case ServerArenaEventType:
		arenaMessage, ok := message.Data.(ServerArenaMessageEventData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandleServerArenaMessageEvent(arenaMessage)
	case InitialMessageEventType:
		initialMessage, ok := message.Data.(InitialMessageEventData)
		if !ok {
			return BadTypeError{}
		}
		return handler.HandleInitialMessageEvent(initialMessage)
	default:
		return fmt.Errorf("unexpected web.MessageType: %#v", message.MessageType)
	}
}
