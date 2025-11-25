package main

import (
	core "codeberg.org/ijnakashiar/LibreRiichi/core"
	web "codeberg.org/ijnakashiar/LibreRiichi/core/web"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
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
