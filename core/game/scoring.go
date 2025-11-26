package game

type Scoring struct {
	// Game Index
	Points [4]uint32
}

func InitScoring(initial uint32) Scoring {
	return Scoring{
		Points: [4]uint32{initial, initial, initial, initial},
	}
}
