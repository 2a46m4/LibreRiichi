package core

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

// Agent represents a player that has joined the arena.
type Agent struct {
	Name       string    `json:"name"`
	Id         uuid.UUID `json:"id"`
	Connection ConnChan  `json:"-"`
}

// A location where players gather. Controls the flow of the game
type Arena struct {
	Agents []Agent
	Game   MahjongGame

	JoinChannel chan Agent
}

func (arena Arena) Broadcast(data any) error {
	marshalledData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for _, player := range arena.Agents {
		player.Connection.WriteChannel <- marshalledData
	}

	return nil
}

func (arena Arena) Send(data ArenaMessage) error {
	return nil
}

func (arena Arena) Loop() {

	for {
		// Check for new join requests
		select {
		case newRequest := <-arena.JoinChannel:
			err := arena.JoinArena(newRequest, true)
			if err != nil {
				// Send a error back to the request
				continue
			}

			err = arena.Broadcast(
				ArenaMessage{
					MessageType: PlayerJoinedEventType,
					Data:        PlayerJoinedEventData{},
					VisibleTo:   GLOBAL,
				})
			if err != nil {
				panic(err)
			}

		default:
			break
		}

		// Check for new messages from players
		for _, player := range arena.Agents {
			select {
			case msgReceived := <-player.Connection.DataChannel:
				if err, ok := msgReceived.(error); ok {
					// Handle problematic connection here
					panic(err)
				}

				Message := ArenaMessage{}
				err := json.Unmarshal(msgReceived.([]byte), &Message)
				if err != nil {
					panic(err)
				}

				switch Message.MessageType {
				case StartGameActionType:
					err := arena.StartArena()
					if err != nil {
						continue
					}

					arena.GameLoop()
					arena.EndArena()

				case QuitActionType:
					err := arena.EndArena()
					if err != nil {
						continue
					}
				}

			default:
				continue
			}
		}
	}
}

// Adds an agnet to the arena.
func (arena *Arena) JoinArena(agent Agent, joinAsPlayer bool) error {
	if !joinAsPlayer {
		panic("NYI")
	}

	err := arena.Game.JoinArena(len(arena.Agents))
	if err != nil {
		return err
	}

	arena.Agents = append(arena.Agents, agent)
	return nil
}

// StartArena is called when a game should be started. It broadcasts a start round message to the connected players
func (arena Arena) StartArena() error {
	setups, err := arena.Game.StartNewGame()
	if err != nil {
		return err
	}

	// Send over the data
	for idx, setup := range setups {
		data, err := json.Marshal(ArenaMessage{
			MessageType: 0,
			Data:        nil,
		})
		if err != nil {
			panic(err)
		}
		arena.Agents[idx].Connection.WriteChannel <- data
	}

	return nil
}

func (arena Arena) GameLoop() {
	dataChannels := make([]chan any, len(arena.Agents))
	for _, agent := range arena.Agents {
		dataChannels = append(dataChannels, (agent.Connection.DataChannel))
	}
	inputChannel := FanIn(dataChannels)
	gameContinue := true

	for gameContinue {

		input := <-inputChannel
		if err, ok := input.Data.(error); ok {
			// Handle problematic connection here
			panic(err)
		}

		var action PlayerAction
		err := action.DecodeAction(input.Data.([]byte))
		if err != nil {
			log.Println(err)
			continue
		}
		if action.FromPlayer != uint8(input.I) {
			continue
		}

		var actionResults []ActionResult
		actionResults, gameContinue = arena.Game.RespondToAction(action)

		// Send the results to the players
		for _, actionResult := range actionResults {
			if actionResult.VisibleTo == GLOBAL {
				for _, agent := range arena.Agents {
					data, err := json.Marshal(actionResult)
					if err != nil {
						panic(err)
					}

					agent.Connection.WriteChannel <- data
				}
			}
		}
	}

	arena.Game.GetGameResults()
}

// FinishRoundArena is called when the arena round should be finished. It broadcasts an end round message to the connected players
func (arena Arena) FinishRoundArena() {

}

// EndArena is called when the arena is finished and all players should be disconnected
func (arena Arena) EndArena() error {
	return nil
}
