package core

type MessageType uint8

const (
	RESPONSE MessageType = iota
	REQUEST
	EVENT
)

type Message struct {
	MessageType  MessageType   `json:"message_type"`
	MessageIndex uint          `json:"message_index"`
	Data         ServerMessage `json:"data"`
}

func IsResponse(msgType MessageType) bool {
	return msgType >= MessageType(GENERICRESPONSE) && msgType <= MessageType(ARENAINFORESPONSE)
}

func IsEvent(msgType MessageType) bool {
	return msgType == MessageType(SERVERARENAEVENT)
}

func IsAction(msgType MessageType) bool {
	return msgType >= MessageType(INITIALMESSAGEACTION) && msgType <= MessageType(ARENAINFOACTION)
}

func Receive() {
	// Validate correctness of index
}

func Send() {
	// Pass in some kind of index
}
