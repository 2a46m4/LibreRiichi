package core

import (
	"errors"
	"net"
	"os"
	"time"
)

// A channel wrapper for a connection
type ConnChan struct {
	// Data that will be received
	DataChannel  chan any
	CloseChannel chan any
	WriteChannel chan []byte
	connection   net.Conn
}

// Create a ConnChan from a connection
func MakeChannel(conn net.Conn) ConnChan {
	ret := ConnChan{
		make(chan any),
		make(chan any),
		make(chan []byte),
		conn,
	}

	go func() {
		buffer := make([]byte, 1024)
		for {
			conn.SetDeadline(time.Now().Add(time.Second))
			read, err := conn.Read(buffer)
			if errors.Is(err, os.ErrDeadlineExceeded) {
				continue
			} else if err != nil {
				ret.DataChannel <- err
			}

			select {
			case <-ret.CloseChannel:
				close(ret.DataChannel)
				conn.Close()
				return
			case toWrite := <-ret.WriteChannel:
				conn.SetDeadline(time.Now().Add(time.Second))
				conn.Write(toWrite)
			default:
			}

			ret.DataChannel <- buffer[:read]

		}

	}()

	return ret
}

// Sends a message through the data channel
func (conn ConnChan) Send(data []byte) {
	conn.WriteChannel <- data
}

func (conn ConnChan) Close() {
	conn.CloseChannel <- Unit
}

func (conn ConnChan) Recv() any {
	return <-conn.DataChannel
}

// Send the data to all of the channels
func Broadcast(data []byte, conns []ConnChan) {
	for _, conn := range conns {
		conn.WriteChannel <- data
	}
}

// To be called by the data receiver when it's done
func (conn ConnChan) CloseConnChan() {
	conn.CloseChannel <- Unit
	close(conn.CloseChannel)
	close(conn.WriteChannel)
}
