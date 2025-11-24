package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	winresult "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/win_result"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/yaku"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
)

type MahjongRoundData struct {
	scoring   Scoring
	tileState TileState
	// Player direction
	windState WindState
	turnState TurnState
	// Round direction
	roundWind   Wind
	roundNumber uint8
}

func InitMahjongRoundData() MahjongRoundData {
	return MahjongRoundData{
		scoring:   InitScoring(25000),
		tileState: CreateNewRound(),
		windState: WindState(0),
		turnState: InitTurnState(),
	}
}

func (data *MahjongRoundData) IncrementRound() {
	data.windState.IncrementWind()
	data.tileState = CreateNewRound()
	data.turnState = InitTurnState()
	data.roundNumber += 1
	// TODO: Implement switching round winds
}

func (data *MahjongRoundData) CheckNaki(playerIdx uint8) (info []MessageSendInfo) {
	for i := range uint8(4) {
		if playerIdx == i {
			continue
		}

		actions := PotentialActionEvent{}

		// Check for Kan
		kanResult := CheckKan(playerIdx, i, data.tileState)
		if kanResult != nil {
			actions.Actions = append(actions.Actions, kanResult)
		}

		// Check for Pon
		ponResult := CheckPon(playerIdx, i, data.tileState)
		if ponResult != nil {
			actions.Actions = append(actions.Actions, ponResult)
		}

		chiiResult := CheckChii(playerIdx, i, data.tileState)
		if chiiResult != nil {
			actions.Actions = append(actions.Actions, chiiResult)
		}

		ronResult := CheckRon(playerIdx, i, data)
		if ronResult != nil {
			actions.Actions = append(actions.Actions, ronResult)
		}

		info = append(info, MessageSendInfo{
			Events: []BoardEvent{actions},
			SendTo: i,
		})
	}

	return info
}

func CheckKan(discardedPlayerIdx, playerIdx uint8, tileState TileState) Action {
	discardedTile := tileState.DiscardPile[discardedPlayerIdx].Last()

	if tileState.Hands[playerIdx].ClosedHand.HasTileN(discardedTile) == 3 {
		return Kan{
			TileToKan: discardedTile,
		}
	}

	return nil
}

func CheckChii(discardedPlayerIdx, playerIdx uint8, tileState TileState) Action {
	if (discardedPlayerIdx+1)%4 != playerIdx {
		return nil
	}

	discardedTile := tileState.DiscardPile[discardedPlayerIdx].Last()
	playerHand := tileState.Hands[playerIdx].ClosedHand

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

func CheckPon(discardedPlayerIdx, playerIdx uint8, tileState TileState) Action {
	discardedTile := tileState.DiscardPile[discardedPlayerIdx].Last()
	playerHand := tileState.Hands[playerIdx].ClosedHand

	if playerHand.HasTileN(discardedTile) == 2 {
		return Pon{
			TileToPon: discardedTile,
		}
	}

	return nil
}

func CheckRon(discardedPlayerIdx, playerIdx uint8, data *MahjongRoundData) Action {
	discardedTile := data.tileState.DiscardPile[discardedPlayerIdx].Last()

	yakuContext := YakuContext{
		IsSelfDrawn:                false,
		HasCalledRiichi:            false,
		IsIppatsu:                  false,
		IsLastTileDrawnOrDiscarded: false,
		IsDeadWallCall:             false,
		IsFromOpponentKanCall:      false,
		IsDoubleRiichi:             false,
		IsTenhou:                   false,
		IsChiihou:                  false,
		IsHandOpen:                 false,
		HandInRiichi:               false,
	}

	// Check if adding the discarded tile would complete a winning hand
	hand := &data.tileState.Hands[playerIdx]
	canWin := CheckHandCanWin(hand, yakuContext, discardedTile)
	if canWin {
		list, points, err := CheckYakuAndScore(hand, yakuContext, discardedTile)
		if err != nil {
			panic("Shouldn't get here")
		}

		yaku := []winresult.Yaku{}
		for _, y := range list.Yakus {
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
					winresult.PointsTransfer{
						To:     playerIdx,
						From:   discardedPlayerIdx,
						Amount: points.Ron,
					},
				},
			},
		}
	}

	return nil
}

func CheckHandWaits(playerIdx uint8, data *MahjongRoundData) []Tile {
	if data.tileState.Hands[playerIdx].FullHand() {
		panic("Wrong use of check hand waits function")
	}

	return nil
}
