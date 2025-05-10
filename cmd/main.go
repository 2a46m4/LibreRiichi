package main

import (
	core "codeberg.org/ijnakashiar/LibreRiichi/core/arena"
	"github.com/google/uuid"
)

type ServerConfig struct {
	PortNumber uint16
}

func main() {
	rooms := map[uuid.UUID]core.Arena{}
}
