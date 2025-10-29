package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type GameActionHandler struct {
	game *MahjongGame
}

// NewGameActionHandler creates a new GameActionHandler with the given MahjongGame
func NewGameActionHandler(game *MahjongGame) GameActionHandler {
	return GameActionHandler{
		game: game,
	}
}

// HandleChii implements core.ActionHandler.
func (g GameActionHandler) HandleChii(action Chii, playerIdx uint8) (MessageSendInfo, error) {
	// Player calls Chii to claim a discarded tile
	// TilesInHand contains the two tiles from player's hand that combine with the discarded tile
	g.game.TileState.Chii(playerIdx, action.TilesInHand)

	// Create board event for the Chii call
	event := PlayerActionEvent{
		Action:     action,
		FromPlayer: playerIdx,
	}

	return MessageSendInfo{
		Events: []BoardEvent{event},
		SendTo: playerIdx,
	}, nil
}

// HandleDraw implements core.ActionHandler.
func (g GameActionHandler) HandleDraw(action Draw, playerIdx uint8) (MessageSendInfo, error) {
	// Player draws a tile from the wall
	g.game.TileState.Draw(playerIdx)

	// Create board event for the draw
	event := PlayerActionEvent{
		Action:     action,
		FromPlayer: playerIdx,
	}

	return MessageSendInfo{
		Events: []BoardEvent{event},
		SendTo: playerIdx,
	}, nil
}

// HandleKan implements core.ActionHandler.
func (g GameActionHandler) HandleKan(action Kan, playerIdx uint8) (MessageSendInfo, error) {
	// Player calls Kan - this could be several types:
	// - DaiminKan (great kan from another player's discard)
	// - ShouminKan (small kan from self-drawn tile)
	// - AnKan (closed kan from hand only)

	// For now, implement as ShouminKan (from self)
	g.game.TileState.ShouminKan(playerIdx, action.TileToKan)

	// Create board event for the Kan call
	event := PlayerActionEvent{
		Action:     action,
		FromPlayer: playerIdx,
	}

	return MessageSendInfo{
		Events: []BoardEvent{event},
		SendTo: playerIdx,
	}, nil
}

// HandlePon implements core.ActionHandler.
func (g GameActionHandler) HandlePon(action Pon, playerIdx uint8) (MessageSendInfo, error) {
	// Player calls Pon to claim a discarded tile
	// Need to determine which player discarded the tile (previous player in turn order)
	againstPlayer := (playerIdx + 3) % 4 // player who discarded the tile

	g.game.TileState.Pon(playerIdx, againstPlayer)

	// Create board event for the Pon call
	event := PlayerActionEvent{
		Action:     action,
		FromPlayer: playerIdx,
	}

	return MessageSendInfo{
		Events: []BoardEvent{event},
		SendTo: playerIdx,
	}, nil
}

// HandleRiichi implements core.ActionHandler.
func (g GameActionHandler) HandleRiichi(action Riichi, playerIdx uint8) (MessageSendInfo, error) {
	// Player declares Riichi (ready hand)
	g.game.TileState.Riichi(playerIdx, action.TileToRiichi)

	// Create board event for the Riichi declaration
	event := PlayerActionEvent{
		Action:     action,
		FromPlayer: playerIdx,
	}

	return MessageSendInfo{
		Events: []BoardEvent{event},
		SendTo: playerIdx,
	}, nil
}

// HandleRon implements core.ActionHandler.
func (g GameActionHandler) HandleRon(action Ron, playerIdx uint8) (MessageSendInfo, error) {
	// Player wins by Ron (claiming another player's discard)
	g.game.TileState.Ron(playerIdx, action.TileToRon)

	// Create board event for the Ron win
	event := PlayerActionEvent{
		Action:     action,
		FromPlayer: playerIdx,
	}

	return MessageSendInfo{
		Events: []BoardEvent{event},
		SendTo: playerIdx,
	}, nil
}

// HandleSkip implements core.ActionHandler.
func (g GameActionHandler) HandleSkip(action Skip, playerIdx uint8) (MessageSendInfo, error) {
	// Player skips an action (passes on calling a tile)
	// This is usually used when a player could call a discarded tile but chooses not to

	// Create board event for the skip action
	event := PlayerActionEvent{
		Action:     action,
		FromPlayer: playerIdx,
	}

	return MessageSendInfo{
		Events: []BoardEvent{event},
		SendTo: playerIdx,
	}, nil
}

// HandleToss implements core.ActionHandler.
func (g GameActionHandler) HandleToss(action Toss, playerIdx uint8) (MessageSendInfo, error) {
	// Player discards a tile
	g.game.TileState.Discard(playerIdx, action.TileToToss)

	// Create board event for the discard
	event := PlayerActionEvent{
		Action:     action,
		FromPlayer: playerIdx,
	}

	return MessageSendInfo{
		Events: []BoardEvent{event},
		SendTo: playerIdx,
	}, nil
}

// HandleTsumo implements core.ActionHandler.
func (g GameActionHandler) HandleTsumo(action Tsumo, playerIdx uint8) (MessageSendInfo, error) {
	// Player wins by Tsumo (self-drawn winning tile)
	g.game.TileState.Tsumo(playerIdx, action.TileToTsumo)

	// Create board event for the Tsumo win
	event := PlayerActionEvent{
		Action:     action,
		FromPlayer: playerIdx,
	}

	return MessageSendInfo{
		Events: []BoardEvent{event},
		SendTo: playerIdx,
	}, nil
}

func x() {
	var x ActionHandler[MessageSendInfo, uint8] = GameActionHandler{}
}
