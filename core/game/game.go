package core

import (
	"errors"
	"fmt"

	core "codeberg.org/ijnakashiar/LibreRiichi/core"
)

type MahjongGame struct {
	Players []Player
	// Maps Player Index → Order
	PlayerToOrder []uint8
	// Maps Order → Player Index
	OrderToPlayer []uint8
	// Current turn lasts until everyone has finished their possible actions
	CurrentTurn        uint8
	CurrentTurnActions []ActionResult
	CurrentTurnPlayed  bool
	PostTurnActions    []ActionResult
	PostTurnPlayed     bool

	RoundWind Wind
	GameEnded bool

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
	game.CurrentTurn = 0
	game.CurrentTurnPlayed = false

	game.RoundWind = South
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

func (game *MahjongGame) currentPlayerIdx() uint8 {
	return game.OrderToPlayer[game.CurrentTurn]
}

func (game *MahjongGame) nextPlayerIdx() uint8 {
	return game.OrderToPlayer[game.CurrentTurn+1]
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
		setup[idx] = make([]Setup, 0, 7)
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
func (game *MahjongGame) RespondToAction(action PlayerAction) ([]ActionResult, bool) {
	switch action.Action {
	case CHII:
		return game.handleChii(action)
	case KAN:
		return game.handleKan(action)
	case PON:
	case RIICHI:
	case RON:
	case SKIP:

	case TOSS:
		return game.handleToss(action)
	case TSUMO:
	default:
		panic(fmt.Sprintf("unexpected core.ActionType: %#v", action.Action))
	}
	return nil, false
}

func (game *MahjongGame) handleChii(action PlayerAction) ([]ActionResult, bool) {
	chiiData := action.Data

	onTile, ok := chiiData["OnTile"].(int)
	if !ok {
		return nil, true
	}
	chiiSequence, ok := chiiData["ChiiSequence"].([]int)
	if !ok {
		return nil, true
	}
	if !game.CurrentTurnPlayed {
		return nil, true
	}
	if action.FromPlayer != game.nextPlayerIdx() {
		return nil, true
	}
	err := game.Players[action.FromPlayer].Chii(
		Tile(onTile),
		[2]Tile{
			Tile(chiiSequence[0]),
			Tile(chiiSequence[1]),
		})
	if err != nil {
		return nil, true
	}

	return []ActionResult{
		{action, GLOBAL},
	}, true
}

func (game *MahjongGame) handleKan(action PlayerAction) ([]ActionResult, bool) {
	return nil, false
}

func (game *MahjongGame) handleToss(action PlayerAction) ([]ActionResult, bool) {
	tossData := action.Data

	onTileInt, ok := tossData["OnTile"].(int)
	if !ok {
		return nil, true
	}
	onTile := Tile(onTileInt)
	if game.CurrentTurnPlayed {
		return nil, true
	}
	if action.FromPlayer != game.currentPlayerIdx() {
		return nil, true
	}
	err := game.Players[action.FromPlayer].Toss(onTile)
	if err != nil {
		return nil, true
	}

	return []ActionResult{
		{action, GLOBAL},
	}, true
}

// Checks the post-toss actions that can be made
func (game *MahjongGame) checkPostTossActions() []ActionResult {
	for idx := range 4 {

	}

	return nil
}

// Return the game results
func (MahjongGame) GetGameResults() (GameResult, error) {
	return GameResult{}, nil
}

// Returns the maximum amount of players
func (MahjongGame) GetMaxPlayers() int {
	return 4
}
