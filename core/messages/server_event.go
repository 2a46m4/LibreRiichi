package core

//go:generate go run ../generate_message.go -- ServerEvent

type ServerEvent interface {
	serverEventImpl()
}

type ServerArenaEvent struct {
	ArenaMessage ArenaEvent `json:"arena_message"` // wrap
}
