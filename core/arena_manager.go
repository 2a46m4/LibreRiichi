package core

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/google/uuid"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game"
	core "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
)

// TODO: This can potentially use RCU

type ArenaList struct {
	arena map[uuid.UUID]*Arena
	name  map[string]uuid.UUID
	logger *slog.Logger

	sync.RWMutex
}

var GlobalArenaList ArenaList = ArenaList{}

type EmptyNameError struct{}

type ArenaNotFoundError struct {
	arena_name string
}

type SameNameError struct {
	arena_name string
}

func InitializeMap() {
	GlobalArenaList.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	GlobalArenaList.logger.Info("Initializing map")
	GlobalArenaList.arena = make(map[uuid.UUID]*Arena)
	GlobalArenaList.name = make(map[string]uuid.UUID)
	// Debug initial room
	// TODO: Remove
	err := CreateAndAddArena("a")
	if err != nil {
		panic("Couldn't create debug arena")
	}
	
	arena, err := GetArenaFromName("a")
	if err != nil {
		panic("Couldn't find arena from name")
	}
	arena.HandleAddAIArenaAction(core.AddAIArenaAction{}, 0)
	arena.HandleAddAIArenaAction(core.AddAIArenaAction{}, 0)
	arena.HandleAddAIArenaAction(core.AddAIArenaAction{}, 0)
	
}

func (e EmptyNameError) Error() string {
	return "Arena name cannot be empty"
}

func (e ArenaNotFoundError) Error() string {
	return fmt.Sprintf("Arena %v is not found", e.arena_name)
}

func (e SameNameError) Error() string {
	return fmt.Sprintf("There is already an arena with the same name: %v", e.arena_name)
}

func ListArenas() []string {
	GlobalArenaList.RLock()
	defer GlobalArenaList.RUnlock()

	result := make([]string, 0, len(GlobalArenaList.name))
	for name := range GlobalArenaList.name {
		result = append(result, name)
	}
	GlobalArenaList.logger.Info("Listing arenas: ", "names", GlobalArenaList.name, "data", result)

	return result
}

func GetArenaFromName(name string) (*Arena, error) {
	log.Println("Getting arena from name")
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
		fmt.Println("Did not find arena: ", name)
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

	if len(name) == 0 {
		return EmptyNameError{}
	}

	_, exists := GlobalArenaList.name[name]
	if exists {
		return SameNameError{name}
	}

	newUUID := uuid.New()
	GlobalArenaList.name[name] = newUUID

	game := NewMahjongGame()

	newArena := CreateArena(name, newUUID, game)
	GlobalArenaList.arena[newUUID] = &newArena

	fmt.Println("Created arena. Now arena map is:", GlobalArenaList.name)
	return nil
}
