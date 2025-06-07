package web

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SetupHTTP() error {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/game", serveGame)
	http.ListenAndServe(":3000", nil)

	return nil
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	log.Println(r.Method)

	http.ServeFile(w, r, "static/home.html")
}

func serveGame(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	log.Println(r.Method)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panicln(err)
	}

	conn.WriteMessage(websocket.TextMessage, []byte("Hello world"))

	// http.ServeFile(w, r, "static/home.html")
}
