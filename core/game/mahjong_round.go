package game

type MahjongRound struct {
	roundState RoundState
	data       MahjongRoundData
}

func InitMahjongRound(roundNumber uint8) MahjongRound {
	return MahjongRound{
		roundState: *InitRoundState(),
		data:       InitMahjongRoundData(roundNumber),
	}
}
