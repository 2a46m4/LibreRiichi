package core

import (
	"errors"

	core "codeberg.org/ijnakashiar/LibreRiichi/core"
)

type MahjongGame struct {
	Players []Player
	// Maps Player Index → Order
	PlayerToOrder []uint8
	// Maps Order → Player Index
	OrderToPlayer []uint8
	RoundWind     Wind
	GameEnded     bool

	Tiles        [136]Tile
	LiveWall     []Tile
	Dora         []Tile
	DoraRevealed int
	UraDora      []Tile
	KanDraw      []Tile
	KansDrawn    int
}

// Sets up the game and the tiles
func (game *MahjongGame) setupGame() {
	core.PermuteArray(game.PlayerToOrder)
	for idx, order := range game.PlayerToOrder {
		game.OrderToPlayer[order] = uint8(idx)
	}
	game.Tiles = [136]Tile(GetTileList())
	core.PermuteArray(game.Tiles[:])
	game.GameEnded = false

	tileItr := 0
	for idx, order := range game.PlayerToOrder {
		player := &game.Players[idx]
		player = &Player{
			Points:   25000,
			SeatWind: Wind(order) + East,
		}
		player.FreshHand(game.Tiles[tileItr : tileItr+13])
		tileItr += 13
	}

	game.Dora = game.Tiles[tileItr : tileItr+5]
	game.DoraRevealed = 1
	tileItr += 5
	game.UraDora = game.Tiles[tileItr : tileItr+5]
	tileItr += 5
	game.KanDraw = game.Tiles[tileItr : tileItr+4]
	game.KansDrawn = 0
	tileItr += 4
	game.LiveWall = game.Tiles[tileItr:]
}

func (game MahjongGame) JoinArena(PlayerIdx int) error {
	if PlayerIdx >= game.GetMaxPlayers() {
		return errors.New("No more space")
	}

	game.Players = append(game.Players, Player{})
	game.PlayerToOrder = append(game.PlayerToOrder, uint8(PlayerIdx))
	game.OrderToPlayer = append(game.OrderToPlayer, 0)

	return nil
}

// Returns data to send to clients when a new game can be started, otherwise an error
func (game MahjongGame) StartNewGame() ([][]Setup, error) {
	if len(game.Players) != 4 {
		return nil, errors.New("Not enough players")
	}

	game.setupGame()
	setup := make([][]Setup, 4)
	startingPoints := [4]uint32{
		game.Players[0].Points,
		game.Players[1].Points,
		game.Players[2].Points,
		game.Players[3].Points,
	}

	for idx, player := range game.Players {
		setup[idx] = make([]Setup, 7)
		setup[idx] = append(setup[idx],
			Setup{
				Type:     INITIAL_TILES,
				ToPlayer: uint8(idx),
				Data:     player.ClosedHand,
			},
			Setup{
				Type:     DORA,
				ToPlayer: uint8(idx),
				Data:     game.Dora[0],
			},
			Setup{
				Type:     PLAYER_NUMBER,
				ToPlayer: uint8(idx),
				Data:     uint8(idx),
			},
			Setup{
				Type:     PLAYER_ORDER,
				ToPlayer: uint8(idx),
				Data:     game.PlayerToOrder,
			},
			Setup{
				Type:     ROUND_NUMBER,
				ToPlayer: uint8(idx),
				Data:     0,
			},
			Setup{
				Type:     ROUND_WIND,
				ToPlayer: uint8(idx),
				Data:     game.RoundWind,
			},
			Setup{
				Type:     STARTING_POINTS,
				ToPlayer: uint8(idx),
				Data:     startingPoints,
			},
		)
	}

	return setup, nil
}

// Returns the updated game state and things to notify when a player action is taken
// Additionally returns whether the game should continue
func RespondToAction(action PlayerAction) ([]ActionResult, bool) {
	return nil, false
}

// Return the game results
func GetGameResults() (GameResult, error) {
	return GameResult{}, nil
}

// Returns the maximum amount of players
func (MahjongGame) GetMaxPlayers() int {
	return 4
}
