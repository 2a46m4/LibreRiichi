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
	StartGameAndRound() ([]MessageSendInfo, error)
	ContinueRound() ([]MessageSendInfo, error)
	HandleEvent(action Action, arenaIdx uint8) ([]MessageSendInfo, error)
	IsInGame() bool
	HasRoundEnded() bool
	// Send round results
	RoundEndCleanup() ([]MessageSendInfo, error)
	ShouldContinueRound() bool
	// Send game results
	GameEndCleanup() ([]MessageSendInfo, error)
}

// A location where players gather. Controls the flow of the game,
// directing messages to players, requesting input/ouput
type Arena struct {
	agents     []Client
	spectators []Client
	game       Game

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
		GameStarted: !arena.game.IsInGame(),
		DateCreated: arena.DateCreated,
	}
}

func (arena *Arena) Send(data ArenaEvent, sendTo uint8) {
	arena.agents[sendTo].GetRecv() <- ServerArenaEvent{
		ArenaMessage: data,
	}
}

func (arena *Arena) SendBoardEvent(data BoardEvent, sendTo uint8) {
	arena.Send(ArenaBoardEvent{BoardEvent: data}, sendTo)
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

	thisPlayer := len(arena.agents) - 1
	for i := range len(arena.agents) {
		if i == thisPlayer {
			continue
		}
		arena.Send(data, uint8(i))
	}
	agent.SetArena(arena)
	return nil
}

// StartArena is called by a client when a game should be started. It broadcasts a start round message to the connected players
func (arena *Arena) HandleStartGameActionData(data StartGameActionData, fromPlayer uint8) (UnitType, error) {
	arena.Lock()
	defer arena.Unlock()

	arena.log.Info("Handle start game called")

	if len(arena.agents) != 4 {
		return Unit, errors.New("Not enough agents")
	}

	sendInfo, err := arena.game.StartGameAndRound()
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

	arena.log.Info("Finish start game handle")

	return Unit, nil
}

func (arena *Arena) HandlePlayerActionData(data PlayerActionData, fromPlayer uint8) (UnitType, error) {
	arena.Lock()
	defer arena.Unlock()
	arena.log.Info("Driving game")
	sendInfos, err := arena.game.HandleEvent(data.Action, fromPlayer)

	if err != nil {
		arena.log.Info("Error: ", "Msg", err.Error())
		return Unit, err
	}

	if arena.game.HasRoundEnded() {
		infos, err := arena.game.RoundEndCleanup()
		if err != nil {
			return Unit, err
		}
		sendInfos = append(sendInfos, infos...)
	}

	if !arena.game.ShouldContinueRound() {
		// End the round by calling cleanup and sending
		// the data.
		infos, err := arena.game.GameEndCleanup()
		if err != nil {
			return Unit, err
		}
		sendInfos = append(sendInfos, infos...)
		for _, sendInfo := range sendInfos {
			for _, event := range sendInfo.Events {
				arena.SendBoardEvent(event, sendInfo.SendTo)
			}
		}	

		// TODO: Cleanup the arena itself
		return Unit, err

	} else { // Continue the round

		// New round data
		info, err := arena.game.ContinueRound()
		if err != nil {
			return Unit, err
		}
		sendInfos = append(sendInfos, info...)

		for _, sendInfo := range sendInfos {
			for _, event := range sendInfo.Events {
				arena.SendBoardEvent(event, sendInfo.SendTo)
			}
		}
	}

		

	return Unit, err
}

func (arena *Arena) HandlePlayerQuitActionData(data PlayerQuitActionData, fromPlayer uint8) (UnitType, error) {
	arena.Lock()
	defer arena.Unlock()
	if arena.game.IsInGame() {
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

// TODO: Add a finish game section, which cleans up resources and tells the clients that the arena will stop

// Gives a copy of the game
func (arena *Arena) QueryGame() Game {
	arena.Mutex.Lock()
	defer arena.Mutex.Unlock()
	return arena.game
}
