package winresult

import (
	"codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
)

type WinResult struct {
	Yakus          []Yaku           `json:"yakus,omitempty"`
	WinningTile    tile.Tile        `json:"winning_tile,omitempty"`
	WonByRon       bool             `json:"won_by_ron,omitempty"`
	PointsTransfer []PointsTransfer `json:"points_transfer,omitempty"`
}

type PointsTransfer struct {
	To     uint8 `json:"to,omitempty"`
	From   uint8 `json:"from,omitempty"`
	Amount uint  `json:"amount,omitempty"`
}

type Yaku struct {
	Name string `json:"name,omitempty"`
	Han  int    `json:"han,omitempty"`
}
