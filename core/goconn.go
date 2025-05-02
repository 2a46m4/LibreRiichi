package core

import (
	"errors"
	"net"
	"os"
	"time"
)

// A channel wrapper for a connection
type ConnChan struct {
	DataChannel  chan any
	CloseChannel chan any
	WriteChannel chan []byte
	connection   net.Conn
}

// Create a ConnChan from a connection
func MakeChannel(conn net.Conn) ConnChan {
	ret := ConnChan{
		make(chan any, 0),
		make(chan any, 0),
		make(chan []byte, 0),
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

// To be called by the data receiver when it's done
func (conn ConnChan) CloseConnChan() {
	conn.CloseChannel <- Unit
	close(conn.CloseChannel)
	close(conn.WriteChannel)
}
