package core

import (
	"encoding/json"

	core "codeberg.org/ijnakashiar/LibreRiichi/core"
	game "codeberg.org/ijnakashiar/LibreRiichi/core/game"
	msg "codeberg.org/ijnakashiar/LibreRiichi/core/msg"
	"github.com/google/uuid"
)

type Player struct {
	Name       string        `json:"name"`
	Id         uuid.UUID     `json:"id"`
	Connection core.ConnChan `json:"-"`
}

// A location where players gather. Controls the flow of the game
type Arena struct {
	Players []Player
	Game    game.MahjongGame

	JoinChannel chan Player
}

func (arena Arena) Broadcast(data any) error {
	marshalledData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for _, player := range arena.Players {
		player.Connection.WriteChannel <- marshalledData
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
				// Send a error back to the request
				continue
			}

			err = arena.Broadcast(PlayerJoined(newRequest))
			if err != nil {
				panic(err)
			}

		default:
			break
		}

		// Check for new messages from players
		for _, player := range arena.Players {
			select {
			case msgReceived := <-player.Connection.DataChannel:
				if err, ok := msgReceived.(error); ok {
					// Handle problematic connection here
					panic(err)
				}

				Message := msg.Message{}
				err := json.Unmarshal(msgReceived.([]byte), Message)
				if err != nil {
					panic(err)
				}

				switch Message.MessageType {
				case msg.StartGameAction:
					err := arena.StartArena()
					if err != nil {
						continue
					}

					arena.GameLoop()
					arena.EndArena()

				case msg.QuitAction:
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

// Adds a player to the arena.
func (arena *Arena) JoinArena(player Player, joinAsPlayer bool) error {
	if !joinAsPlayer {
		panic("NYI")
	}

	err := arena.Game.JoinArena(len(arena.Players))
	if err != nil {
		return err
	}

	arena.Players = append(arena.Players, player)
	return nil
}

// StartArena is called when a game should be started. It broadcasts a start round message to the connected players
func (arena Arena) StartArena() error {
	setups, err := arena.Game.StartNewGame()
	if err != nil {
		return err
	}

	for setup := range setup {

	}

	return err
}

func (arena Arena) GameLoop() {

}

// FinishRoundArena is called when the arena round should be finished. It broadcasts an end round message to the connected players
func (arena Arena) FinishRoundArena() {

}

// EndArena is called when the arena is finished and all players should be disconnected
func (arena Arena) EndArena() error {
	return nil
}
