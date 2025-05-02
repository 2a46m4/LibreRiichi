package msg

type MessageType uint8

const (
	PlayerJoinedEvent MessageType = iota
	GameStartedEvent

	StartGameAction
	QuitAction
)

type Message struct {
	MessageType MessageType `json:"message_type"`
	Data        any
}
