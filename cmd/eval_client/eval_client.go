package main

import (
	"github.com/gorilla/websocket"
	"log"
)

func main() {
	url := "ws://localhost:3000/game"
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Dial failed: %v (status: %v)", err, resp)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello WebSocket!"))
	if err != nil {
		log.Println("Write error:", err)
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read error:", err)
	} else {
		log.Printf("Received: %s", msg)
	}
}
