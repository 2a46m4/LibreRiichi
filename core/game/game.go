package game

import (
	"log"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	"github.com/looplab/fsm"
)

type MahjongGame struct {
	gameState    GameState
	ordering     Ordering
}

func NewMahjongGame() *MahjongGame {
	return &MahjongGame{
		gameState:  *InitGameState(),
	}
}

// We should be calling the function inside
func (game *MahjongGame) StartGame() (messages []MessageSendInfo, err error) {
	err = game.gameState.Transition("start-game")
	sendInfo, _ := game.gameState.GetReturn()
	ChangeToArenaIdx(sendInfo, game.ordering)
	return sendInfo, err
}

func (game *MahjongGame) StartRound() (msgs []MessageSendInfo, err error) {
	err = game.gameState.Transition("start-round", true)
	sendInfo, _ := game.gameState.GetReturn()
	ChangeToArenaIdx(sendInfo, game.ordering)
	return sendInfo, err
}

func (game *MahjongGame) ContinueRound() (msgs []MessageSendInfo, err error) {
    err = game.gameState.Transition("start-round", false)
    sendInfo, _ := game.gameState.GetReturn()
    ChangeToArenaIdx(sendInfo, game.ordering)
    return nil, nil
}


func (game *MahjongGame) HandleEvent(action Action, arenaIdx uint8) (msgs []MessageSendInfo, err error) {
    gameIdx := game.ordering.GameIdx(arenaIdx)
    err = game.gameState.Transition("handle-event", &game.mahjongRound, action, gameIdx)
    _, isNoTransition := err.(fsm.NoTransitionError)
    if !isNoTransition {
	log.Println("Error occurred: ", err)
    } else {
	err = nil
    }
    sendInfo, _ := game.gameState.GetReturn()
    ChangeToArenaIdx(sendInfo, game.ordering)
    return sendInfo, err
}

func (game *MahjongGame) RoundEnd() (msgs []MessageSendInfo, err error) {
    // Do the increment post-round
    // TODO
    // game.mahjongRound.roundState


    return nil, nil
}

func (game *MahjongGame) RoundEnded() bool {
	return game.gameState.Current() == "in-game"
}

func (game *MahjongGame) GameEnd() (msgs []MessageSendInfo, err error) {
	return nil, nil
}

func (game *MahjongGame) GameEnded() bool {
	return game.gameState.Current() == "finished-game"
}
