package game

import (
	"errors"
	"slices"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/hand"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type TileState struct {
	LiveWall
	DeadWall    DeadWall
	Hands       [4]Hand
	DiscardPile [4]DiscardPile
}

func CreateNewRound() (game TileState) {
	game.newRound()
	return game
}

func (game *TileState) newRound() {
	// Get the complete set of tiles and shuffle them
	tiles := GetTileList()
	PermuteArray(tiles)

	for i := range 70 {
		game.LiveWall.tiles[i] = tiles[i]
	}
	game.LiveWall.index = 0

	deadWallTiles := tiles[70:84]
	game.DeadWall.Reset(deadWallTiles)
	game.DeadWall.revealDora()

	for i := range 4 {
		game.Hands[i].ClosedHand.Add(tiles[84+i*13:84+(i+1)*13]...)
	}
}

func (game *TileState) IncrementRound() {
	game.newRound()
}

func (game *TileState) Draw(playerIdx uint8) Draw {
	tile := game.GetLiveTile()
	game.Hands[playerIdx].Draw(tile)
	return Draw{DrawnTile: tile}
}

func (game *TileState) Discard(playerIdx uint8, tile Tile) Toss {
	game.Hands[playerIdx].Discard(tile)
	game.DiscardPile[playerIdx].Add(tile)
	return Toss{TileToToss: tile}
}

func (game *TileState) Pon(playerIdx uint8, againstPlayer uint8) Pon {
	tile := game.DiscardPile[againstPlayer].Remove()
	game.Hands[playerIdx].Pon(tile)
	return Pon{TileToPon: tile}
}

func (game *TileState) DaiminKan(playerIdx uint8, againstPlayer uint8) Kan {
	tile := game.DiscardPile[againstPlayer].Remove()
	game.Hands[playerIdx].DaiminKan(tile)
	return Kan{TileToKan: tile}
}

func (game *TileState) ShouminKan(playerIdx uint8, tile Tile) Kan {
	game.Hands[playerIdx].ShouminKan(tile)
	return Kan{TileToKan: tile}
}

func (game *TileState) AnKan(playerIdx uint8, tile Tile) Kan {
	game.Hands[playerIdx].AnKan(tile)
	return Kan{TileToKan: tile}
}

func (game *TileState) Chii(playerIdx uint8, tiles [2]Tile) Chii {
	tile := game.DiscardPile[(playerIdx+3)%4].Remove()
	allTiles := []Tile{tile, tiles[0], tiles[1]}
	slices.Sort(allTiles)
	game.Hands[playerIdx].Chii(allTiles[0], tiles)
	return Chii{TileToChii: tile, TilesInHand: tiles}
}

func (game *TileState) Riichi(playerIdx uint8, tile Tile) Riichi {
	game.Hands[playerIdx].Riichi(tile)
	return Riichi{TileToRiichi: tile}
}

func (game *TileState) Ron(playerIdx uint8, tile Tile) Ron {
	// TODO: Implement Ron logic and calculate WinResult
	return Ron{TileToRon: tile}
}

func (game *TileState) Tsumo(playerIdx uint8, tile Tile) Tsumo {
	// TODO: Implement Tsumo logic
	return Tsumo{TileToTsumo: tile}
}

func (game *TileState) Try(action Action, playerIdx uint8, againstPlayer uint8, isClosedKan ...bool) error {
	switch a := action.(type) {
	case Draw:
		if game.LiveWall.End() {
			return errors.New("No more tiles to draw")
		}
	case Toss:
		if !game.Hands[playerIdx].FullHand() {
			return errors.New("Can't discard when the hand is not full")
		}
	case Pon:
		game.Pon(playerIdx, againstPlayer)
	case Kan:
		if len(isClosedKan) > 0 && isClosedKan[0] {
			game.AnKan(playerIdx, a.TileToKan)
		} else if len(isClosedKan) > 0 && !isClosedKan[0] {
			game.ShouminKan(playerIdx, a.TileToKan)
		} else {
			game.DaiminKan(playerIdx, againstPlayer)
		}
	case Chii:
		game.Chii(playerIdx, a.TilesInHand)
	case Riichi:
		game.Riichi(playerIdx, a.TileToRiichi)
	case Tsumo:
		game.Tsumo(playerIdx, a.TileToTsumo)
	case Ron:
		game.Ron(playerIdx, a.TileToRon)
	case Skip:
		// Skip action doesn't modify the game state
	default:
		return nil
	}
}
