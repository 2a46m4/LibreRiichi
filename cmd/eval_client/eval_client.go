package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {

	url := "ws://localhost:3000/game"
	var err error = errors.New("Hello")
	var conn *websocket.Conn
	var resp *http.Response
	conns := [4]*websocket.Conn{}
	for i := range 4 {
		header := http.Header{}
		header.Set("Origin", "http://localhost")
		conn, resp, err = websocket.DefaultDialer.Dial(url, header)
		if err != nil {
			log.Printf("Dial failed: %v (status: %v)", err, resp)
		} else {
			conns[i] = conn
			fmt.Println("New conn")
		}
	}

	defer conn.Close()

	for err != nil {
		err = conn.WriteMessage(websocket.TextMessage, []byte("Hello WebSocket!"))
		log.Println("Write error:", err)
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read error:", err)
	} else {
		log.Printf("Received: %s", msg)
	}
}
