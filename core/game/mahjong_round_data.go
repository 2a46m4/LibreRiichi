package game

import (
	"fmt"
	"slices"

	meldfinder "codeberg.org/ijnakashiar/LibreRiichi/core/game/meld_finder"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	core "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/hand"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/yaku"
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
	canWin := CheckHandCanWin(hand, discardedTile, yakuContext)
	if canWin {
		// TODO: Calculate actual WinResult with yaku and scores
		// Iterate over each yaku, checking which we can get
		
		
		return Ron{
			TileToRon: discardedTile,
			WinResult: WinResult{}, // TODO
		}
	}

	return nil
}

// Checks if the hand currently has yaku if it is a full hand, or if it will with an extra tile
func GetHandYaku(
	playerIdx uint8,
	data *MahjongRoundData,
	yakuContext YakuContext,
	extraTile ...Tile,
) YakuType {
	if len(extraTile) == 0 {
		if !data.tileState.Hands[playerIdx].FullHand() {
			panic("Wrong use of get hand yaku function")
		}

		validYakus := CheckYaku(playerIdx, data, extraTile...)
		fmt.Println(validYakus)

	} else if len(extraTile) == 1 {
		data.tileState.Hands[playerIdx].ClosedHand.Add(extraTile...)

		validYakus := CheckYaku(playerIdx, data, extraTile...)
		fmt.Println(validYakus)

		data.tileState.Hands[playerIdx].ClosedHand.RemoveTile(extraTile...)

	} else {
		panic("Wrong use of get hand yaku function")
	}

	return NO_YAKU
}

func CheckHandWaits(playerIdx uint8, data *MahjongRoundData) []Tile {
	if data.tileState.Hands[playerIdx].FullHand() {
		panic("Wrong use of check hand waits function")
	}

	return nil
}

func CheckKokushiMusou(hand core.Hand) YakuType {
	if hand.Closed() {
		typeList := []Tile{Manzu, Pinzu, Souzu}
		numberList := []uint8{1, 9}
		tileList := []Tile{}
		for _, t := range typeList {
			for _, n := range numberList {
				tileList = append(tileList, MakeNumberTile(t, n))
			}
		}
		if hand.ClosedHand.HasTile(tileList...) && slices.Contains(tileList, extraTile) {
			return KOKUSHI_MUSOU_THIRTEEN_WAITS_YAKU, 0
		}
		hand.ClosedHand.Add(extraTile)
		totalValid := 0
		for _, tile := range tileList {
			count := hand.ClosedHand.HasTileN(tile)
			if count == 1 || count == 2 {
				totalValid += count
			}
		}
		if totalValid == 14 {
			return KOKUSHI_MUSOU_YAKU, 0
		}
		hand.ClosedHand.Pop(1)
	}
}
