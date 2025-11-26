package game

import (
	"slices"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/hand"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type TileData struct {
	LiveWall
	DeadWall    DeadWall
	Hands       [4]Hand
	DiscardPile [4]DiscardPile
}

func CreateNewRound() (game TileData) {
	game.newRound()
	return game
}

func (game *TileData) newRound() {
	// Get the complete set of tiles and shuffle them
	tiles := GetTileList()
	PermuteArray(tiles)

	for i := range 70 {
		game.LiveWall.tiles[i] = tiles[i]
	}
	game.LiveWall.index = 0

	deadWallTiles := tiles[70:84]
	game.DeadWall.Reset(deadWallTiles)
	game.DeadWall.dora.revealDora()

	for i := range 4 {
		game.Hands[i].ClosedHand.Add(tiles[84+i*13 : 84+(i+1)*13]...)
	}
}

func (game *TileData) IncrementRound() {
	game.newRound()
}

// Modifies the tile state and returns the tile drawn
func (game *TileData) Draw(playerIdx uint8) Draw {
	tile := game.GetLiveTile()
	game.Hands[playerIdx].Draw(tile)
	return Draw{DrawnTile: tile}
}

func (game *TileData) Discard(playerIdx uint8, tile Tile) Toss {
	game.Hands[playerIdx].Discard(tile)
	game.DiscardPile[playerIdx].Add(tile)
	return Toss{TileToToss: tile}
}

func (game *TileData) Pon(playerIdx uint8, againstPlayer uint8) Pon {
	tile := game.DiscardPile[againstPlayer].Remove()
	game.Hands[playerIdx].Pon(tile)
	return Pon{TileToPon: tile}
}

func (game *TileData) DaiminKan(playerIdx uint8, againstPlayer uint8) Kan {
	tile := game.DiscardPile[againstPlayer].Remove()
	game.Hands[playerIdx].DaiminKan(tile)
	return Kan{TileToKan: tile}
}

func (game *TileData) ShouminKan(playerIdx uint8, tile Tile) Kan {
	game.Hands[playerIdx].ShouminKan(tile)
	return Kan{TileToKan: tile}
}

func (game *TileData) AnKan(playerIdx uint8, tile Tile) Kan {
	game.Hands[playerIdx].AnKan(tile)
	return Kan{TileToKan: tile}
}

func (game *TileData) Chii(playerIdx uint8, tiles [2]Tile) Chii {
	tile := game.DiscardPile[(playerIdx+3)%4].Remove()
	allTiles := []Tile{tile, tiles[0], tiles[1]}
	slices.Sort(allTiles)
	game.Hands[playerIdx].Chii(allTiles[0], tiles)
	return Chii{TileToChii: tile, TilesInHand: tiles}
}

func (game *TileData) Riichi(playerIdx uint8, tile Tile) Riichi {
	game.Hands[playerIdx].Riichi(tile)
	return Riichi{TileToRiichi: tile}
}

func (game *TileData) Ron(playerIdx uint8, tile Tile) Ron {
	// TODO: Implement Ron logic and calculate WinResult
	return Ron{TileToRon: tile}
}

func (game *TileData) Tsumo(playerIdx uint8, tile Tile) Tsumo {
	// TODO: Implement Tsumo logic
	return Tsumo{TileToTsumo: tile}
}
