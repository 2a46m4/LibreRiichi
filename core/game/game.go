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

func (game *MahjongGame) InitGame() {
	*game = MahjongGame {
		ScoringState: InitScoring(25000),
		TileState:    CreateNewRound(),
		RoundState:   RoundState{},
		WindState:    0,
		Ordering:     InitRandomOrdering(),
		TurnState:    InitTurnState(),
	}
}

func (game *MahjongGame) StartGame() (MessageSendInfo, error) {
	
}

func (game *MahjongGame) HandleEvent(action Action) (MessageSendInfo, error) {
	return MessageSendInfo{}, nil
}


