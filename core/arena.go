package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

// A location where players gather. Controls the flow of the game,
// directing messages to players, requesting input/ouput
type Arena struct {
	Agents      []*Client
	Spectators  []*Client
	GameStarted bool
	Game        MahjongGame
	// AwaitingInputs []??? that stores the list of agents that it is waiting on

	sync.Mutex
}

type MessageSendInfo struct {
	Events     []ArenaBoardEventData
	Visibility Visibility
	SendTo     uint8
}

func (arena *Arena) Send(data ArenaMessage, visibility Visibility, sendTo uint8) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	switch visibility {
	case GLOBAL:
		for _, player := range arena.Agents {
			player.Recv <- Message{
				MessageType: ServerArenaEventType,
				Data:        bytes,
			}
		}
	case PARTIAL:
		arena.Agents[sendTo].Recv <- Message{
			MessageType: ServerArenaEventType,
			Data:        bytes,
		}

		altMessage, err := GetAltMessage(data)
		if err != nil {
			return err
		}

		altBytes, err := json.Marshal(altMessage)
		if err != nil {
			return err
		}

		for idx, player := range arena.Agents {
			if idx == int(sendTo) {
				continue
			}
			player.Recv <- Message{
				MessageType: ServerArenaEventType,
				Data:        altBytes,
			}
		}
	case PLAYER:
		arena.Agents[sendTo].Recv <- Message{
			MessageType: ServerArenaEventType,
			Data:        bytes,
		}
	default:
		panic(fmt.Sprintf("unexpected core.Visibility: %#v", sendTo))
	}

	return nil
}

func CreateArena() Arena {
	return Arena{
		Agents:      []*Client{},
		Spectators:  []*Client{},
		Game:        MahjongGame{},
		GameStarted: false,
		Mutex:       sync.Mutex{},
	}
}

func (arena *Arena) JoinArena(agent *Client, joinAsPlayer bool) error {
	if !joinAsPlayer {
		panic("NYI")
	}

	arena.Lock()
	defer arena.Unlock()

	arena.Agents = append(arena.Agents, agent)

	data := PlayerJoinedEventData{
		Client: *agent,
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = arena.Send(
		ArenaMessage{
			MessageType: PlayerJoinedEventType,
			Data:        bytes,
		}, GLOBAL, 0)

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

	sendInfos, gameContinue := arena.Game.GetNextEvent()

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

// TODO: Implement ServerArenaHandler
// StartArena is called when a game should be started. It broadcasts a start round message to the connected players
func (arena *Arena) HandleStartGameAction(StartGameActionData) error {
	arena.Lock()
	defer arena.Unlock()

	if arena.GameStarted {
		return errors.New("Game already started")
	}

	if len(arena.Agents) != 4 {
		return errors.New("Not enough agents")
	}

	setups, err := arena.Game.StartNewGame()
	if err != nil {
		return err
	}

	// Send over the setups for each player
	for idx, setup := range setups {

		err = arena.Send(ArenaMessage{
			MessageType: ArenaBoardEventType,
			Data: ArenaBoardEventData{
				BoardEvent{
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

	arena.GameStarted = true
	return arena.driveGame()
}

func (arena *Arena) HandlePlayerAction(data PlayerActionData) error {
	arena.Lock()
	defer arena.Unlock()

	sendInfos, valid := arena.Game.RespondToAction(data)
	if !valid {
		return errors.New("Invalid move")
	}

	for _, sendInfo := range sendInfos {
		for _, event := range sendInfo.Events {
			arena.Send(ArenaMessage{
				MessageType: ArenaBoardEventType,
				Data:        event,
			}, sendInfo.Visibility, sendInfo.SendTo)
		}
	}

	err := arena.driveGame()
	if err != nil {
		panic("TODO: Error handling")
	}

	return nil
}

func (arena *Arena) HandlePlayerQuitAction(data PlayerQuitActionData) error {
	panic("NYI")
}

// FinishRoundArena is called when the arena round should be finished. It broadcasts an end round message to the connected players
func (arena *Arena) FinishRoundArena() {
	arena.Game.GetGameResults()
}

// EndArena is called when the arena is finished and all players should be disconnected
func (arena *Arena) EndArena() error {
	return nil
}
