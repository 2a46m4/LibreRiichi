package core

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type ArenaList struct {
	arena map[uuid.UUID]*Arena
	name  map[string]uuid.UUID

	sync.RWMutex
}

var GlobalArenaList ArenaList = ArenaList{}

type ArenaNotFoundError struct {
	arena_name string
}

type SameNameError struct {
	arena_name string
}

func InitializeMap() {
	GlobalArenaList.arena = make(map[uuid.UUID]*Arena)
	GlobalArenaList.name = make(map[string]uuid.UUID)
}

func (e ArenaNotFoundError) Error() string {
	return fmt.Sprintf("Arena %v is not found", e.arena_name)
}

func (e SameNameError) Error() string {
	return fmt.Sprintf("There is already an arena with the same name: %v", e.arena_name)
}

func GetArenaFromName(name string) (*Arena, error) {
	GlobalArenaList.RLock()
	defer GlobalArenaList.RUnlock()

	uuid, ok := GlobalArenaList.name[name]
	if !ok {
		return nil, ArenaNotFoundError{name}
	}

	arena, ok := GlobalArenaList.arena[uuid]
	if !ok {
		return nil, ArenaNotFoundError{name}
	}

	return arena, nil
}

func GetArena(uuid uuid.UUID) (*Arena, error) {
	GlobalArenaList.RLock()
	defer GlobalArenaList.RUnlock()

	arena, ok := GlobalArenaList.arena[uuid]
	if !ok {
		return nil, ArenaNotFoundError{uuid.String()}
	}

	return arena, nil
}

func RemoveArena(name string) error {
	GlobalArenaList.Lock()
	defer GlobalArenaList.Unlock()

	uuid, ok := GlobalArenaList.name[name]
	if !ok {
		return ArenaNotFoundError{name}
	}

	delete(GlobalArenaList.name, name)
	delete(GlobalArenaList.arena, uuid)
	return nil
}

func RemoveArenaUUID(uuid uuid.UUID) {
	GlobalArenaList.Lock()
	defer GlobalArenaList.Unlock()
	delete(GlobalArenaList.arena, uuid)
}

func GetArenaUUID(name string) (uuid.UUID, error) {
	GlobalArenaList.RLock()
	defer GlobalArenaList.RUnlock()

	uuid, ok := GlobalArenaList.name[name]
	if !ok {
		return uuid, ArenaNotFoundError{name}
	}

	return uuid, nil
}

func CreateAndAddArena(name string) error {
	GlobalArenaList.Lock()
	defer GlobalArenaList.Unlock()

	_, exists := GlobalArenaList.name[name]
	if exists {
		return SameNameError{name}
	}

	newUUID := uuid.New()
	GlobalArenaList.name[name] = newUUID

	newArena := CreateArena()
	GlobalArenaList.arena[newUUID] = &newArena

	return nil
}
