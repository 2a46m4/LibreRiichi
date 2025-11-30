package game

import (
	"log"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	"github.com/looplab/fsm"
)

type MahjongGame struct {
	gameState    GameState
	mahjongRound MahjongRound
	ordering     Ordering
	firstRound   bool
}

func NewMahjongGame() *MahjongGame {
	return &MahjongGame{
		gameState:  *InitGameState(),
		firstRound: true,
	}
}

func (game *MahjongGame) StartGame() (messages []MessageSendInfo, err error) {
	err = game.gameState.Transition("start-game", game, game.firstRound)
	sendInfoRaw, _ := game.gameState.GetReturn()
	sendInfo := sendInfoRaw.([]MessageSendInfo)
	ChangeToArenaIdx(sendInfo, game.ordering)
	return sendInfo, err
}

func (game *MahjongGame) StartRound() (msgs []MessageSendInfo, err error) {
	err = game.gameState.Transition("start-round", &game.mahjongRound, game.firstRound)
	sendInfoRaw, _ := game.gameState.GetReturn()
	sendInfo := sendInfoRaw.([]MessageSendInfo)
	ChangeToArenaIdx(sendInfo, game.ordering)
	return sendInfo, err
}

func (game *MahjongGame) HandleEvent(action Action, arenaIdx uint8) (msgs []MessageSendInfo, err error) {
    err = game.gameState.Transition("handle-event", &game.mahjongRound, &game.ordering, action, arenaIdx)
    _, isNoTransition := err.(fsm.NoTransitionError)
    if !isNoTransition {
	log.Println("Error occurred: ", err)
    } else {
	err = nil
    }
    sendInfoRaw, _ := game.gameState.GetReturn()
    sendInfo := sendInfoRaw.([]MessageSendInfo)
    ChangeToArenaIdx(sendInfo, game.ordering)
    return sendInfo, err
}

func (game *MahjongGame) RoundEnd() error {
	// Do the increment post-round
	game.mahjongRound.data.IncrementRound()

	if game.firstRound {
		game.firstRound = false
	}

	return nil
}

func (game *MahjongGame) RoundEnded() bool {
	return game.gameState.Current() == "in-game"
}

func (game *MahjongGame) GameEnd() error {
	return nil
}

func (game *MahjongGame) GameEnded() bool {
	return game.gameState.Current() == "finished-game"
}
