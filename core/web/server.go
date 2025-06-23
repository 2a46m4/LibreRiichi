package web

import (
	"fmt"

	core "codeberg.org/ijnakashiar/LibreRiichi/core"
	util "codeberg.org/ijnakashiar/LibreRiichi/core/util"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Server struct {
	Names map[string]uuid.UUID
	Rooms map[uuid.UUID]*core.Arena

	ServerConfig struct {
		PortNumber uint16
	}
}

func (server Server) AcceptConnection(conn *websocket.Conn) {
	fmt.Println("Got connection")
	go func() {
		client, err := core.MakeClient(util.MakeChannelFromWebsocket(conn))
		if err != nil {
			fmt.Println("Client fail")
			conn.Close()
			return
		}

		go client.Loop()
	}()
}
