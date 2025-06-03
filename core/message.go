package core

import (
	"encoding/json"
)

type MessageType uint8

const (
	// Messages that are sent from game to player
	PlayerJoinedEventType MessageType = iota
	GameStartedEventType
	SetupEventType
	PlayerActionEventType

	// Messages that are sent from player to game
	StartGameActionType
	PlayerActionType
	QuitActionType
)

// ArenaMessage are messages that are sent between clients and server
type ArenaMessage struct {
	MessageType MessageType      `json:"message_type"`
	Data        ArenaMessageData `json:"data"`
	VisibleTo   Visibility       `json:"visibility"`
}

type ArenaMessageData interface{ arenaMessageDataImpl() }

type PlayerJoinedEventTypeData struct{}

func (PlayerJoinedEventTypeData) arenaMessageDataImpl() {}

type GameStartedEventTypeData struct{}

func (GameStartedEventTypeData) arenaMessageDataImpl() {}

type SetupEventTypeData struct {
	Setup Setup `json:"setup"`
}

func (SetupEventTypeData) arenaMessageDataImpl() {}

type StartGameActionTypeData struct{}

func (StartGameActionTypeData) arenaMessageDataImpl() {}

// TODO: Change ActionResult to be this type
type PlayerActionEventTypeData struct {
	ActionResult
}

func (PlayerActionEventTypeData) arenaMessageDataImpl() {}

type PlayerActionTypeData struct{}

func (PlayerActionTypeData) arenaMessageDataImpl() {}

type QuitActionTypeData struct{}

func (QuitActionTypeData) arenaMessageDataImpl() {}

type PlayerJoinedEventData struct {
	Agent Agent `json:"agent"`
}

func (PlayerJoinedEventData) arenaMessageDataImpl() {}

// Selects the correct type to write into
func arenaToDataMap(msgType MessageType) ArenaMessageData {
	var ArenaToDataMap = []ArenaMessageData{
		PlayerJoinedEventData{},
		GameStartedEventTypeData{},
		SetupEventTypeData{},
		StartGameActionTypeData{},
		PlayerActionTypeData{},
		QuitActionTypeData{},
	}

	return ArenaToDataMap[msgType]
}

func (arena *ArenaMessage) DecodeArenaMessage(rawData []byte) error {
	var raw struct {
		MessageType MessageType     `json:"message_type"`
		Data        json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	arena.MessageType = raw.MessageType
	data := arenaToDataMap(raw.MessageType)
	if err := json.Unmarshal(raw.Data, &data); err != nil {
		return err
	}
	arena.Data = data
	return nil
}
