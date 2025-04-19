package core

import "net"

// A location where players gather. Controls the flow of the game
type Arena struct {
	Players []net.Conn
	Game    Game
}
