package core

//go:generate go run ../generate_message.go -- Action

type Action interface {
	Data()
}

type Ron struct {
	TileToRon Tile      `json:"tile_to_ron"`
	WinResult WinResult `json:"win_result"` // TODO: Remove
}

type Tsumo struct {
	TileToTsumo Tile `json:"tile_to_tsumo"`
}

type Riichi struct {
	TileToRiichi Tile `json:"tile_to_riichi"`
}

type Toss struct {
	TileToToss Tile `json:"tile_to_toss"`
}

type Skip struct {
	ActionToSkip Action `json:"action_to_skip"` // wrap
}

type Pon struct {
	TileToPon Tile `json:"tile_to_pon"`
}

type Kan struct {
	TileToKan Tile `json:"tile_to_kan"`
}

type Chii struct {
	TileToChii  Tile    `json:"tile_to_chii"`
	TilesInHand [2]Tile `json:"tiles_in_hand"`
}

type Draw struct {
	DrawnTile Tile `json:"drawn_tile"`
}
