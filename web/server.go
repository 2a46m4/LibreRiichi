package web

import (
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
		err := conn.WriteJSON(ServerMessage{
			MessageType: InitialMessage,
			Data:        nil,
		})
		if err != nil {
			fmt.Println("Failed at write")
			conn.Close()
			return
		}

		ret := ServerMessage{}
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
		err = ret.DecodeMessage(buffer)
		if err != nil {
			conn.Close()
			return
		}

		if ret.MessageType != InitialMessageReturn {
			conn.Close()
			return
		}
		res, ok := ret.Data.(*InitialMessageReturnData)
		if !ok {
			conn.Close()
			return
		}

		uuid, exist := server.Names[res.Name]
		if exist {
			fmt.Println("Not found")
			// Consider a failure message type here instead
			conn.Close()
			return
		}

		client, err := core.MakeClient(res.Name, server.Rooms[uuid], conn)
		if err != nil {
			fmt.Println("Client fail")
			conn.Close()
			return
		}

		go client.Loop()
	}()

}
