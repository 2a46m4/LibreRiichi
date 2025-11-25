package core

import "time"

//go:generate go run ../generate_message.go -- ServerResponse

type ServerResponse interface {
	serverResponseImpl()
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

type GameInfoResponse struct {
}
