package web

import (
	"encoding/json"
	"fmt"

	"codeberg.org/ijnakashiar/LibreRiichi/core"
)

type MessageType uint8

const (
	// Messages that are sent from server to client
	InitialMessage MessageType = iota
	ArenaMessage

	// Messages that are sent from client to server
	InitialMessageReturn
)

type Message struct {
	MessageType MessageType     `json:"message_type"`
	Data        json.RawMessage `json:"data"`
}

// This interface handles different message types
type ServerHandler interface {
	HandlerInitialMessageReturn(InitialMessageReturnData) error
}

type ClientHandler interface {
	HandleInitialMessage(InitialMessageData) error
	HandleArenaMessage(ArenaMessageData) error
}

type InitialMessageData struct{}

type ArenaMessageData struct {
	ArenaMessage core.ArenaMessage
}

type InitialMessageReturnData struct {
	Name string `json:"name"`
	Room string `json:"room"`
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
	case InitialMessageReturn:
		initialMessageReturn := InitialMessageReturnData{}
		err := json.Unmarshal(raw.Data, &initialMessageReturn)
		if err != nil {
			return err
		}
		return handler.HandlerInitialMessageReturn(initialMessageReturn)
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
	case ArenaMessage:
		arenaMessage := ArenaMessageData{}
		err := json.Unmarshal(raw.Data, &arenaMessage)
		if err != nil {
			return err
		}
		return handler.HandleArenaMessage(arenaMessage)
	case InitialMessage:
		initialMessage := InitialMessageData{}
		err := json.Unmarshal(raw.Data, &initialMessage)
		if err != nil {
			return err
		}
		return handler.HandleInitialMessage(initialMessage)
	default:
		return fmt.Errorf("unexpected web.MessageType: %#v", raw.MessageType)
	}
}
