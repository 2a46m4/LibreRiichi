package winresult

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/yaku"
)

type WinResult struct {
	Yakus       YakuType
	// WinningHand Hand
	WinningTile Tile
	WonByRon    bool
}
