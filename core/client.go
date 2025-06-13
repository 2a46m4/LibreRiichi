package core

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Name       string    `json:"name"`
	ID         uuid.UUID `json:"id"`
	Connection ConnChan  `json:"-"`
	Recv       chan Message
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

			err := ServerDecodeAndDispatch(&client, recv.([]byte))
			if err != nil {
				continue
			}
		}
	}
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
	return nil
}

func (client Client) GetSendChannel() chan<- Message {
	return client.Recv
}
