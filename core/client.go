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

	room *Arena `json:"-"`
}

func MakeClient(name string, room *Arena, connection *websocket.Conn) (Client, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return Client{}, err
	}
	client := Client{
		Name:       name,
		ID:         uuid,
		Connection: MakeChannelFromWebsocket(connection),
		room:       room,
	}
	fmt.Println("Making new client", client)
	return client, nil
}

func (client Client) Loop() {
	fmt.Println(client.Name, client.ID, client.Connection, client.room)
	for {
		msg, err := json.Marshal("Hello")
		if err != nil {
			panic(err)
		}
		sent := client.Connection.SendNonBlock(msg)
		if sent {
			fmt.Println("Sent")
		}

		data, recv := client.Connection.RecvNonBlock()
		if recv {
			fmt.Println("Data received by client:", string(data.([]byte)))
		}
	}
}
