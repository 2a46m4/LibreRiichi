package core

import "github.com/google/uuid"

type AgentInfo struct {
	Name  string    `json:"name"`
	ID    uuid.UUID `json:"id"`
	Order uint8     `json:"order"`
}
