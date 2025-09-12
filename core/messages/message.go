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

type WrongTypeError struct {
	expected MessageType
	was      MessageType
}

func (e WrongTypeError) Error() string {
	return fmt.Sprintf("Wrong type: expected %v, got %v", e.expected, e.was)
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

	m.MessageType = tempData.MessageType
	m.MessageIndex = tempData.MessageIndex
	switch tempData.MessageType {
	case EVENT:
		data := ServerEventUnpacker{}
		err := json.Unmarshal(tempData.Data, &data)
		if err != nil {
			return err
		}
		m.Data = data.ServerEvent
	case REQUEST:
		data := ServerActionUnpacker{}
		err := json.Unmarshal(tempData.Data, &data)
		if err != nil {
			return err
		}
		m.Data = data.ServerAction
	case RESPONSE:
		data := ServerResponseUnpacker{}
		err := json.Unmarshal(tempData.Data, &data)
		if err != nil {
			return err
		}
		m.Data = data.ServerResponse
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
		fmt.Println("Wrong index: ", msg)
		return nil, WrongIndexError{
			Wanted: index,
			Got:    msg.MessageIndex,
		}
	}

	if msg.MessageType != REQUEST {
		fmt.Println("Wrong message type: ", msg)
		return nil, WrongTypeError{
			expected: REQUEST,
			was:      msg.MessageType,
		}
	}
	return msg.Data.(ServerAction), nil
}

func SendEvent(event ServerEvent, index uint) {
	// Pass in some kind of index
}
