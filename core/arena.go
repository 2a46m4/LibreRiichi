package core

import "net"

// A location where players gather. Controls the flow of the game
type Arena struct {
	Players     []net.Conn
	Game        MahjongGame
	JoinChannel chan net.Conn
}

func (arena Arena) Loop() {
	for {
		newRequest := <-arena.JoinChannel
		err := arena.JoinArena(newRequest, true)
		if err != nil {
			// Send a error back to the request
			continue
		}

		if len(arena.Players) == arena.Game.GetMaxPlayers() {
			arena.StartArena()
			arena.GameLoop()
			arena.EndArena()
		}
	}
}

// Adds a player to the arena.
func (arena *Arena) JoinArena(player net.Conn, joinAsPlayer bool) error {
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
func (arena Arena) StartArena() {

}

func (arena Arena) GameLoop() {

}

// FinishRoundArena is called when the arena round should be finished. It broadcasts an end round message to the connected players
func (arena Arena) FinishRoundArena() {

}

// EndArena is called when the arena is finished and all players should be disconnected
func (arena Arena) EndArena() {

}
