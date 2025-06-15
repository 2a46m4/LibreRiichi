package core

import (
	"errors"
	"fmt"
)

type MahjongState uint8

const (
	CURRENT_TURN MahjongState = iota
	CURRENT_TURN_PLAYED
	POST_TURN_PLAYED
	GAME_ENDED
)

// TODO: With the pending game actions stored in the game, we don't
// have to re-check a lot of the actions

type MahjongGame struct {
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
	PendingActions []MessageSendInfo
}

// ==================== ERRORS ====================
type GameEndError struct{}

func (GameEndError) Error() string { return "Game ended" }

// ==================== PRIVATE FUNCTIONS ====================

// Sets up the game and the tiles for the very start of the game
func (game *MahjongGame) setupGame() {
	game.Players = make([]Player, 4)
	game.PlayerToOrder = make([]uint8, 4)
	for i := range game.PlayerToOrder {
		game.PlayerToOrder[i] = uint8(i)
	}
	game.OrderToPlayer = make([]uint8, 0)

	PermuteArray(game.PlayerToOrder)
	for idx, order := range game.PlayerToOrder {
		game.OrderToPlayer[order] = uint8(idx)
	}
	game.Tiles = [136]Tile(GetTileList())
	PermuteArray(game.Tiles[:])
	game.CurrentTurnOrder = 3         // To initiate the first draw
	game.GameState = POST_TURN_PLAYED // To initiate the first draw
	game.RoundWind = East

	tileItr := 0
	for idx, order := range game.PlayerToOrder {
		player := &game.Players[idx]
		*player = Player{
			Points:   25000,
			SeatWind: Wind(order) + East,
		}
		player.FreshHand(game.Tiles[tileItr : tileItr+13])
		tileItr += 13
	}

	game.Dora = game.Tiles[tileItr : tileItr+5]
	game.DoraRevealed = 1
	tileItr += 5
	game.UraDora = game.Tiles[tileItr : tileItr+5]
	tileItr += 5
	game.KanDraw = game.Tiles[tileItr : tileItr+4]
	game.KansDrawn = 0
	tileItr += 4
	game.LiveWall = game.Tiles[tileItr:]
	game.TileIdx = 0

	game.Results = nil
	game.PendingActions = nil
}

func (game *MahjongGame) drawNewTile() (Tile, error) {
	if len(game.LiveWall) == 0 {
		return Invalid, GameEndError{}
	}
	if game.GameState != POST_TURN_PLAYED {
		return Invalid, errors.New("Not in right state to draw")
	}

	tile := game.LiveWall[game.TileIdx]
	game.TileIdx += 1
	return tile, nil
}

func (game MahjongGame) lastTile() (Tile, error) {
	if game.TileIdx == 0 {
		return Invalid, errors.New("No last tile")
	}

	switch game.GameState {
	case CURRENT_TURN:
		// The drawn tile
		return game.LiveWall[game.TileIdx-1], nil
	case CURRENT_TURN_PLAYED:
		// The tile just discarded
		return game.LiveWall[game.TileIdx-1], nil
	case POST_TURN_PLAYED:
		return game.LiveWall[game.TileIdx-1], nil
	case GAME_ENDED:
		return game.LiveWall[game.TileIdx-1], nil
	default:
		return Invalid, nil
	}
}

func (game *MahjongGame) currentPlayer() *Player {
	return &game.Players[game.currentPlayerIdx()]
}

func (game MahjongGame) currentPlayerIdx() uint8 {
	return game.OrderToPlayer[game.CurrentTurnOrder]
}

func (game MahjongGame) nextPlayerIdx() uint8 {
	return game.OrderToPlayer[game.CurrentTurnOrder+1]
}

func (game *MahjongGame) incrementTurn() {
	game.CurrentTurnOrder = (game.CurrentTurnOrder + 1) % 4
}

// Returns the index of the pending action
func (game MahjongGame) findAction(action ActionData, fromPlayer uint8) (int, error) {
	for i, pendingAction := range game.PendingActions {
		if pendingAction.ActionData == action &&
			pendingAction.uint8 == fromPlayer {
			return i, nil
		}
	}

	return 0, errors.New("Can't find action")
}

func encodeBoardEvent(eventType BoardEventType, data any) ArenaBoardEventData {
	return ArenaBoardEventData{
		BoardEvent{
			EventType: eventType,
			Data:      data,
		},
	}
}

func makeMessage(sendTo uint8, data ...ArenaBoardEventData) MessageSendInfo {
	return MessageSendInfo{
		Events: data,
		SendTo: Visibility(sendTo),
	}
}

// ==================== PUBLIC FUNCTIONS ====================

// Returns data to send to clients when a new game can be started, otherwise an error
func (game *MahjongGame) StartNewGame() ([][]Setup, error) {

	game.setupGame()
	setup := make([][]Setup, 4)
	startingPoints := [4]uint32{
		game.Players[0].Points,
		game.Players[1].Points,
		game.Players[2].Points,
		game.Players[3].Points,
	}

	for idx, player := range game.Players {
		setup[idx] = make([]Setup, 0, 7)
		setup[idx] = append(setup[idx],
			Setup{
				Type: INITIAL_TILES,
				Data: player.ClosedHand,
			},
			Setup{
				Type: DORA,
				Data: game.Dora[0],
			},
			Setup{
				Type: PLAYER_NUMBER,
				Data: uint8(idx),
			},
			Setup{
				Type: PLAYER_ORDER,
				Data: game.PlayerToOrder,
			},
			Setup{
				Type: ROUND_NUMBER,
				Data: 0,
			},
			Setup{
				Type: ROUND_WIND,
				Data: game.RoundWind,
			},
			Setup{
				Type: STARTING_POINTS,
				Data: startingPoints,
			})
	}

	return setup, nil
}

// Returns the next events in the game, and if the game should end.
func (game *MahjongGame) GetNextEvent() (actions []MessageSendInfo, shouldEnd bool) {
	switch game.GameState {

	case CURRENT_TURN: // The current player can make a toss move
		// We should only reach this state when someone makes a post-turn action like pon.
		// Then the player only has the choice to discard or kan

		// TODO: Check if the player can make a kan
		actions = []MessageSendInfo{
			makeMessage(
				game.currentPlayerIdx(),
				encodeBoardEvent(
					PotentialActionEventType,
					TossData{Invalid},
				)),
		}
		shouldEnd = false

	case CURRENT_TURN_PLAYED: // Get post-toss actions
		// We should wait for all post toss actions to finish before moving to the next turn
		var err error
		actions, err = game.getPostTossActions()
		if err != nil {
			panic(err)
		}
		game.PendingActions = actions

		if len(actions) == 0 {
			game.GameState = POST_TURN_PLAYED
			// TODO: Get next event again here?
		}

		shouldEnd = false

	case POST_TURN_PLAYED: // The post-toss has been played, we should progress to the next turn
		game.GameState = CURRENT_TURN
		game.incrementTurn()
		tile, err := game.drawNewTile()
		if errors.Is(err, GameEndError{}) {
			game.GameState = GAME_ENDED
			return nil, true
		}

		actions = []MessageSendInfo{
			makeMessage(
				game.currentPlayerIdx(),
				encodeBoardEvent(PlayerActionEventType, PlayerActionEventData{
					ActionData{
						ActionType: DRAW,
						Data:       DrawData{tile},
					},
					game.currentPlayerIdx(),
				})),

			// {
			// 	ActionPerformed: PlayerAction{
			// 		Action:     DRAW,
			// 		FromPlayer: game.currentPlayerIdx(),
			// 		Data:       DrawData{DrawnTile: tile},
			// 	},
			// 	IsPotential: false,
			// 	VisibleTo:   Visibility(game.currentPlayerIdx()),
			// },
			// {
			// 	ActionPerformed: PlayerAction{
			// 		Action:     TOSS,
			// 		FromPlayer: game.currentPlayerIdx(),
			// 		Data:       TossData{Invalid},
			// 	},
			// 	IsPotential: true,
			// 	VisibleTo:   Visibility(game.currentPlayerIdx()),
			// },
		}

		for _, discard := range game.currentPlayer().GetRiichiDiscards() {
			actions = append(actions, ActionResult{
				ActionPerformed: PlayerAction{
					Action:     RIICHI,
					FromPlayer: game.currentPlayerIdx(),
					Data: RiichiData{
						TileToRiichi: discard,
					},
				},
				IsPotential: true,
				VisibleTo:   Visibility(game.currentPlayerIdx()),
			})
		}

		shouldEnd = false

	case GAME_ENDED:
		actions = nil
		shouldEnd = true
	default:
	}
	return actions, shouldEnd
}

// Updates the game state and returns the things to notify
// Additionally returns whether the move was valid
// Performs no validation of the action data structure
func (game *MahjongGame) RespondToAction(action PlayerActionData) ([]MessageSendInfo, bool) {

	// TODO: Finish all the cases
	switch action.ActionType {
	case CHII:
		return game.handleChii(action)
	case KAN:
		return game.handleKan(action)
	case PON:
		return game.handlePon(action)
	case RIICHI:
		return game.handleRiichi(action)
	case RON:
		return game.handleRon(action)
	case SKIP:
		return game.handleSkip(action)
	case TOSS:
		return game.handleToss(action)
	case TSUMO:
		return game.handleTsumo(action)
	default:
		panic(fmt.Sprintf("unexpected core.ActionType: %#v", action.Action))
	}
}

func (game *MahjongGame) handleChii(action PlayerAction) ([]ActionResult, bool) {
	chiiData := action.Data.(ChiiData)

	onTile := chiiData.TileToChii
	chiiSequence := chiiData.TilesInHand

	last, err := game.lastTile()
	if err != nil || onTile != last {
		return nil, false
	}

	if game.GameState != CURRENT_TURN_PLAYED {
		return nil, false
	}
	if action.FromPlayer != game.nextPlayerIdx() {
		return nil, false
	}
	err = game.Players[action.FromPlayer].Chii(
		Tile(onTile),
		[2]Tile{
			Tile(chiiSequence[0]),
			Tile(chiiSequence[1]),
		})
	if err != nil {
		return nil, false
	}

	game.CurrentTurnOrder = action.FromPlayer

	return []ActionResult{
		{action, false, GLOBAL},
	}, true
}

func (game *MahjongGame) handleKan(action PlayerAction) (actions []ActionResult, validMove bool) {
	kanData := action.Data.(KanData)
	switch game.GameState {
	case CURRENT_TURN: // Ankan

		if action.FromPlayer != game.currentPlayerIdx() {
			break
		}

		err := game.Players[action.FromPlayer].Ankan(kanData.TileToKan)
		if err != nil {
			break
		}

		actions = []ActionResult{
			{
				ActionPerformed: action,
				IsPotential:     false,
				VisibleTo:       GLOBAL,
			},
		}
		validMove = true

	case CURRENT_TURN_PLAYED: // Daiminkan
		if action.FromPlayer == game.currentPlayerIdx() {
			break
		}

		lastTile, err := game.lastTile()
		if err != nil || lastTile != kanData.TileToKan {
			break
		}

		err = game.Players[action.FromPlayer].Daiminkan(kanData.TileToKan)
		if err != nil {
			break
		}

		game.CurrentTurnOrder = action.FromPlayer
		actions = []ActionResult{
			{action, false, GLOBAL},
		}
		validMove = true

	case POST_TURN_PLAYED: // Invalid
	case GAME_ENDED: // Invalid
	}
	return actions, validMove
}

func (game *MahjongGame) handlePon(action PlayerAction) (actions []ActionResult, validMove bool) {
	ponData := action.Data.(PonData)
	last, err := game.lastTile()
	onTile := ponData.TileToPon
	if err != nil || onTile != last {
		return nil, false
	}
	if game.GameState != CURRENT_TURN_PLAYED {
		return nil, false
	}
	if action.FromPlayer != game.nextPlayerIdx() {
		return nil, false
	}

	err = game.Players[action.FromPlayer].Pon(onTile)
	if err != nil {
		return nil, false
	}
	game.CurrentTurnOrder = action.FromPlayer

	return []ActionResult{
		{action, false, GLOBAL},
	}, true

}

func (game *MahjongGame) handleRon(action PlayerAction) ([]ActionResult, bool) {
	ronData := action.Data.(RonData)

	if action.FromPlayer == game.currentPlayerIdx() {
		return nil, false
	}
	_, err := game.findAction(action)
	if err != nil {
		return nil, false
	}

	result, err := game.Players[action.FromPlayer].Ron(ronData.TileToRon)
	if err != nil {
		return nil, false
	}

	gameResult := GenerateGameResult(result, action.FromPlayer)
	err = gameResult.Apply(game)
	if err != nil {
		return nil, false
	}

	game.Results = &gameResult
	game.GameState = GAME_ENDED
	return []ActionResult{
		{action, false, GLOBAL},
	}, true
}

func (game *MahjongGame) handleRiichi(action PlayerAction) ([]ActionResult, bool) {
	riichiData := action.Data.(RiichiData)

	tileDrawn, err := game.lastTile()
	if err != nil || riichiData.TileToRiichi != tileDrawn {
		return nil, false
	}
	if game.GameState != CURRENT_TURN {
		return nil, false
	}
	if action.FromPlayer != game.currentPlayerIdx() {
		return nil, false
	}

	err = game.Players[action.FromPlayer].Riichi(riichiData.TileToRiichi)
	if err != nil {
		return nil, false
	}

	game.GameState = CURRENT_TURN_PLAYED

	return []ActionResult{
		{action, false, GLOBAL},
	}, true
}

func (game *MahjongGame) handleSkip(action PlayerAction) ([]ActionResult, bool) {
	skipData := action.Data.(SkipData)

	// We aren't finding the skip action itself but the action that is being skipped
	idx, err := game.findAction(PlayerAction{
		Action:     skipData.ActionToSkip.Action,
		FromPlayer: action.FromPlayer,
		Data:       skipData.ActionToSkip.Data,
	})
	if err != nil {
		return nil, false
	}
	Remove(&game.PendingPostActions, idx)
	return []ActionResult{
		{action, false, Visibility(action.FromPlayer)},
	}, true
}

func (game *MahjongGame) handleToss(action PlayerAction) ([]ActionResult, bool) {
	tossData := action.Data.(TossData)

	onTile := tossData.TileToToss
	if game.GameState != CURRENT_TURN {
		return nil, false
	}
	if action.FromPlayer != game.currentPlayerIdx() {
		return nil, false
	}
	err := game.Players[action.FromPlayer].Toss(onTile)
	if err != nil {
		return nil, false
	}

	game.GameState = CURRENT_TURN_PLAYED
	actions := []ActionResult{
		{action, false, GLOBAL},
	}

	return actions, true
}

func (game *MahjongGame) handleTsumo(action PlayerAction) ([]ActionResult, bool) {
	tsumoData := action.Data.(TsumoData)

	if action.FromPlayer != game.currentPlayerIdx() {
		return nil, false
	}
	last, err := game.lastTile()
	if err != nil || tsumoData.TileToTsumo != last {
		return nil, false
	}

	result, err := game.Players[action.FromPlayer].Tsumo(tsumoData.TileToTsumo)
	if err != nil {
		return nil, false
	}

	gameResult := GenerateGameResult(result, action.FromPlayer)
	err = gameResult.Apply(game)
	if err != nil {
		return nil, false
	}

	game.Results = &gameResult
	game.GameState = GAME_ENDED
	return []ActionResult{
		{action, false, GLOBAL},
	}, true
}

// Checks the post-toss actions that can be made
func (game *MahjongGame) getPostTossActions() ([]MessageSendInfo, error) {
	if game.GameState != CURRENT_TURN_PLAYED {
		return nil, errors.New("Incorrect state")
	}

	if len(game.PendingPostActions) != 0 {
		return game.PendingPostActions, nil
	}

	tileTossed, err := game.lastTile()
	if err != nil {
		panic(err)
	}

	nextPlayerIdx := game.nextPlayerIdx()
	nextPlayer := game.Players[nextPlayerIdx]
	moves := make([]ActionResult, 0)

	// Helper that appends a potential move
	appendMove := func(action ActionType, forPlayer uint8, data ActionData) {
		moves = append(moves,
			ActionResult{
				ActionPerformed: PlayerAction{
					Action:     action,
					FromPlayer: forPlayer,
					Data:       data,
				},
				IsPotential: true,
				VisibleTo:   Visibility(forPlayer),
			})
	}

	// Iterate through all possible combinations of Chii
	{
		tileNum := tileTossed.GetTileNumber()

		// Call when the chii move is valid
		appendChiiMove := func(chiiSequence [2]Tile) {
			appendMove(CHII, nextPlayerIdx, ChiiData{
				TileToChii:  tileTossed,
				TilesInHand: chiiSequence,
			})
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
			appendMove(KAN, uint8(idx), KanData{
				TileToKan: tileTossed,
			})
		}

		if player.TestPon(tileTossed) == nil {
			appendMove(PON, uint8(idx), PonData{
				TileToPon: tileTossed,
			})
		}

		if player.TestRon(tileTossed) == nil {
			appendMove(RON, uint8(idx), RonData{
				TileToRon: tileTossed,
			})
		}
	}

	return moves, nil
}

// Return the game results
func (MahjongGame) GetGameResults() (GameResult, error) {
	return GameResult{}, nil
}

// Returns the maximum amount of players
func (MahjongGame) GetMaxPlayers() int {
	return 4
}
