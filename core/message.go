package core

import (
	"encoding/json"
)

type MessageType uint8

const (
	PlayerJoinedEventType MessageType = iota
	GameStartedEventType
	SetupEventType
	GameEventType

	StartGameActionType
	PlayerActionType
	QuitActionType
)

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

type GameEventTypeData struct {
	ActionPerformed PlayerAction `json:"action"`
	// Whether this is an action that a player can take, not an action that a player took
	IsPotential bool       `json:"is_potential"`
	VisibleTo   Visibility `json:"visibility"`
}

func (GameEventTypeData) arenaMessageDataImpl() {}

type PlayerActionTypeData struct{}

func (PlayerActionTypeData) arenaMessageDataImpl() {}

type QuitActionTypeData struct{}

func (QuitActionTypeData) arenaMessageDataImpl() {}

type PlayerJoinedEventData struct {
	Agent Agent `json:"agent"`
}

func (PlayerJoinedEventData) arenaMessageDataImpl() {}

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
