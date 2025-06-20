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
	ListArenasEventType
	CreateArenaEventType

	// Messages that are sent from client to server
	InitialMessageActionType
	JoinArenaActionType
	ServerArenaActionType
	ListArenasActionType
	CreateArenaActionType
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
type ServerHandler[Return any] interface {
	HandleInitialMessageAction(InitialMessageActionData) (Return, error)
	HandleJoinArenaAction(JoinArenaActionData) (Return, error)
	HandleServerArenaAction(ServerArenaActionData) (Return, error)
	HandleCreateArenaAction(CreateArenaActionData) (Return, error)
}

type ClientHandler[Return any] interface {
	HandleInitialMessageEvent(InitialMessageEventData) (Return, error)
	HandleJoinArenaEvent() (Return, error)
	HandleServerArenaMessageEvent(ServerArenaMessageEventData) (Return, error)
}

type InitialMessageEventData struct{}

type ServerArenaMessageEventData struct {
	ArenaMessage ArenaMessage `json:"arena_message"`
}

type JoinArenaEventData struct {
	Success bool `json:"success"`
}

type ListArenasEventData struct {
	ArenaList []uuid.UUID `json:"arena_list"`
}

type CreateArenaEventData struct {
	Success bool `json:"success"`
}

type InitialMessageActionData struct {
	Name string `json:"name"`
}

type ServerArenaActionData struct {
	ArenaMessage ArenaMessage `json:"arena_message"`
}

type JoinArenaActionData struct {
	ArenaName string    `json:"arena_name"`
	ArenaID   uuid.UUID `json:"arena_id"`
}

type ListArenasActionData struct{}

type CreateArenaActionData struct {
	ArenaName string `json:"arena_name"`
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
	default:
		return fmt.Errorf("unexpected web.MessageType: %#v", raw.MessageType)
	}

	return nil
}

// This function decodes the message type and then dispatches the
// correct handler based on the message type
func ServerDispatch[Return any](handler ServerHandler[Return], message Message) (ret Return, err error) {
	switch message.MessageType {
	case InitialMessageActionType:
		initialMessageReturn, ok := message.Data.(InitialMessageActionData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleInitialMessageAction(initialMessageReturn)
	case ServerArenaActionType:
		serverArenaAction, ok := message.Data.(ServerArenaActionData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleServerArenaAction(serverArenaAction)
	case CreateArenaActionType:
		data, ok := message.Data.(CreateArenaActionData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleCreateArenaAction(data)
	default:
		return ret, fmt.Errorf("unexpected web.MessageType: %#v", message.MessageType)
	}
}

func ClientDispatch[Return any](handler ClientHandler[Return], message Message) (ret Return, err error) {
	switch message.MessageType {
	case ServerArenaEventType:
		arenaMessage, ok := message.Data.(ServerArenaMessageEventData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleServerArenaMessageEvent(arenaMessage)
	case InitialMessageEventType:
		initialMessage, ok := message.Data.(InitialMessageEventData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleInitialMessageEvent(initialMessage)
	default:
		return ret, fmt.Errorf("unexpected web.MessageType: %#v", message.MessageType)
	}
}
