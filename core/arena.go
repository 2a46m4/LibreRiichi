package core

import (
	"encoding/json"
	"log"
)

// A location where players gather. Controls the flow of the game,
// directing messages to players, requesting input/ouput
type Arena struct {
	Agents     []*Client
	Spectators []*Client
	Game       MahjongGame

	JoinChannel chan *Client
}

func (arena Arena) Send(data ArenaMessage, sendTo Visibility) error {
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

func (arena Arena) Loop() {

	for {
		// Check for new join requests
		select {
		case newRequest := <-arena.JoinChannel:
			err := arena.JoinArena(newRequest, true)
			if err != nil {
				// TODO: Send a error back to the request
				continue
			}

			data := PlayerJoinedEventData{
				Client: *newRequest,
			}
			bytes, err := json.Marshal(data)
			if err != nil {
				continue
			}

			err = arena.Send(
				ArenaMessage{
					MessageType: PlayerJoinedEventType,
					Data:        bytes,
				}, GLOBAL)
			if err != nil {
				panic(err)
			}

			newRequest.Arena = &arena
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
func (arena *Arena) JoinArena(agent *Client, joinAsPlayer bool) error {
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

	// Send over the setups for each player
	for idx, setup := range setups {
		err := arena.Send(ArenaMessage{
			MessageType: SetupEventType,
			Data:        SetupEventTypeData{setup[idx]},
			VisibleTo:   Visibility(idx),
		})
		if err != nil {
			panic(err)
		}
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
func (arena Arena) FinishRoundArena() {

}

// EndArena is called when the arena is finished and all players should be disconnected
func (arena Arena) EndArena() error {
	return nil
}
