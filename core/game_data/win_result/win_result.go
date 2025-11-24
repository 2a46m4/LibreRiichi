package winresult

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/yaku"
)

type WinResult struct {
	Yakus          YakuList
	WinningTile    Tile
	WonByRon       bool
	PointsTransfer []PointsTransfer
}

type PointsTransfer struct {
	To     uint8
	From   uint8
	Amount uint
}
