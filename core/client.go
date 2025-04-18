package core

import (
	"net"

	"github.com/google/uuid"
)

type Client struct {
	name       string
	id         uuid.UUID
	connection net.Conn

	room Arena
}

func (client Client) Loop() {
	for {

	}
}
