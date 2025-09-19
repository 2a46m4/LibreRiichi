package core

import "github.com/google/uuid"

//go:generate go run ../generate_message.go -- ArenaEvent

type ArenaEvent interface {
	ArenaEventImpl()
}

type PlayerJoinedEvent struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

type PlayerQuitEvent struct {
	Name string `json:"name"`
}

type GameStartedEvent struct{}

type ArenaBoardEvent struct {
	// For handling generic games, this should be replaced
	BoardEvent BoardEvent `json:"board_event"` // wrap
}
