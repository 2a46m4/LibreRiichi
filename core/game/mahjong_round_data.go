package game

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
