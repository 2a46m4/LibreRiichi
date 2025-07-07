package core

import (
	"errors"
	"fmt"
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
	uuid uuid.UUID

	sync.Mutex
}

type MessageSendInfo struct {
	Events     []ArenaBoardEventData
	Visibility Visibility
	SendTo     uint8
}

func (arena *Arena) GetArenaInfo() ArenaInfoResponseData {
	arena.Lock()
	defer arena.Unlock()

	agents := make([]AgentInfo, 0)
	for _, agent := range arena.agents {
		agents = append(agents, AgentInfo{Name: agent.Name})
	}

	return ArenaInfoResponseData{
		Success:     true,
		Name:        arena.Name,
		Agents:      agents,
		GameStarted: arena.gameStarted,
		DateCreated: arena.DateCreated,
	}
}

func (arena *Arena) Send(data ArenaMessage, visibility Visibility, sendTo uint8) error {
	switch visibility {
	case GLOBAL:
		for i, player := range arena.agents {
			fmt.Println("Sending index: ", i)
			player.Recv <- Message{
				MessageType: ServerArenaEventType,
				Data:        ServerArenaMessageEventData{ArenaMessage: data},
			}
		}

	case PARTIAL:
		arena.agents[sendTo].Recv <- Message{
			MessageType: ServerArenaEventType,
			Data:        ServerArenaMessageEventData{ArenaMessage: data},
		}

		altMessage, err := GetAltMessage(data)
		if err != nil {
			return err
		}

		for idx, player := range arena.agents {
			fmt.Println("Sending index: ", idx)
			if idx == int(sendTo) {
				continue
			}
			player.Recv <- Message{
				MessageType: ServerArenaEventType,
				Data:        ServerArenaMessageEventData{ArenaMessage: altMessage},
			}
		}

	case PLAYER:
		arena.agents[sendTo].Recv <- Message{
			MessageType: ServerArenaEventType,
			Data:        ServerArenaMessageEventData{ArenaMessage: data},
		}
	case EXCLUDE:
		for i, player := range arena.agents {
			fmt.Println("Exclude: sending index: ", i)
			if i == int(sendTo) {
				fmt.Println("Exclude: Skipping: ", i)
				continue
			}
			fmt.Println("Exclude: Continuing with: ", i)
			player.Recv <- Message{
				MessageType: ServerArenaEventType,
				Data:        ServerArenaMessageEventData{ArenaMessage: data},
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

	data := PlayerJoinedEventData{
		Name: agent.Name,
		ID:   agent.ID,
	}

	err := arena.Send(
		ArenaMessage{
			MessageType: PlayerJoinedEventType,
			Data:        data,
		}, EXCLUDE, uint8(len(arena.agents)-1))

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
		for event := range sendInfo.Events {
			arena.Send(ArenaMessage{
				MessageType: ArenaBoardEventType,
				Data:        event,
			}, sendInfo.Visibility, sendInfo.SendTo)
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

// TODO: Implement ServerArenaHandler
// StartArena is called when a game should be started. It broadcasts a start round message to the connected players
func (arena *Arena) HandleStartGameAction(data StartGameActionData, fromPlayer uint8) error {
	arena.Lock()
	defer arena.Unlock()

	if arena.gameStarted {
		return errors.New("Game already started")
	}

	if len(arena.agents) != 4 {
		return errors.New("Not enough agents")
	}

	setups, err := arena.game.StartNewGame()
	if err != nil {
		return err
	}

	// Send over the setups for each player
	for idx, setup := range setups {

		err = arena.Send(ArenaMessage{
			MessageType: ArenaBoardEventType,
			Data: ArenaBoardEventData{
				BoardEvent: BoardEvent{
					EventType: GameSetupEventType,
					Data: GameSetupEventData{
						Setup: setup,
					},
				},
			},
		}, PLAYER, uint8(idx))

		if err != nil {
			panic(err)
		}
	}

	arena.gameStarted = true
	return arena.driveGame()
}

func (arena *Arena) HandlePlayerAction(data PlayerActionData, fromPlayer uint8) error {
	arena.Lock()
	defer arena.Unlock()

	sendInfos, err := ActionDecode(&arena.game, data.ActionData, fromPlayer)
	if err != nil {
		return err
	}

	for _, sendInfo := range sendInfos {
		for _, event := range sendInfo.Events {
			arena.Send(ArenaMessage{
				MessageType: ArenaBoardEventType,
				Data:        event,
			}, sendInfo.Visibility, sendInfo.SendTo)
		}
	}

	err = arena.driveGame()
	if err != nil {
		panic("TODO: Error handling")
	}

	return nil
}

func (arena *Arena) HandlePlayerQuitAction(data PlayerQuitActionData, fromPlayer uint8) error {
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
		arena.Send(ArenaMessage{
			MessageType: PlayerQuitEventType,
			Data: PlayerQuitEventData{
				Name: agent.Name,
			},
		}, GLOBAL, 0)
	}
	return nil
}

// FinishRoundArena is called when the arena round should be finished. It broadcasts an end round message to the connected players
func (arena *Arena) FinishRoundArena() {
	arena.game.GetGameResults()
}

// EndArena is called when the arena is finished and all players should be disconnected
func (arena *Arena) EndArena() error {
	return nil
}
