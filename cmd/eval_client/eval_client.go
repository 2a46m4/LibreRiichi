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

	var idx [4]uint = [4]uint{0, 0, 0, 0}

	// Initial hello
	{
		MakeMessage := func(name string, i int) []byte {
			msg := Message{
				MessageType:  REQUEST,
				MessageIndex: idx[i],
				Data: InitialMessageAction{
					Name: name,
				},
			}
			idx[i]++

			bytes, err := json.Marshal(msg)
			fmt.Println(string(bytes))
			if err != nil {
				panic(err)
			}
			return bytes
		}

		for i := range 4 {
			log.Println("Writing message")
			err = conns[i].WriteMessage(websocket.TextMessage, MakeMessage("Hello"+strconv.Itoa(i), i))
			if err != nil {
				log.Println("Write error:", err)
			}
		}

		// TODO: Assert index
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
		MakeMessage := func(i int) []byte {
			msg := Message{
				MessageType:  REQUEST,
				MessageIndex: idx[i],
				Data: CreateArenaAction{
					ArenaName: "arena",
				},
			}
			idx[i]++

			bytes, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}
			return bytes
		}

		log.Println("Sending create arena request")
		err = conns[0].WriteMessage(websocket.TextMessage, MakeMessage(0))
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

			MakeMessage := func(i int) []byte {
				msg := Message{
					MessageType:  REQUEST,
					MessageIndex: idx[i],
					Data: JoinArenaAction{
						ArenaName: "arena",
					},
				}
				idx[i]++

				bytes, err := json.Marshal(msg)
				if err != nil {
					panic(err)
				}
				return bytes
			}

			log.Println("Sending join arena request: ", i)
			err = conns[i].WriteMessage(websocket.TextMessage, MakeMessage(i))
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

		MakeMessage := func(i int) []byte {
			msg := Message{
				MessageType:  REQUEST,
				MessageIndex: idx[i],
				Data: ServerArenaAction{
					ArenaAction: StartGameActionData{},
				},
			}
			idx[i]++

			bytes, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}
			return bytes
		}

		log.Println("Sending start game request")
		err = conns[0].WriteMessage(websocket.TextMessage, MakeMessage(0))
		_, msg, err := conns[0].ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
		} else {
			log.Printf("Received: %s", msg)
		}
	}

	time.Sleep(10 * time.Second)

}
