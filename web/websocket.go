package web

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SetupHTTP(accept func(conn *websocket.Conn)) error {
	upgrader.CheckOrigin =
		func(r *http.Request) bool {
			return strings.HasPrefix(
				r.Header.Get("Origin"),
				"http://localhost",
			)
		}

	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		log.Println(r.Method)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		accept(conn)
	})
	http.ListenAndServe(":3000", nil)

	return nil
}
