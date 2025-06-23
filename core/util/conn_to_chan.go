package core

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// A channel wrapper for a connection
type ConnChan struct {
	// Data that will be received from the connection will be sent
	// through this channel
	DataChannel chan any
	// Receiving
	CloseChannel chan UnitType
	WriteChannel chan []byte
}

// Create a ConnChan from a connection
func MakeChannel(conn net.Conn) ConnChan {
	ret := ConnChan{
		make(chan any),
		make(chan UnitType),
		make(chan []byte),
	}

	go func() {
		buffer := make([]byte, 1024)
		for {
			select {
			case <-ret.CloseChannel:
				close(ret.DataChannel)
				err := conn.Close()
				if err != nil {
					fmt.Println(err)
				}
				return
			default:
				read, err := conn.Read(buffer)
				if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
					ret.DataChannel <- err
				} else {
					ret.DataChannel <- buffer[:read]
				}
			}
		}
	}()

	go func() {
		select {
		case <-ret.CloseChannel:
			return
		case toWrite := <-ret.WriteChannel:
			conn.SetDeadline(time.Now().Add(time.Second))
			_, err := conn.Write(toWrite)
			if err != nil && errors.Is(err, net.ErrClosed) {
				// TODO: Handle quit
				return
			}
		default:
		}
	}()

	return ret
}

// Create one from a WebSocket
// TODO: Make this more generic by having the user pass functions to handle the messages
func MakeChannelFromWebsocket(conn *websocket.Conn) ConnChan {
	ret := ConnChan{
		make(chan any),
		make(chan UnitType),
		make(chan []byte),
	}

	// Incoming channel
	go func() {
		for {
			select {
			case <-ret.CloseChannel:
				close(ret.DataChannel)
				err := conn.Close()
				if err != nil {
					fmt.Println(err)
				}
				return
			default:
				fmt.Println("Waiting for message")
				msgType, buffer, err := conn.ReadMessage()
				fmt.Println("Recved message", string(buffer))
				if err != nil {
					ret.DataChannel <- err
					continue
				}

				switch msgType {
				case websocket.TextMessage:
					ret.DataChannel <- buffer
				case websocket.BinaryMessage, websocket.PingMessage, websocket.PongMessage:
					continue
				case websocket.CloseMessage:
					close(ret.DataChannel)
					conn.Close()
					return
				}
			}
		}
	}()

	// Outgoing channel
	go func() {
		for {
			select {
			case <-ret.CloseChannel:
				return
			case toWrite, ok := <-ret.WriteChannel:
				if !ok {
					conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(0, "Closed conection"))
					return
				}

				err := conn.WriteMessage(websocket.TextMessage, toWrite)
				if err != nil {
					panic(err)
				}
			default:
			}
		}
	}()

	return ret
}

// Sends a message through the data channel
// TODO: Error handling
func (conn ConnChan) Send(data []byte) {
	conn.WriteChannel <- data
}

func (conn ConnChan) SendNonBlock(data []byte) bool {
	select {
	case conn.WriteChannel <- data:
		return true
	default:
		return false
	}
}

func (conn ConnChan) Close() {
	conn.CloseChannel <- Unit
}

// TODO: Error handling
func (conn ConnChan) Recv() any {
	return <-conn.DataChannel
}

func (conn ConnChan) RecvChan() <-chan any {
	return conn.DataChannel
}

func (conn ConnChan) RecvNonBlock() (any, bool) {
	select {
	case data := <-conn.DataChannel:
		return data, true
	default:
		return nil, false
	}
}

// Send the data to all of the channels
func Broadcast(data []byte, conns []ConnChan) {
	for _, conn := range conns {
		conn.WriteChannel <- data
	}
}

// To be called by the data receiver when it's done
func (conn ConnChan) CloseConnChan() {
	close(conn.CloseChannel)
	close(conn.WriteChannel)
}
