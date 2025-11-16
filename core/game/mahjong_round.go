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

func (round *MahjongRound) GetReturn() (ret any, ok bool) {
	ret, ok = round.roundState.RoundFSM.Metadata("return")
	round.roundState.RoundFSM.DeleteMetadata("return")
	return ret, ok
}
