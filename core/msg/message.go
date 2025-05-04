package msg

type MessageType uint8

const (
	PlayerJoinedEventType MessageType = iota
	GameStartedEventType

	StartGameActionType
	PlayerActionType
	QuitActionType
)

type Message struct {
	MessageType MessageType `json:"message_type"`
	Data        any
}
