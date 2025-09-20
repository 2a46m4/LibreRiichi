package core

//go:generate go run ../generate_message.go -- ServerResponse

type ServerAction interface {
	serverActionImpl()
}

type InitialMessageAction struct {
	Name string `json:"name"`
}

type JoinArenaAction struct {
	ArenaName string `json:"arena_name"`
}

type ServerArenaAction struct {
	ArenaAction ArenaAction `json:"arena_action"` // wrap
}

type ListArenasAction struct{}

type CreateArenaAction struct {
	ArenaName string `json:"arena_name"`
}

type ArenaInfoAction struct{}

type AddAIArenaAction struct{}

type RemoveAIArenaAction struct{}