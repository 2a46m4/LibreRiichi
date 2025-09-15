package core

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
	"github.com/google/uuid"
)

// A location where players gather. Controls the flow of the game,
// directing messages to players, requesting input/ouput
type Arena struct {
	agents      []*Client
	spectators  []*Client
	gameStarted bool
	game        MahjongGame
	// AwaitingInputs []??? that stores the list of agents that it is waiting on

	DateCreated time.Time
	Name        string
	uuid        uuid.UUID

	sync.Mutex
}

// Information needed to send a messsage
type MessageSendInfo struct {
	Events     []ArenaBoardEvent
	Visibility Visibility
	SendTo     uint8
}

type InfoList []MessageSendInfo

func (list *InfoList) Add(data ...MessageSendInfo) *InfoList {
	*list = append(*list, data...)
	return list
}

func (list *InfoList) AddGlobalMessage(data ...ArenaBoardEvent) *InfoList {
	global := GlobalMessage()
	global.Add(data...)
	*list = append(*list, global)
	return list
}

func (list *InfoList) AddPrivateMessage(sendTo uint8, data ...ArenaBoardEvent) *InfoList {
	private := PrivateMessage(sendTo)
	private.Add(data...)
	*list = append(*list, private)
	return list
}

func (list *InfoList) AddPartialMessage(sendTo uint8, data ...ArenaBoardEvent) *InfoList {
	partial := PartialMessage(sendTo)
	partial.Add(data...)
	*list = append(*list, partial)
	return list
}

func GlobalMessage() MessageSendInfo {
	return MessageSendInfo{
		Events:     []ArenaBoardEvent{},
		Visibility: GLOBAL,
		SendTo:     0,
	}
}

func PrivateMessage(sendTo uint8) MessageSendInfo {
	return MessageSendInfo{
		Events:     []ArenaBoardEvent{},
		Visibility: PLAYER,
		SendTo:     sendTo,
	}
}

func PartialMessage(sendTo uint8) MessageSendInfo {
	return MessageSendInfo{
		Events:     []ArenaBoardEvent{},
		Visibility: PARTIAL,
		SendTo:     sendTo,
	}
}

// Add an event to the send info
func (info *MessageSendInfo) Add(data ...ArenaBoardEvent) *MessageSendInfo {
	info.Events = append(info.Events, data...)
	return info
}

func (arena *Arena) GetArenaInfo() ArenaInfoResponse {
	arena.Lock()
	defer arena.Unlock()

	agents := make([]AgentInfo, 0)
	for _, agent := range arena.agents {
		agents = append(agents, AgentInfo{Name: agent.Name})
	}

	return ArenaInfoResponse{
		Success:     true,
		Name:        arena.Name,
		Agents:      agents,
		GameStarted: arena.gameStarted,
		DateCreated: arena.DateCreated,
	}
}

// TODO: This should be explicit
func (arena *Arena) Send(data ArenaEvent, visibility Visibility, sendTo uint8) error {
	switch visibility {
	case GLOBAL:
		for i, player := range arena.agents {
			fmt.Println("Sending index: ", i)
			player.Recv <- ServerArenaEvent{
				ArenaMessage: data,
			}
		}

	case PARTIAL:
		arena.agents[sendTo].Recv <- ServerArenaEvent{ArenaMessage: data}

		altMessage, err := GetAltMessage(data)
		if err != nil {
			return err
		}

		for idx, player := range arena.agents {
			fmt.Println("Sending index: ", idx)
			if idx == int(sendTo) {
				continue
			}
			player.Recv <- ServerArenaEvent{
				ArenaMessage: altMessage,
			}
		}

	case PLAYER:
		arena.agents[sendTo].Recv <- ServerArenaEvent{
			ArenaMessage: data,
		}
	case EXCLUDE:
		for i, player := range arena.agents {
			fmt.Println("Exclude: sending index: ", i)
			if i == int(sendTo) {
				fmt.Println("Exclude: Skipping: ", i)
				continue
			}
			fmt.Println("Exclude: Continuing with: ", i)
			player.Recv <- ServerArenaEvent{
				ArenaMessage: data,
			}
		}
	default:
		panic(fmt.Sprintf("unexpected core.Visibility: %#v", sendTo))
	}

	return nil
}

func CreateArena(name string, uuid uuid.UUID) Arena {
	return Arena{
		agents:      make([]*Client, 0),
		spectators:  make([]*Client, 0),
		gameStarted: false,
		game:        MahjongGame{},
		DateCreated: time.Now(),
		Mutex:       sync.Mutex{},
		Name:        name,
		uuid:        uuid,
	}
}

func (arena *Arena) JoinArena(agent *Client, joinAsPlayer bool) error {
	if !joinAsPlayer {
		panic("NYI")
	}

	arena.Lock()
	defer arena.Unlock()

	arena.agents = append(arena.agents, agent)

	data := PlayerJoinedEvent{
		Name: agent.Name,
		ID:   agent.ID,
	}

	err := arena.Send(
		data, EXCLUDE, uint8(len(arena.agents)-1))

	if err != nil {
		panic(err)
	}

	agent.Arena = arena

	return nil
}

func (arena *Arena) DriveGame() error {
	arena.Lock()
	defer arena.Unlock()
	return arena.driveGame()
}

// Drives the game forward
func (arena *Arena) driveGame() error {

	sendInfos, gameContinue := arena.game.GetNextEvent()

	if !gameContinue {
		arena.FinishRoundArena()
		return nil
	}

	// Send the event to the players
	for _, sendInfo := range sendInfos {
		for _, event := range sendInfo.Events {
			arena.Send(event, sendInfo.Visibility, sendInfo.SendTo)
		}
	}

	return nil
}

func (arena *Arena) getPlayerIdx(client *Client) (uint8, error) {
	arena.Lock()
	defer arena.Unlock()

	for i, ptr := range arena.agents {
		if ptr == client {
			return uint8(i), nil
		}
	}
	return 0, errors.New("not found")
}

// StartArena is called when a game should be started. It broadcasts a start round message to the connected players
func (arena *Arena) HandleStartGameActionData(data StartGameActionData, fromPlayer uint8) (UnitType, error) {
	arena.Lock()
	defer arena.Unlock()

	log.Println("Handle start game called")

	if arena.gameStarted {
		return Unit, errors.New("Game already started")
	}

	if len(arena.agents) != 4 {
		return Unit, errors.New("Not enough agents")
	}

	setups, err := arena.game.StartNewGame()
	if err != nil {
		return Unit, err
	}

	// Send over the setups for each player
	for idx, setup := range setups {

		err = arena.Send(ArenaBoardEvent{
			BoardEvent: GameSetupEvent{
				Setup: setup,
			}}, PLAYER, uint8(idx))

		if err != nil {
			panic(err)
		}
	}

	arena.gameStarted = true
	return Unit, arena.driveGame()
}

func (arena *Arena) HandlePlayerActionData(data PlayerActionData, fromPlayer uint8) (UnitType, error) {
	arena.Lock()
	defer arena.Unlock()

	sendInfos, err := ActionDecode(&arena.game, data.Action, fromPlayer)
	if err != nil {
		return Unit, err
	}

	for _, sendInfo := range sendInfos {
		for _, event := range sendInfo.Events {
			arena.Send(event, sendInfo.Visibility, sendInfo.SendTo)
		}
	}

	err = arena.driveGame()
	if err != nil {
		panic("TODO: Error handling")
	}

	return Unit, err
}

func (arena *Arena) HandlePlayerQuitActionData(data PlayerQuitActionData, fromPlayer uint8) (UnitType, error) {
	arena.Lock()
	defer arena.Unlock()
	if arena.gameStarted {
		// Replace with AI

	} else if len(arena.agents) == 1 {
		fmt.Println("Removing arena")
		RemoveArena(arena.Name)
	} else {
		agent := arena.agents[fromPlayer]
		Remove(&arena.agents, uint(fromPlayer))
		arena.Send(PlayerQuitEvent{
			Name: agent.Name,
		}, GLOBAL, 0)
	}
	return Unit, nil
}

// FinishRoundArena is called when the arena round should be finished. It broadcasts an end round message to the connected players
func (arena *Arena) FinishRoundArena() {
	arena.game.GetGameResults()
}

// EndArena is called when the arena is finished and all players should be disconnected
func (arena *Arena) EndArena() error {
	return nil
}
