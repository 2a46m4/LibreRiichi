package core

import (
	"errors"
	"log/slog"
	"os"
	"sync"
	"time"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
	"github.com/google/uuid"
)

type Client interface {
	GetName() string
	GetID() uuid.UUID
	SetArena(*Arena)
	GetRecv() chan<- any
	IsAI() bool
}

type Game interface {
	StartGame() ([]MessageSendInfo, error)
	StartRound() ([]MessageSendInfo, error)
	HandleEvent(action Action, arenaIdx uint8) ([]MessageSendInfo, error)
	RoundEnd() error
	GameEnd() error
	RoundEnded() bool
	GameEnded() bool
}

// A location where players gather. Controls the flow of the game,
// directing messages to players, requesting input/ouput
type Arena struct {
	agents      []Client
	spectators  []Client
	gameStarted bool
	game        Game

	DateCreated time.Time
	Name        string
	uuid        uuid.UUID

	log *slog.Logger

	sync.Mutex
}

type RoomFullError struct{}

func (RoomFullError) Error() string {
	return "Room is full"
}

func (arena *Arena) GetArenaInfo() ArenaInfoResponse {
	arena.Lock()
	defer arena.Unlock()

	agents := make([]AgentInfo, 0)
	for i, agent := range arena.agents {
		agents = append(agents, AgentInfo{
			Name:  agent.GetName(),
			Order: uint8(i),
		})
	}

	return ArenaInfoResponse{
		Success:     true,
		Name:        arena.Name,
		Agents:      agents,
		GameStarted: arena.gameStarted,
		DateCreated: arena.DateCreated,
	}
}

func (arena *Arena) Send(data ArenaEvent, sendTo uint8) {
	arena.agents[sendTo].GetRecv() <- ServerArenaEvent{
		ArenaMessage: data,
	}
}

func (arena *Arena) SendBoardEvent(data BoardEvent, sendTo uint8) {
	arena.agents[sendTo].GetRecv() <- ServerArenaEvent{
		ArenaMessage: ArenaBoardEvent{BoardEvent: data},
	}
}

func (arena *Arena) getPlayerIdx(client Client) (uint8, error) {
	arena.Lock()
	defer arena.Unlock()

	for i, ptr := range arena.agents {
		if ptr == client {
			return uint8(i), nil
		}
	}
	return 0, errors.New("not found")
}

func CreateArena(name string, uuid uuid.UUID, game Game) Arena {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	return Arena{
		agents:      make([]Client, 0),
		spectators:  make([]Client, 0),
		gameStarted: false,
		game:        game,
		DateCreated: time.Now(),
		Mutex:       sync.Mutex{},
		Name:        name,
		uuid:        uuid,
		log:         logger,
	}
}

func (arena *Arena) JoinArena(agent Client) error {
	arena.Lock()
	defer arena.Unlock()

	if len(arena.agents) >= 4 {
		return RoomFullError{}
	}

	arena.agents = append(arena.agents, agent)
	data := PlayerJoinedEvent{
		AgentInfo: AgentInfo{
			Name:  agent.GetName(),
			ID:    agent.GetID(),
			Order: uint8(len(arena.agents) - 1),
		},
	}

	thisPlayer := len(arena.agents)-1
	for i := range len(arena.agents) {
		if i == thisPlayer {
			continue
		}
		arena.Send(data, uint8(i))
	}
	agent.SetArena(arena)
	return nil
}

// Drives the game forward
func (arena *Arena) driveGame(action Action, fromPlayer uint8) error {
	arena.log.Info("Driving game")

	var sendInfos []MessageSendInfo
	var err error

	if arena.game.RoundEnded() {
		arena.FinishRoundArena()
		if arena.game.GameEnded() {
			arena.FinishGameArena()
			return nil
		}
		
		sendInfos, err = arena.game.StartRound()
	} else {
		sendInfos, err = arena.game.HandleEvent(action, fromPlayer)
	}
	
	if err != nil {
		arena.log.Info("Error: ", err)
		return err
	}

	for _, sendInfo := range sendInfos {
		for _, event := range sendInfo.Events {
			arena.SendBoardEvent(event, sendInfo.SendTo)
		}
	}

	return nil
}

// StartArena is called by a client when a game should be started. It broadcasts a start round message to the connected players
func (arena *Arena) HandleStartGameActionData(data StartGameActionData, fromPlayer uint8) (UnitType, error) {
	arena.Lock()
	defer arena.Unlock()

	arena.log.Info("Handle start game called")

	if arena.gameStarted {
		return Unit, errors.New("Game already started")
	}

	if len(arena.agents) != 4 {
		return Unit, errors.New("Not enough agents")
	}

	sendGameInfo, err := arena.game.StartGame()
	sendRoundInfo, err := arena.game.StartRound()
	sendInfo := append(sendGameInfo, sendRoundInfo...)
	if err != nil {
		return Unit, err
	}

	for i := range 4 {
		arena.Send(GameStartedEvent{}, uint8(i))
	}

	for _, info := range sendInfo {
		for _, event := range info.Events {
			arena.SendBoardEvent(event, info.SendTo)
		}
	}

	arena.gameStarted = true
	
	return Unit, nil
}

func (arena *Arena) HandlePlayerActionData(data PlayerActionData, fromPlayer uint8) (UnitType, error) {
	arena.Lock()
	defer arena.Unlock()

	err := arena.driveGame(data.Action, fromPlayer)
	if err != nil {
		panic("TODO: Error handling")
	}

	return Unit, err
}

func (arena *Arena) HandlePlayerQuitActionData(data PlayerQuitActionData, fromPlayer uint8) (UnitType, error) {
	arena.Lock()
	defer arena.Unlock()
	if arena.gameStarted {
		// TODO: Replace with AI
	} else if len(arena.agents) == 1 {
		// TODO: Cleanup
		arena.log.Info("Removing arena")
		RemoveArena(arena.Name)
	} else {
		agent := arena.agents[fromPlayer]
		Remove(&arena.agents, uint(fromPlayer))
		for i := range len(arena.agents) {
			arena.Send(PlayerQuitEvent{
				Name: agent.GetName(),
			}, uint8(i))
		}
	}
	return Unit, nil
}

func (arena *Arena) HandleAddAIArenaAction(data AddAIArenaAction, fromPlayer uint8) (UnitType, error) {
	aiClient, err := MakeComputerClient()
	if err != nil {
		return Unit, err
	}

	err = arena.JoinArena(&aiClient)
	if err != nil {
		return Unit, err
	}
	go aiClient.Loop()
	return Unit, nil
}

func (arena *Arena) HandleRemoveAIArenaAction(data RemoveAIArenaAction, fromPlayer uint8) (UnitType, error) {
	// TODO: Send Die struct to client
	arena.log.Info("NYI")
	return Unit, nil
}

func (arena *Arena) HandleGameInfoActionData(data GameInfoActionData, fromPlayer uint8) (UnitType, error) {
	arena.log.Info("NYI")
	return Unit, nil
}

// FinishRoundArena is called when the arena round should be finished. It broadcasts an end round message to the connected players
func (arena *Arena) FinishRoundArena() {
	panic("TODO")
	// arena.game.GetGameResults()
}

func (arena *Arena) FinishGameArena() {
	panic("TODO")
	// arena.game.GetGameResults()
}

// Gives a copy of the game
func (arena *Arena) QueryGame() Game {
	arena.Mutex.Lock()
	defer arena.Mutex.Unlock()
	return arena.game
}
