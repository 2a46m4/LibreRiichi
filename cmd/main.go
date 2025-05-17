package main

import (
	"fmt"

	core "codeberg.org/ijnakashiar/LibreRiichi/core"
	"github.com/google/uuid"
)

type ServerConfig struct {
	PortNumber uint16
}

func main() {
	rooms := map[uuid.UUID]core.Arena{}
	fmt.Println(rooms)
}
