package core

import (
	"encoding/json"
	"fmt"
)

type MessageType uint8

const (
	RESPONSE MessageType = iota
	REQUEST
	EVENT
)

type Message struct {
	MessageType  MessageType           `json:"message_type"`
	MessageIndex uint                  `json:"message_index"`
	Data         ServerMessageUnpacker `json:"data"`
}

type WrongIndexError struct {
	Wanted uint
	Got    uint
}

func (e WrongIndexError) Error() string {
	return fmt.Sprintf("Wrong index: wanted %v but got %v", e.Wanted, e.Got)
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

// Gets a message and validates it
func Receive(bytes []byte, index uint) (ServerMessage, error) {
	msg := Message{}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return nil, err
	}
	if index != msg.MessageIndex {
		fmt.Println("Wrong index")
		return nil, WrongIndexError{
			Wanted: index,
			Got:    msg.MessageIndex,
		}
	}

	return msg.Data, nil
}

func Send() {
	// Pass in some kind of index
}
