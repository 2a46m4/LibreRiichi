package game

import (
	"log/slog"
	"os"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	gameresult "codeberg.org/ijnakashiar/LibreRiichi/core/game/game_result"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

// Essentially a thin wrapper over game state and changes the ordering
type MahjongGame struct {
	ordering Ordering
	*slog.Logger
	roundState RoundState
}

type startGame struct{}
type hasGameStarted struct{}
type startRound struct{}
type hasRoundStarted struct{}
type handleEvent struct {
	event Action
	from  uint8
}
type shouldContinueRound struct{}
type shouldGameEnd struct{}

// Retrieve round and game end messages that should be sent
type roundEndCleanup struct{}
type gameEndCleanup struct{}

func NewMahjongGame() *MahjongGame {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	game := &MahjongGame{
		ordering:   InitRandomOrdering(),
		Logger:     logger,
		roundState: InitRoundState(),
	}
	go game.roundState.gameLoop()
	return game
}

func (game *MahjongGame) StartGameAndRound() (messages []MessageSendInfo, err error) {
	game.roundState.Inbox <- startGame{}
	data := <-game.roundState.Reply
	switch data := data.(type) {
	case nil:
		break
	case error:
		return nil, data
	default:
		panic("Wrong type returned")
	}

	game.roundState.Inbox <- startRound{}
	data = <-game.roundState.Reply
	switch data := data.(type) {
	case []MessageSendInfo:
		ChangeToArenaIdx(data, game.ordering)
		return data, nil
	case error:
		return nil, data
	default:
		panic("Wrong type returned")
	}
}

func (game *MahjongGame) ContinueRound() (msgs []MessageSendInfo, err error) {
	game.roundState.Inbox <- startRound{}
	data := <-game.roundState.Reply
	switch data := data.(type) {
	case []MessageSendInfo:
		ChangeToArenaIdx(data, game.ordering)
		return data, nil
	case error:
		return nil, data
	default:
		panic("Wrong type returned")
	}
}

func (game *MahjongGame) HandleEvent(action Action, arenaIdx uint8) (msgs []MessageSendInfo, err error) {
	gameIdx := game.ordering.GameIdx(arenaIdx)
	game.roundState.Inbox <- handleEvent{
		event: action,
		from:  gameIdx,
	}
	data := <-game.roundState.Reply
	switch data := data.(type) {
	case []MessageSendInfo:
		ChangeToArenaIdx(data, game.ordering)
		return data, nil
	case error:
		return nil, data
	default:
		panic("Wrong type returned")
	}
}

func (game *MahjongGame) IsInGame() bool {
	game.roundState.Inbox <- hasGameStarted{}
	return (<-game.roundState.Reply).(bool)
}

func (game *MahjongGame) HasRoundEnded() bool {
	game.roundState.Inbox <- hasRoundStarted{}
	return !(<-game.roundState.Reply).(bool)
}

func (game *MahjongGame) RoundEndCleanup() (msgs []MessageSendInfo, err error) {
	game.roundState.Inbox <- roundEndCleanup{}
	data := <-game.roundState.Reply
	switch data := data.(type) {
	case gameresult.RoundResult:
		return sendToAll(RoundEndEvent{ RoundResult: data }) , nil
	case error:
		panic("Shouldn't be here")
	}
	// Do the increment post-round, but probably in the round itself
	return nil, nil
}

func (game *MahjongGame) ShouldContinueRound() bool {
	game.roundState.Inbox <- shouldContinueRound{}
	return (<-game.roundState.Reply).(bool)
}

// TODO: Send game results
func (game *MahjongGame) GameEndCleanup() (msgs []MessageSendInfo, err error) {
	return nil, nil
}

func sendToAll(events... BoardEvent) (msgs []MessageSendInfo) {
	for i := range uint8(4) {
		msgs = append(msgs, MessageSendInfo{
			Events: events,
			SendTo: i,
		})
	}
	return msgs
}
