package game

import (
	"fmt"

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

func InitGame() *MahjongGame {
	return &MahjongGame {
		ScoringState: InitScoring(25000),
		TileState:    CreateNewRound(),
		RoundState:   RoundState{},
		WindState:    0,
		Ordering:     InitRandomOrdering(),
		TurnState:    InitTurnState(),
	}
}

func (game *MahjongGame) StartGame() ([]MessageSendInfo, error) {
	return nil, nil
}

func (game *MahjongGame) StartRound() ([]MessageSendInfo, error) {
	return nil, nil
}

func (game *MahjongGame) HandleEvent(action Action, arenaIdx uint8) ([]MessageSendInfo, error) {
	gameIdx := game.Ordering.ArenaToGame[arenaIdx]
	fmt.Println(gameIdx)

	return nil, nil
}

func (game *MahjongGame) RoundEnded() bool {
	return false
}

func (game *MahjongGame) GameEnded() bool {
	return false
}
