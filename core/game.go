package core

import (
	"errors"
	"log"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type MahjongState uint8

const (
	CURRENT_TURN MahjongState = iota
	CURRENT_TURN_PLAYED
	POST_TURN_PLAYED
	GAME_ENDED
)

var TossAction = Toss{TileToToss: Invalid}
var TossPotential = PotentialActionEvent{Actions: []Action{TossAction}}

// TODO: With the pending game actions stored in the game, we don't
// have to re-check a lot of the actions

type MahjongGameOld struct {
	Players []Player
	// Maps Player Index → Order
	PlayerToOrder []uint8
	// Maps Order → Player Index
	OrderToPlayer []uint8

	LiveWall []Tile
	Dora     []Tile
	UraDora  []Tile
	KanDraw  []Tile
	Tiles    [136]Tile

	// Current turn lasts until everyone has finished their possible actions
	// CurrentTurnOrder is in range (0, 4)
	CurrentTurnOrder uint8
	GameState        MahjongState
	RoundWind        Wind

	// Represents the next, undrawn tile
	TileIdx      uint8
	DoraRevealed uint8
	KansDrawn    uint8

	Results *GameResult // If game has finished, store the results here

	// The list of potential actions that need to be either taken or skipped
	// Need to attach a timer to them
	PendingActions []PendingAction

	// List of actions performed
	RecordedActions []ActionPerformed

	DiscardPiles [4][]Tile
}

// Returns the next events in the game, and if the game should end.
func (game *MahjongGameOld) GetNextEvent() (actions InfoList, shouldEnd bool) {
	log.Println("Getting next event")

	switch game.GameState {

	case CURRENT_TURN: // The current player can make a toss move
		log.Println("Current turn")

		// We should only reach this state when someone makes a post-turn action like pon.
		// Then the player only has the choice to discard or kan

		// TODO: Check if the player can make a kan
		private := PrivateMessage(game.currentPlayerIdx())
		private.Add(ArenaBoardEvent{BoardEvent: TossPotential})
		actions.Add(private)

		shouldEnd = false

	case CURRENT_TURN_PLAYED: // Get post-toss actions
		log.Println("Current turn played")

		// We should wait for all post toss actions to finish before moving to the next turn
		pendingActions, err := game.getPostTossActions()
		if err != nil {
			panic(err)
		}
		game.PendingActions = pendingActions

		if len(pendingActions) == 0 {
			game.GameState = POST_TURN_PLAYED
			return game.GetNextEvent()
		}

		// Group pending actions by player
		playerActions := make(map[uint8][]Action)
		for _, pendingAction := range pendingActions {
			playerActions[pendingAction.fromPlayer] = append(playerActions[pendingAction.fromPlayer], pendingAction.Action)
		}

		// Send all actions for each player in a single message
		for playerIdx, actionsForPlayer := range playerActions {
			private := PrivateMessage(playerIdx)
			private.Add(ArenaBoardEvent{
				BoardEvent: PotentialActionEvent{
					Actions: actionsForPlayer}})
			actions.Add(private)
		}

		shouldEnd = false

	case POST_TURN_PLAYED: // The post-toss has been played, we should progress to the next turn
		log.Println("Post turn played")

		tile, err := game.drawNewTile()

		game.GameState = CURRENT_TURN
		game.incrementTurn()

		if errors.Is(err, GameEndError{}) {
			game.GameState = GAME_ENDED
			return nil, true
		} else if err != nil {
			panic(err)
		}

		// Inform of draw and potential toss action
		partial := PartialMessage(game.currentPlayerIdx())
		partial.Add(ArenaBoardEvent{BoardEvent: PlayerActionEvent{
			Action:     Draw{DrawnTile: tile},
			FromPlayer: game.currentPlayerIdx(),
		}})
		actions.Add(partial)

		private := PrivateMessage(game.currentPlayerIdx())
		private.Add(ArenaBoardEvent{BoardEvent: TossPotential})
		actions.Add(private)

		// Get potential for performing a Riichi
		for _, discard := range game.currentPlayer().GetRiichiDiscards() {
			partial := PrivateMessage(game.currentPlayerIdx())
			partial.Add(ArenaBoardEvent{BoardEvent: PotentialActionEvent{Actions: []Action{Riichi{TileToRiichi: discard}}}})
			actions.Add(partial)
		}

		shouldEnd = false

	case GAME_ENDED:
		log.Println("Game ended")

		actions = nil
		shouldEnd = true
	default:
	}
	return actions, shouldEnd
}

// Updates the game state and returns the things to notify
// Additionally returns whether the move was valid
// Performs no validation of the action data structure
// func (game *MahjongGame) RespondToAction(action PlayerActionData) (InfoList, bool) {
// }

func (game *MahjongGameOld) HandleChii(chiiData Chii, fromPlayer uint8) (infos InfoList, err error) {
	onTile := chiiData.TileToChii
	chiiSequence := chiiData.TilesInHand

	last, err := game.lastTile()
	if err != nil || onTile != last {
		return nil, BadActionError{}
	}

	if game.GameState != CURRENT_TURN_PLAYED {
		return nil, BadActionError{}
	}
	if fromPlayer != game.nextPlayerIdx() {
		return nil, BadActionError{}
	}
	err = game.Players[fromPlayer].Chii(
		Tile(onTile),
		[2]Tile{
			Tile(chiiSequence[0]),
			Tile(chiiSequence[1]),
		})
	if err != nil {
		return nil, BadActionError{}
	}

	game.CurrentTurnOrder = fromPlayer

	return *infos.AddGlobalMessage(
		ArenaBoardEvent{BoardEvent: PlayerActionEvent{
			Action: Chii{
				TileToChii:  chiiData.TileToChii,
				TilesInHand: chiiData.TilesInHand},
			FromPlayer: fromPlayer}}), nil
}

func (game *MahjongGameOld) HandleKan(kanData Kan, fromPlayer uint8) (info InfoList, err error) {
	switch game.GameState {
	case CURRENT_TURN: // Ankan

		if fromPlayer != game.currentPlayerIdx() {
			return nil, BadActionError{}
		}

		err := game.Players[fromPlayer].Ankan(kanData.TileToKan)
		if err != nil {
			break
		}

		global := GlobalMessage()
		global.Add(ArenaBoardEvent{BoardEvent: PlayerActionEvent{Action: Kan{TileToKan: kanData.TileToKan}, FromPlayer: fromPlayer}})
		info.Add(global)

	case CURRENT_TURN_PLAYED: // Daiminkan
		if fromPlayer == game.currentPlayerIdx() {
			break
		}

		lastTile, err := game.lastTile()
		if err != nil || lastTile != kanData.TileToKan {
			break
		}

		err = game.Players[fromPlayer].Daiminkan(kanData.TileToKan)
		if err != nil {
			break
		}

		game.CurrentTurnOrder = fromPlayer

		global := GlobalMessage()
		global.Add(ArenaBoardEvent{
			BoardEvent: PlayerActionEvent{
				Action:     Kan{TileToKan: kanData.TileToKan},
				FromPlayer: fromPlayer}})
		info.Add(global)

	case POST_TURN_PLAYED: // Invalid
		err = BadActionError{}
	case GAME_ENDED: // Invalid
		err = BadActionError{}
	}
	return info, err
}

func (game *MahjongGameOld) HandlePon(ponData Pon, fromPlayer uint8) (info InfoList, err error) {
	last, err := game.lastTile()
	onTile := ponData.TileToPon
	if err != nil || onTile != last {
		return nil, BadActionError{}
	}
	if game.GameState != CURRENT_TURN_PLAYED {
		return nil, BadActionError{}
	}
	if fromPlayer != game.nextPlayerIdx() {
		return nil, BadActionError{}
	}

	err = game.Players[fromPlayer].Pon(onTile)
	if err != nil {
		return nil, BadActionError{}
	}
	game.CurrentTurnOrder = fromPlayer

	global := GlobalMessage()
	global.Add(ArenaBoardEvent{PlayerActionEvent{ponData, fromPlayer}})
	info.Add(global)

	return info, nil
}

func (game *MahjongGameOld) HandleRon(ronData Ron, fromPlayer uint8) (info InfoList, err error) {

	if fromPlayer == game.currentPlayerIdx() {
		return nil, BadActionError{}
	}
	_, err = game.findAction(ronData, fromPlayer)
	if err != nil {
		return nil, BadActionError{}
	}

	result, err := game.Players[fromPlayer].Ron(ronData.TileToRon)
	if err != nil {
		return nil, BadActionError{}
	}

	gameResult := GenerateGameResult(result, fromPlayer)
	// err = gameResult.Apply(game)
	if err != nil {
		return nil, BadActionError{}
	}

	game.Results = &gameResult
	game.GameState = GAME_ENDED
	global := GlobalMessage()
	global.Add(ArenaBoardEvent{PlayerActionEvent{ronData, fromPlayer}})
	info.Add(global)

	return info, nil
}

func (game *MahjongGameOld) HandleRiichi(riichiData Riichi, fromPlayer uint8) (info InfoList, err error) {

	tileDrawn, err := game.lastTile()
	if err != nil || riichiData.TileToRiichi != tileDrawn {
		return nil, BadActionError{}
	}
	if game.GameState != CURRENT_TURN {
		return nil, BadActionError{}
	}
	if fromPlayer != game.currentPlayerIdx() {
		return nil, BadActionError{}
	}

	err = game.Players[fromPlayer].Riichi(riichiData.TileToRiichi)
	if err != nil {
		return nil, BadActionError{}
	}

	game.GameState = CURRENT_TURN_PLAYED

	global := GlobalMessage()
	global.Add(ArenaBoardEvent{PlayerActionEvent{riichiData, fromPlayer}})
	info.Add(global)

	return info, nil
}

func (game *MahjongGameOld) HandleSkip(skipData Skip, fromPlayer uint8) (info InfoList, err error) {

	// We aren't finding the skip action itself but the action that is being skipped
	idx, err := game.findAction(
		skipData.ActionToSkip,
		fromPlayer,
	)
	if err != nil {
		return nil, BadActionError{}
	}
	// TODO: Check if the action is skippable, e.g. a toss is not skippable
	Remove(&game.PendingActions, idx)

	// Inform the player that it's skipped
	private := PrivateMessage(fromPlayer)
	private.Add(ArenaBoardEvent{PlayerActionEvent{skipData, fromPlayer}})
	info.Add(private)

	return info, nil

}

func (game *MahjongGameOld) HandleToss(tossData Toss, fromPlayer uint8) (info InfoList, err error) {

	onTile := tossData.TileToToss
	if game.GameState != CURRENT_TURN {
		return nil, BadActionError{}
	}
	if fromPlayer != game.currentPlayerIdx() {
		return nil, BadActionError{}
	}
	err = game.Players[fromPlayer].Toss(onTile)
	if err != nil {
		return nil, BadActionError{}
	}

	game.GameState = CURRENT_TURN_PLAYED
	global := GlobalMessage()
	global.Add(ArenaBoardEvent{PlayerActionEvent{tossData, fromPlayer}})
	info.Add(global)

	return info, nil
}

func (game *MahjongGameOld) HandleTsumo(tsumoData Tsumo, fromPlayer uint8) (info InfoList, err error) {

	if fromPlayer != game.currentPlayerIdx() {
		return nil, BadActionError{}
	}
	last, err := game.lastTile()
	if err != nil || tsumoData.TileToTsumo != last {
		return nil, BadActionError{}
	}

	result, err := game.Players[fromPlayer].Tsumo(tsumoData.TileToTsumo)
	if err != nil {
		return nil, BadActionError{}
	}

	gameResult := GenerateGameResult(result, fromPlayer)
	// err = gameResult.Apply(game)
	if err != nil {
		return nil, BadActionError{}
	}

	game.Results = &gameResult
	game.GameState = GAME_ENDED

	global := GlobalMessage()
	global.Add(ArenaBoardEvent{PlayerActionEvent{tsumoData, fromPlayer}})
	info.Add(global)

	return info, nil
}

func (game *MahjongGameOld) HandleDraw(drawData Draw, fromPlayer uint8) (InfoList, error) {
	panic("NYI")
}

// Checks the post-toss actions that can be made
func (game *MahjongGameOld) getPostTossActions() ([]PendingAction, error) {
	if game.GameState != CURRENT_TURN_PLAYED {
		return nil, errors.New("Incorrect state")
	}

	if len(game.PendingActions) != 0 {
		return game.PendingActions, nil
	}

	tileTossed, err := game.lastTile()
	log.Println("Tile just tossed: ", tileTossed)
	if err != nil {
		panic(err)
	}

	nextPlayerIdx := game.nextPlayerIdx()
	nextPlayer := game.Players[nextPlayerIdx]
	moves := make([]PendingAction, 0)

	// Helper that appends a potential move
	appendMove := func(action Action, forPlayer uint8) {
		moves = append(moves,
			PendingAction{action, forPlayer})
	}

	// Iterate through all possible combinations of Chii
	{
		tileNum := tileTossed.GetTileNumber()

		// Call when the chii move is valid
		appendChiiMove := func(chiiSequence [2]Tile) {
			appendMove(Chii{
				TileToChii:  tileTossed,
				TilesInHand: chiiSequence,
			}, nextPlayerIdx)
		}

		if tileNum <= 6 { // 6, 7, 8
			chiiSequence := [2]Tile{tileTossed + 1, tileTossed + 2}
			if nextPlayer.TestChii(tileTossed, chiiSequence) == nil {
				appendChiiMove(chiiSequence)
			}
		}
		if tileNum >= 2 { // 0, 1, 2
			chiiSequence := [2]Tile{tileTossed - 1, tileTossed - 2}
			if nextPlayer.TestChii(tileTossed, chiiSequence) == nil {
				appendChiiMove(chiiSequence)
			}
		}
		if tileNum >= 1 && tileNum <= 7 { // Middle
			chiiSequence := [2]Tile{tileTossed + 1, tileTossed - 1}
			if nextPlayer.TestChii(tileTossed, chiiSequence) == nil {
				appendChiiMove(chiiSequence)
			}
		}
	}

	// Iterate through all kans, pons, and rons
	for idx, player := range game.Players {
		if player.TestDaiminkan(tileTossed) == nil {
			appendMove(Kan{
				TileToKan: tileTossed,
			}, uint8(idx))
		}

		if player.TestPon(tileTossed) == nil {
			appendMove(Pon{
				TileToPon: tileTossed,
			}, uint8(idx))
		}

		if player.TestRon(tileTossed) == nil {
			appendMove(Ron{
				TileToRon: tileTossed,
			}, uint8(idx))
		}
	}

	return moves, nil
}

// Return the game results
func (MahjongGameOld) GetGameResults() (GameResult, error) {
	return GameResult{}, nil
}

// Returns the maximum amount of players
func (MahjongGameOld) GetMaxPlayers() int {
	return 4
}
