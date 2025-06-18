package core

import "github.com/google/uuid"

// TODO: Temporary
type ArenaList struct {
	arena Arena
}

var GlobalArenaList ArenaList = ArenaList{}

func InitializeArenaList() {

}

func GetArena(uuid.UUID) *Arena {
	return &GlobalArenaList.arena
}
