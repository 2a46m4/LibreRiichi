package web

import (
	"encoding/json"

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

type ServerMessage struct {
	MessageType MessageType       `json:"message_type"`
	Data        ServerMessageData `json:"data"`
}

type ServerMessageData interface{ serverMessageDataImpl() }

type InitialMessageData struct{}

func (InitialMessageData) serverMessageDataImpl() {}

type ArenaMessageData struct {
	ArenaMessage core.ArenaMessage
}

func (ArenaMessageData) serverMessageDataImpl() {}

type InitialMessageReturnData struct {
	Name string `json:"name"`
	Room string `json:"room"`
}

func (InitialMessageReturnData) serverMessageDataImpl() {}

func (msg *ServerMessage) DecodeMessage(rawData []byte) error {
	var raw struct {
		MessageType MessageType     `json:"message_type"`
		Data        json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	var MessageToDataMap = []ServerMessageData{
		&InitialMessageData{},
		&ArenaMessageData{},
		&InitialMessageReturnData{},
	}

	msg.MessageType = raw.MessageType
	data := MessageToDataMap[raw.MessageType]
	if err := json.Unmarshal(raw.Data, data); err != nil {
		return err
	}
	msg.Data = data
	return nil
}
