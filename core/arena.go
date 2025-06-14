package core

import (
	"encoding/json"
	"errors"
	"log"
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
	events []ArenaBoardEventData
	SendTo Visibility
}

func (arena *Arena) Send(data ArenaMessage, sendTo Visibility) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if sendTo == GLOBAL {
		for _, player := range arena.Agents {
			player.Recv <- Message{
				MessageType: ServerArenaEventType,
				Data:        bytes,
			}
		}
	} else {
		arena.Agents[sendTo].Recv <- Message{
			MessageType: ServerArenaEventType,
			Data:        bytes,
		}
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

	err := arena.Game.JoinArena(len(arena.Agents))
	if err != nil {
		return err
	}

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
		}, GLOBAL)

	if err != nil {
		panic(err)
	}

	agent.Arena = arena

	return nil
}

// StartArena is called when a game should be started. It broadcasts a start round message to the connected players
func (arena *Arena) StartArena() error {
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
		}, Visibility(idx))

		if err != nil {
			panic(err)
		}
	}

	arena.GameStarted = true
	return arena.driveGame()
}

func (arena *Arena) DriveGame() error {
	arena.Lock()
	defer arena.Unlock()
	return arena.driveGame()
}

// Drives the game forward
func (arena *Arena) driveGame() error {

	events, gameContinue := arena.Game.GetNextEvent()

	if !gameContinue {
		arena.FinishRoundArena()
		return nil
	}

	// Send the event to the players
	for _, event := range events {
		arena.Send(ArenaMessage{
			MessageType: ArenaBoardEventType,
			Data:        event.events,
		}, event.SendTo)
	}

	// TODO: Set timeout here

	input := <-inputChannel
	if err, ok := input.Data.(error); ok {
		// Handle problematic connection here
		panic(err)
	}

	var action PlayerAction
	err := action.DecodeAction(input.Data.([]byte))
	if err != nil {
		log.Println(err)
		goto Rewait
	}
	if action.FromPlayer != uint8(input.I) {
		goto Rewait
	}

	actionResults, validMove := arena.Game.RespondToAction(action)
	if !validMove {
		goto Rewait
	}

	// Send the results to the players
	for _, actionResult := range actionResults {
		arena.Send(ArenaMessage{
			MessageType: PlayerActionEventType,
			Data:        PlayerActionEventTypeData{actionResult},
			VisibleTo:   actionResult.VisibleTo,
		})
	}
	return nil
}

func (arena Arena) GameLoop() {

	// Collect data channels
	dataChannels := make([]chan any, len(arena.Agents))
	for _, agent := range arena.Agents {
		dataChannels = append(dataChannels, (agent.Connection.DataChannel))
	}
	inputChannel := FanIn(dataChannels)

	gameContinue := true
	for {

		var events []ActionResult
		events, gameContinue = arena.Game.GetNextEvent()

		// Send the event to the players
		for _, event := range events {
			arena.Send(ArenaMessage{
				MessageType: PlayerActionEventType,
				Data:        PlayerActionEventTypeData{event},
				VisibleTo:   event.VisibleTo,
			})
		}

		if !gameContinue {
			break
		}

		// Wait on the players to make a response
	Rewait:
		// TODO: Set timeout here

		input := <-inputChannel
		if err, ok := input.Data.(error); ok {
			// Handle problematic connection here
			panic(err)
		}

		var action PlayerAction
		err := action.DecodeAction(input.Data.([]byte))
		if err != nil {
			log.Println(err)
			goto Rewait
		}
		if action.FromPlayer != uint8(input.I) {
			goto Rewait
		}

		actionResults, validMove := arena.Game.RespondToAction(action)
		if !validMove {
			goto Rewait
		}

		// Send the results to the players
		for _, actionResult := range actionResults {
			arena.Send(ArenaMessage{
				MessageType: PlayerActionEventType,
				Data:        PlayerActionEventTypeData{actionResult},
				VisibleTo:   actionResult.VisibleTo,
			})
		}
	}

	arena.Game.GetGameResults()
	// Broadcast game end and results
}

// FinishRoundArena is called when the arena round should be finished. It broadcasts an end round message to the connected players
func (arena *Arena) FinishRoundArena() {

}

// EndArena is called when the arena is finished and all players should be disconnected
func (arena *Arena) EndArena() error {
	return nil
}
