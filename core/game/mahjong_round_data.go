package game

type MahjongRoundData struct {
	tileState TileState
	windState WindState
	turnState TurnState
}

func InitMahjongRoundData(roundIndex uint8) MahjongRoundData {
	return MahjongRoundData{
		tileState: CreateNewRound(),
		windState: WindState(roundIndex % 4),
		turnState: InitTurnState(),
	}
}

func (data *MahjongRoundData) IncrementRound() {
	data.windState.IncrementWind()
	data.tileState = CreateNewRound()
	data.turnState = InitTurnState()
}
