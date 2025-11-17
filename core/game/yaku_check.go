package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type YakuChecker func(playerIdx uint8, data *MahjongRoundData, extraTile ...Tile) bool

var YakuCheckerMap = map[YakuType]YakuChecker {
	
}

func CheckYaku(playerIdx uint8, data *MahjongRoundData, extraTile ...Tile) map[YakuType]bool {
	return nil
}
