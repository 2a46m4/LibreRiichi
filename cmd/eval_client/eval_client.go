package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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
				time.Sleep(100 * time.Millisecond)
			} else {
				conns[i] = conn
				fmt.Println("New conn")
				break
			}
		}
	}

	log.Println("Finished connecting")

	for i := range 4 {
		defer conns[i].Close()
	}

	// Initial hello
	{
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
			err = conns[i].WriteMessage(websocket.TextMessage, MakeMessage("Hello"+strconv.Itoa(i)))
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

	// Create a room
	{
		MakeMessage := func() []byte {
			msg := Message{
				MessageType:  CreateArenaActionType,
				MessageIndex: 0,
				Data: CreateArenaActionData{
					ArenaName: "arena",
				},
			}
			bytes, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}
			return bytes
		}

		log.Println("Sending create arena request")
		err = conns[0].WriteMessage(websocket.TextMessage, MakeMessage())
		_, msg, err := conns[0].ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
		} else {
			log.Printf("Received: %s", msg)
		}
	}

	// Join rooms
	{
		for i := range 4 {

			MakeMessage := func() []byte {
				msg := Message{
					MessageType:  JoinArenaActionType,
					MessageIndex: 0,
					Data: JoinArenaActionData{
						ArenaName: "arena",
					},
				}
				bytes, err := json.Marshal(msg)
				if err != nil {
					panic(err)
				}
				return bytes
			}

			log.Println("Sending join arena request: ", i)
			err = conns[i].WriteMessage(websocket.TextMessage, MakeMessage())
			_, msg, err := conns[i].ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
			} else {
				log.Printf("Received: %s", msg)
			}
		}
	}

	// Start game
	{

		MakeMessage := func() []byte {
			msg := Message{
				MessageType:  ServerArenaActionType,
				MessageIndex: 0,
				Data: ServerArenaActionData{
					ArenaMessage: ArenaMessage{
						MessageType: StartGameActionType,
						Data:        StartGameActionData{},
					},
				},
			}
			bytes, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}
			return bytes
		}

		log.Println("Sending start game request")
		err = conns[0].WriteMessage(websocket.TextMessage, MakeMessage())
		_, msg, err := conns[0].ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
		} else {
			log.Printf("Received: %s", msg)
		}
	}

	time.Sleep(10 * time.Second)

}
