package core

import (
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
	return Client{
		Name:       name,
		ID:         uuid,
		Connection: MakeChannelFromWebsocket(connection),
		room:       room,
	}, nil
}

func (client Client) Loop() {
	fmt.Println(client.Name, client.ID, client.Connection, client.room)
	for {
		fmt.Println("Sending hello")
		client.Connection.Send([]byte("Hello"))
		fmt.Println("Sent")
	}
}
