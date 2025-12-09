package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	winresult "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/win_result"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/yaku"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
)

type MahjongRoundData struct {
	scoring  Scoring
	tileData TileData
	// Player direction
	turnData TurnData
	// Round direction
	roundWind   Wind
	roundNumber uint8
}

func (data *MahjongRoundData) CheckNaki(justDiscarded uint8, againstPlayed uint8) MessageSendInfo {
	if justDiscarded == againstPlayed {
		return MessageSendInfo{}
	}

	actions := PotentialActionEvent{}

	kanResult := CheckKan(justDiscarded, againstPlayed, data.tileData)
	if kanResult != nil {
		actions.Actions = append(actions.Actions, kanResult)
	}

	ponResult := CheckPon(justDiscarded, againstPlayed, data.tileData)
	if ponResult != nil {
		actions.Actions = append(actions.Actions, ponResult)
	}

	chiiResult := CheckChii(justDiscarded, againstPlayed, data.tileData)
	if chiiResult != nil {
		actions.Actions = append(actions.Actions, chiiResult)
	}

	ronResult := CheckRon(justDiscarded, againstPlayed, data)
	if ronResult != nil {
		actions.Actions = append(actions.Actions, ronResult)
	}

	return MessageSendInfo{
		Events: []BoardEvent{actions},
		SendTo: againstPlayed,
	}
}

func CheckKan(discardedPlayerIdx, playerIdx uint8, tileData TileData) Action {
	discardedTile := tileData.DiscardPile[discardedPlayerIdx].Last()

	if tileData.Hands[playerIdx].ClosedHand.HasTileN(discardedTile) == 3 {
		return Kan{
			TileToKan: discardedTile,
		}
	}

	return nil
}

func CheckChii(discardedPlayerIdx, playerIdx uint8, tileData TileData) Action {
	if (discardedPlayerIdx+1)%4 != playerIdx {
		return nil
	}

	discardedTile := tileData.DiscardPile[discardedPlayerIdx].Last()
	playerHand := tileData.Hands[playerIdx].ClosedHand

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
			if playerHand.HasTile(nextTile1) && playerHand.HasTile(nextTile2) {
				return Chii{
					TileToChii:  discardedTile,
					TilesInHand: [2]Tile{nextTile1, nextTile2},
				}
			}
		}

		// Check if player has tile before and after (discarded tile is middle)
		if tileNum >= 1 && tileNum <= 7 {
			prevTile := tileType | Tile(tileNum-1)
			nextTile := tileType | Tile(tileNum+1)
			if playerHand.HasTile(prevTile) && playerHand.HasTile(nextTile) {
				return Chii{
					TileToChii:  discardedTile,
					TilesInHand: [2]Tile{prevTile, nextTile},
				}
			}
		}

		// Check if player has the two previous tiles (discarded tile is last)
		if tileNum >= 2 { // Need room for two tiles before
			prevTile1 := tileType | Tile(tileNum-2)
			prevTile2 := tileType | Tile(tileNum-1)
			if playerHand.HasTile(prevTile1) && playerHand.HasTile(prevTile2) {
				return Chii{
					TileToChii:  discardedTile,
					TilesInHand: [2]Tile{prevTile1, prevTile2},
				}
			}
		}
	}

	return nil
}

func CheckPon(discardedPlayerIdx, playerIdx uint8, tileData TileData) Action {
	discardedTile := tileData.DiscardPile[discardedPlayerIdx].Last()
	playerHand := tileData.Hands[playerIdx].ClosedHand

	if playerHand.HasTileN(discardedTile) == 2 {
		return Pon{
			TileToPon: discardedTile,
		}
	}

	return nil
}

func CheckRon(discardedPlayerIdx, playerIdx uint8, data *MahjongRoundData) Action {
	discardedTile := data.tileData.DiscardPile[discardedPlayerIdx].Last()

	yakuContext := YakuContext{
		IsSelfDrawn:           false,
		IsIppatsu:             false,
		IsLastLiveTile:        false,
		IsDeadWallCall:        false,
		IsFromOpponentKanCall: false,
		IsDoubleRiichi:        false,
		IsTenhou:              false,
		IsChiihou:             false,
		HandInRiichi:          false,
	}

	// Check if adding the discarded tile would complete a winning hand
	hand := &data.tileData.Hands[playerIdx]
	yakuList, pointValue, err := CheckYakuAndScore(hand, yakuContext)
	if err == nil {
		if err != nil {
			panic("Shouldn't get here")
		}

		yaku := []winresult.Yaku{}
		for _, y := range yakuList.Yakus {
			yaku = append(yaku, winresult.Yaku{
				Name: y.YakuName,
				Han:  y.HanValue,
			})
		}

		return Ron{
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

func CheckHandWaits(playerIdx uint8, data *MahjongRoundData) []Tile {
	if data.tileData.Hands[playerIdx].FullHand() {
		panic("Wrong use of check hand waits function")
	}

	return nil
}
