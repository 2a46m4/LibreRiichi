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

func MakeClient(connection *websocket.Conn) (Client, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return Client{}, err
	}

	client := Client{
		Name:       uuid.String(),
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
				panic(err)
			}

			msg := Message{}
			err := json.Unmarshal(recv.([]byte), &msg)
			if err != nil {
				continue
			}
			ServerDispatch(&client, msg)
		}
	}
}

// HandleJoinArenaAction implements ServerHandler.
func (client *Client) HandleJoinArenaAction(data JoinArenaActionData) error {
	if client.Arena != nil {
		return errors.New("Already in an arena")
	}

	arena := GetArena(data.ArenaName)
	err := arena.JoinArena(client, true)
	if err != nil {
		return err
	}
	client.Arena = arena
	return nil
}

func (client *Client) HandleInitialMessageAction(data InitialMessageActionData) error {
	client.Name = data.Name
	bytes, err := json.Marshal(InitialMessageEventData{})
	if err != nil {
		return err
	}
	client.Connection.Send(bytes)
	return nil
}

func (client *Client) HandleServerArenaAction(action ServerArenaActionData) error {
	if client.Arena != nil {
		return ServerArenaDispatch(client.Arena, action.ArenaMessage)
	}
	return errors.New("No arena found")
}

func (client Client) GetSendChannel() chan<- Message {
	return client.Recv
}
