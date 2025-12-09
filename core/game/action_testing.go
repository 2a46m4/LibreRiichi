package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game/meld_finder"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/hand"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
	winresult "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/win_result"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/yaku"
)

func CanAnkan(kan Kan, hand *Hand) bool {
	return !hand.InRiichi && hand.FullHand() && (hand.ClosedHand.HasTileN(kan.TileToKan) == 4)
}

func CanRon(
	discardedPlayerIdx uint8,
	discardedTile Tile,
	playerIdx uint8,
	playerHand *Hand,
	yakuContext YakuContext) *Ron {
	// discardedTile := data.tileData.DiscardPile[discardedPlayerIdx].Last()

	// yakuContext := YakuContext{
	// 	IsSelfDrawn:           false,
	// 	IsIppatsu:             false,
	// 	IsLastLiveTile:        false,
	// 	IsDeadWallCall:        false,
	// 	IsFromOpponentKanCall: false,
	// 	IsDoubleRiichi:        false,
	// 	IsTenhou:              false,
	// 	IsChiihou:             false,
	// 	HandInRiichi:          false,
	// }

	// Check if adding the discarded tile would complete a winning hand
	// hand := &data.tileData.Hands[playerIdx]
	yakuList, pointValue, err := CheckYakuAndScore(playerHand, yakuContext)
	if err == nil {

		yaku := []winresult.Yaku{}
		for _, y := range yakuList.Yakus {
			yaku = append(yaku, winresult.Yaku{
				Name: y.YakuName,
				Han:  y.HanValue,
			})
		}

		return &Ron{
			TileToRon: discardedTile,
			WinResult: winresult.WinResult{
				Yakus:       yaku,
				WinningTile: discardedTile,
				WonByRon:    true,
				PointsTransfer: []winresult.PointsTransfer{
					{
						To:     playerIdx,
						From:   discardedPlayerIdx,
						Amount: pointValue.Ron,
					},
				},
			},
		}
	}

	return nil
}

func GetChiiTargets(
	discardedPlayerIdx uint8, discardedTile Tile,
	playerIdx uint8, playerHand *Hand,
) (list []Chii) {
	if (discardedPlayerIdx+1)%4 != playerIdx {
		return nil
	}

	// discardedTile := tileData.DiscardPile[discardedPlayerIdx].Last()
	// playerHand := tileData.Hands[playerIdx].ClosedHand

	// Chii can only be done by the player immediately after the discarder (playerIdx == (discardedPlayerIdx + 1) % 4)
	// For now, we'll check all potential Chii combinations and return the first valid one

	// Check for sequences where discarded tile is the first tile
	if !discardedTile.IsHonour() {
		tileNum := discardedTile.GetTileNumber()
		tileType := discardedTile & TileMask

		// Check if player has the next two tiles in sequence
		if tileNum <= 6 { // Need room for two tiles after
			nextTile1 := tileType | Tile(tileNum+1)
			nextTile2 := tileType | Tile(tileNum+2)
			if playerHand.ClosedHand.HasTile(nextTile1) && playerHand.ClosedHand.HasTile(nextTile2) {
				list = append(list, Chii{
					TileToChii:  discardedTile,
					TilesInHand: [2]Tile{nextTile1, nextTile2},
				})
			}
		}

		// Check if player has tile before and after (discarded tile is middle)
		if tileNum >= 1 && tileNum <= 7 {
			prevTile := tileType | Tile(tileNum-1)
			nextTile := tileType | Tile(tileNum+1)
			if playerHand.ClosedHand.HasTile(prevTile) && playerHand.ClosedHand.HasTile(nextTile) {
				list = append(list, Chii{
					TileToChii:  discardedTile,
					TilesInHand: [2]Tile{prevTile, nextTile},
				})
			}
		}

		// Check if player has the two previous tiles (discarded tile is last)
		if tileNum >= 2 { // Need room for two tiles before
			prevTile1 := tileType | Tile(tileNum-2)
			prevTile2 := tileType | Tile(tileNum-1)
			if playerHand.ClosedHand.HasTile(prevTile1) && playerHand.ClosedHand.HasTile(prevTile2) {
				list = append(list, Chii{
					TileToChii:  discardedTile,
					TilesInHand: [2]Tile{prevTile1, prevTile2},
				})
			}
		}
	}

	return list
}

func CheckPon(playerHand *Hand, discardedTile Tile) bool {
	return playerHand.ClosedHand.HasTileN(discardedTile) >= 2
}

// Ankan only
func CheckKan(hand *Hand, discardedTile Tile) bool {
	return hand.ClosedHand.HasTileN(discardedTile) == 3
}

// Requires a full hand, representing a hand that just drew a tile
func GetRiichiTargets(hand *Hand) (list []Riichi) {
	if !hand.Closed() {
		return nil
	}

	if hand.InRiichi {
		return nil
	}

	closed := hand.ClosedHand
	for _, tile := range GetTileList() {
		handCopy := append([]Tile{}, closed.GetHand()...)
		for _, curTile := range handCopy {
			// Impossible meld
			if curTile != tile && closed.HasTileN(tile) == 4 {
				continue
			}

			closed.RemoveTile(curTile)
			closed.Add(tile)

			if HasWinningCombination(closed.GetHand(), 0) {
				list = append(list, Riichi{
					TileToRiichi: tile,
				})
			}

			closed.Pop(1)
			closed.Add(curTile)
		}
	}

	return list
}
