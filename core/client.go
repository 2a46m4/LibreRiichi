package core

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Name       string       `json:"name"`
	ID         uuid.UUID    `json:"id"`
	Connection ConnChan     `json:"-"`
	Recv       chan Message `json:"-"`
	Arena      *Arena       `json:"-"`
}

type DispatchResult struct {
	Message Message
	DoSend  bool
}

func NoSend() DispatchResult {
	return DispatchResult{
		Message: Message{},
		DoSend:  false,
	}
}

func FormatMessage(msgType MessageType, data any) DispatchResult {
	return DispatchResult{
		Message: Message{
			MessageType: msgType,
			Data:        data,
		},
		DoSend: true,
	}
}

func MakeClient(name string, connection *websocket.Conn) (Client, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return Client{}, err
	}

	client := Client{
		Name:       name,
		ID:         uuid,
		Connection: MakeChannelFromWebsocket(connection),
		Recv:       make(chan Message, 32),
		Arena:      nil,
	}
	fmt.Println("Making new client", client)
	return client, nil
}

func (client Client) Loop() {
	fmt.Println(client.Name, client.ID, client.Connection)
	for {
		select {
		case send := <-client.Recv:
			bytes, err := json.Marshal(send)
			if err != nil {
				panic(err)
			}
			client.Connection.Send(bytes)
		case recv := <-client.Connection.RecvChan():
			if err, ok := recv.(error); ok {
				fmt.Println(err)
				client.Connection.CloseConnChan()
				return
			}

			msg := Message{}
			err := json.Unmarshal(recv.([]byte), &msg)
			if err != nil {
				fmt.Println("Error unmarshalling:", err)
				continue
			}
			dispatchResult, err := ServerDispatch(&client, msg)
			if err != nil {
				fmt.Println("Problem with message during dispatch:", err)
				continue
			}

			if dispatchResult.DoSend {
				client.GetSendChannel() <- dispatchResult.Message
			}
		}
	}
}

// HandleJoinArenaAction implements ServerHandler.
func (client *Client) HandleJoinArenaAction(data JoinArenaActionData) (DispatchResult, error) {
	if client.Arena != nil {
		return FormatMessage(JoinArenaEventType,
			JoinArenaEventData{
				Success: false,
			}), errors.New("Already in an arena")
	}

	arena, err := GetArenaFromName(data.ArenaName)
	if err != nil {
		return FormatMessage(JoinArenaEventType,
			JoinArenaEventData{
				Success: false,
			}), err
	}

	err = arena.JoinArena(client, true)
	if err != nil {
		return FormatMessage(JoinArenaEventType,
			JoinArenaEventData{
				Success: false,
			}), err
	}

	client.Arena = arena
	return FormatMessage(JoinArenaEventType,
		JoinArenaEventData{
			Success: false,
		}), nil
}

func (client *Client) HandleInitialMessageAction(data InitialMessageActionData) (DispatchResult, error) {
	client.Name = data.Name
	return FormatMessage(InitialMessageEventType,
		InitialMessageEventData{}), nil
}

func (client *Client) HandleServerArenaAction(action ServerArenaActionData) (DispatchResult, error) {
	if client.Arena != nil {
		// TODO: Do something with this result
		ServerArenaDispatch(client.Arena, action.ArenaMessage)
		return DispatchResult{}, nil
	}
	return DispatchResult{}, errors.New("No arena found")
}

func (client *Client) HandleCreateArenaAction(data CreateArenaActionData) (DispatchResult, error) {
	err := CreateAndAddArena(data.ArenaName)
	if err != nil {
		return DispatchResult{}, err
	}
	return FormatMessage(CreateArenaEventType, CreateArenaEventData{
		Success: true,
	}), nil
}

func (client Client) GetSendChannel() chan<- Message {
	return client.Recv
}
