package msg

import "encoding/json"

type MessageType uint8

const (
	PlayerJoinedEventType MessageType = iota
	GameStartedEventType

	StartGameActionType
	PlayerActionType
	QuitActionType
)

type ArenaMessage struct {
	MessageType MessageType      `json:"message_type"`
	Data        ArenaMessageData `json:"data"`
}

type ArenaMessageData interface{ arenaMessageDataImpl() }

type PlayerJoinedEventData struct{}

func (arena *ArenaMessage) DecodeArenaMessage(rawData []byte) {
	var raw struct {
		Action MessageType     `json:"action_type"`
		Data   json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	var ActionToDataMap = []ActionData{
		RonData{},
		TsumoData{},
		RiichiData{},
		TossData{},
		SkipData{},
		PonData{},
		KanData{},
		ChiiData{},
		DrawData{},
	}

	action.Action = raw.Action
	data := ActionToDataMap[raw.Action]
	if err := json.Unmarshal(raw.Data, &data); err != nil {
		return err
	}
	action.Data = data
	return nil
}
