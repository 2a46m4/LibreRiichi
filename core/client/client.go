package core

import (
	arena "codeberg.org/ijnakashiar/LibreRiichi/core/arena"

	"github.com/google/uuid"
)

type Client struct {
	name       string
	id         uuid.UUID
	connection chan []byte

	room *arena.Arena
}

func (client Client) Loop() {
	for {

	}
}
