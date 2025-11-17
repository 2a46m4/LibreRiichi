package game

import (
	"fmt"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type MahjongRoundData struct {
	tileState TileState
	windState WindState
	turnState TurnState
}

func InitMahjongRoundData() MahjongRoundData {
	return MahjongRoundData{
		tileState: CreateNewRound(),
		windState: WindState(0),
		turnState: InitTurnState(),
	}
}

func (data *MahjongRoundData) IncrementRound() {
	data.windState.IncrementWind()
	data.tileState = CreateNewRound()
	data.turnState = InitTurnState()
}

// A quick check to see if the hand can win
func CheckHandCanWin(playerIdx uint8, data *MahjongRoundData, extraTile ...Tile) bool {
	return false
}

// Checks if the hand currently has yaku if it is a full hand, or if it will with an extra tile
func GetHandYaku(playerIdx uint8, data *MahjongRoundData, extraTile ...Tile) YakuType {
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
	
	// 
}

