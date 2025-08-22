package core

//go:generate go run ../generate_message.go -- ServerMessage

import (
	"time"
)

// The server message has three different types
type ServerMessage interface {
	serverMessageImpl()
}

type ServerArenaEvent struct {
	ArenaMessage ArenaEvent `json:"arena_message"` // wrap
}

type GenericResponse struct {
	Success    bool   `json:"success"`
	FailReason string `json:"fail_reason"`
}

type ListArenasResponse struct {
	Success   bool     `json:"success"`
	ArenaList []string `json:"arena_list"`
}

type ArenaInfoResponse struct {
	Success     bool        `json:"success"`
	Name        string      `json:"name"`
	Agents      []AgentInfo `json:"agents"`
	GameStarted bool        `json:"game_started"`
	DateCreated time.Time   `json:"date_created"`
}

type InitialMessageAction struct {
	Name string `json:"name"`
}

type JoinArenaAction struct {
	ArenaName string `json:"arena_name"`
}

type ServerArenaAction struct {
	ArenaAction ArenaAction `json:"arena_action"`
}

type ListArenasAction struct{}

type CreateArenaAction struct {
	ArenaName string `json:"arena_name"`
}

type ArenaInfoAction struct{}
