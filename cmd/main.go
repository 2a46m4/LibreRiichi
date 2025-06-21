package main

import (
	"os"
	"os/signal"
	"syscall"

	core "codeberg.org/ijnakashiar/LibreRiichi/core"
	web "codeberg.org/ijnakashiar/LibreRiichi/core/web"
	"github.com/google/uuid"
)

func main() {
	server := web.Server{
		Rooms:        map[uuid.UUID]*core.Arena{},
		ServerConfig: struct{ PortNumber uint16 }{3000},
	}

	web.SetupHTTP(server.AcceptConnection)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
}
