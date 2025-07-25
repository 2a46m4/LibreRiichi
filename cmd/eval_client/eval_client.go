package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"

	"github.com/gorilla/websocket"
)

func main() {

	url := "ws://localhost:3000/game"
	var err error = nil
	var conn *websocket.Conn
	var resp *http.Response
	conns := [4]*websocket.Conn{}
	for i := range 4 {
		log.Println("i:", i)
		for {
			header := http.Header{}
			header.Set("Origin", "http://localhost")
			conn, resp, err = websocket.DefaultDialer.Dial(url, header)
			if err != nil {
				log.Printf("Dial failed: %v (status: %v)", err, resp)
			} else {
				conns[i] = conn
				fmt.Println("New conn")
				break
			}
		}
	}

	log.Println("Finished connecting")

	defer conn.Close()

	MakeMessage := func(name string) []byte {
		msg := Message{
			MessageType:  InitialMessageActionType,
			MessageIndex: 0,
			Data: InitialMessageActionData{
				Name: name,
			},
		}
		bytes, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}
		return bytes
	}

	for i := range 4 {
		log.Println("Writing message")
		err = conns[i].WriteMessage(websocket.TextMessage, MakeMessage("Hello"+string(i)))
		log.Println("Write error:", err)
	}

	for i := range 4 {
		_, msg, err := conns[i].ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
		} else {
			log.Printf("Received: %s", msg)
		}
	}

}
