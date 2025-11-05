package game

type ScoringState struct {
	Points [4]uint32
}

func InitScoring(initial uint32) ScoringState {
	return ScoringState{
		Points: [4]uint32{initial, initial, initial, initial},
	}
}

