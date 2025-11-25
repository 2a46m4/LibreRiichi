package main

import (
    _ "net/http/pprof"
	"log"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	core "codeberg.org/ijnakashiar/LibreRiichi/core"
	web "codeberg.org/ijnakashiar/LibreRiichi/core/web"
)

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	server := web.Server{
		Rooms:        &core.GlobalArenaList,
		ServerConfig: struct{ PortNumber uint16 }{3000},
	}

	core.InitializeMap()
	web.SetupHTTP(server.AcceptConnection)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
}
