package core

import (
	"encoding/json"
	"net"
	"time"
)

func Read(conn net.Conn, v *any) error {
	err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return err
	}

	bytes := make([]byte, 256)
	_, err = conn.Read(bytes)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		return err
	}

	return nil
}

func Write(conn net.Conn, v *any) error {
	err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	n, err := conn.Write(bytes)
	if err != nil || n != len(bytes) {
		return err
	}

	return nil
}
