package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type MahjongGame struct {
	ScoringState
	TileState
	RoundState
	WindState
	Ordering
	TurnState
}

func InitGame() MahjongGame {
	return MahjongGame {
		ScoringState: InitScoring(25000),
		TileState:    CreateNewRound(),
		RoundState:   RoundState{},
		WindState:    0,
		Ordering:     InitRandomOrdering(),
		TurnState:    InitTurnState(),
	}
}

func DriveGame(action Action) error {
	
	return nil
}


