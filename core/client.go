package core

import (
	"github.com/google/uuid"
)

type Client struct {
	name       string
	id         uuid.UUID
	connection chan []byte

	room *Arena
}

func (client Client) Loop() {
	for {

	}
}
