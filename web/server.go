package web

import (
	"encoding/json"
	"fmt"

	"codeberg.org/ijnakashiar/LibreRiichi/core"
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
		err := conn.WriteJSON(core.Message{
			MessageType: core.InitialMessageEventType,
			Data:        nil,
		})
		if err != nil {
			fmt.Println("Failed at write")
			conn.Close()
			return
		}

		ret := core.Message{}
		messageType, buffer, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Failed at read")
			conn.Close()
			return
		}
		if messageType != websocket.TextMessage {
			fmt.Println("Wrong type")
			conn.Close()
			return
		}
		fmt.Println(string(buffer))
		err = json.Unmarshal(buffer, &ret)
		if err != nil {
			fmt.Println("Unmarshal failure")
			conn.Close()
			return
		}

		if ret.MessageType != core.InitialMessageActionType {
			fmt.Println("Unexpected message")
			conn.Close()
			return
		}
		res, ok := ret.Data.(core.InitialMessageActionData)
		if !ok {
			conn.Close()
			return
		}

		client, err := core.MakeClient(res.Name, conn)
		if err != nil {
			fmt.Println("Client fail")
			conn.Close()
			return
		}

		go client.Loop()
	}()

}
