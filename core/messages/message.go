package core

import (
	"encoding/json"
	"errors"
	"fmt"
)

type MessageType uint8

const (
	RESPONSE MessageType = iota
	REQUEST
	EVENT
)

type Message struct {
	MessageType  MessageType `json:"message_type"`
	MessageIndex uint        `json:"message_index"`
	Data         any         `json:"data"`
}

type WrongIndexError struct {
	Wanted uint
	Got    uint
}

func (e WrongIndexError) Error() string {
	return fmt.Sprintf("Wrong index: wanted %v but got %v", e.Wanted, e.Got)
}

type WrongTypeError struct{}

func (e WrongTypeError) Error() string {
	return fmt.Sprintf("Wrong type")
}

func (m *Message) UnmarshalJSON(data []byte) error {
	tempData := struct {
		MessageType  MessageType     `json:"message_type"`
		MessageIndex uint            `json:"message_index"`
		Data         json.RawMessage `json:"data"`
	}{}

	err := json.Unmarshal(data, &tempData)
	if err != nil {
		return err
	}

	switch tempData.MessageType {
	case EVENT:
		data := ServerEventUnpacker{}
		json.Unmarshal(tempData.Data, &data)
		m.Data = data
	case REQUEST:
		data := ServerActionUnpacker{}
		json.Unmarshal(tempData.Data, &data)
		m.Data = data
	case RESPONSE:
		data := ServerResponseUnpacker{}
		json.Unmarshal(tempData.Data, &data)
		m.Data = data
	default:
		return errors.New("Bad message")
	}
	return nil
}

func ReceiveRequest(bytes []byte, index uint) (ServerAction, error) {
	msg := Message{}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return nil, err
	}
	if index != msg.MessageIndex {
		return nil, WrongIndexError{
			Wanted: index,
			Got:    msg.MessageIndex,
		}
	}

	if msg.MessageType != REQUEST {
		return nil, WrongTypeError{}
	}
	return msg.Data.(ServerAction), nil
}

func SendEvent(event ServerEvent, index uint) {
	// Pass in some kind of index
}
