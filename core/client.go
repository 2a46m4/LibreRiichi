package core

import "net"

type Client struct {
	connection net.Conn
}

func (client Client) Loop() {

}
