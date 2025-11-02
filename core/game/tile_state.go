package game

import (
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

	// Fill the live wall (70 tiles)
	for i := range 70 {
		game.LiveWall.tiles[i] = tiles[i]
	}
	game.LiveWall.index = 0

	// Fill the dead wall (14 tiles total) using Reset method
	deadWallTiles := tiles[70:84]
	game.DeadWall.Reset(deadWallTiles)
}

func (game *TileState) IncrementRound() {
	game.newRound()
}

func (game *TileState) Draw(playerIdx uint8) {
	game.Hands[playerIdx].Draw(game.GetLiveTile())
}

func (game *TileState) Discard(playerIdx uint8, tile Tile) {
	game.Hands[playerIdx].Discard(tile)
	game.DiscardPile[playerIdx].Add(tile)
}

func (game *TileState) Pon(playerIdx uint8, againstPlayer uint8) {
	tile := game.DiscardPile[againstPlayer].Remove()
	game.Hands[playerIdx].Pon(tile)
}

func (game *TileState) DaiminKan(playerIdx uint8, againstPlayer uint8) {
	tile := game.DiscardPile[againstPlayer].Remove()
	game.Hands[playerIdx].DaiminKan(tile)
}

func (game *TileState) ShouminKan(playerIdx uint8, tile Tile) {
	game.Hands[playerIdx].ShouminKan(tile)
}

func (game *TileState) AnKan(playerIdx uint8, tile Tile) {
	game.Hands[playerIdx].AnKan(tile)
}

func (game *TileState) Chii(playerIdx uint8, tiles [2]Tile) {
	tile := game.DiscardPile[(playerIdx+3)%4].Remove()
	allTiles := []Tile{tile, tiles[0], tiles[1]}
	slices.Sort(allTiles)
	game.Hands[playerIdx].Chii(allTiles[0], tiles)
}

func (game *TileState) Riichi(playerIdx uint8, tile Tile) {
	game.Hands[playerIdx].Riichi(tile)
}

func (game *TileState) Ron(playerIdx uint8, tile Tile) {

}

func (game *TileState) Tsumo(playerIdx uint8, tile Tile) {

}
