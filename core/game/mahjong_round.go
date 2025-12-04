package game

// Stores data of a round
type MahjongRound struct {
	roundState RoundState
	data       MahjongRoundData
}

func InitMahjongRound() MahjongRound {
	return MahjongRound{
		roundState: InitRoundState(),
		data:       InitMahjongRoundData(),
	}
}

func (round *MahjongRound) ContinueMahjongRound() {
	round.data.IncrementRound()
	round.roundState.Transition("start-round", round, false)
}


