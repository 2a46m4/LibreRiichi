package core

//go:generate go run ../generate_message.go -- ArenaAction

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type ArenaAction interface {
	ArenaActionImpl()
}

type StartGameActionData struct{}

type PlayerQuitActionData struct{}

type PlayerActionData struct {
	Action Action // wrap
}
