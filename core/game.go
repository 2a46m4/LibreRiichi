package core

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"

	"errors"
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
	PendingActions []PendingAction
}

type PendingAction struct {
	ActionData
	fromPlayer uint8
}

// ==================== ERRORS ====================
type GameEndError struct{}

func (GameEndError) Error() string { return "Game ended" }

type BadActionError struct{}

func (BadActionError) Error() string { return "Bad action" }

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
	for idx, pendingAction := range game.PendingActions {
		if pendingAction.ActionData == action &&
			pendingAction.fromPlayer == fromPlayer {
			return idx, nil
		}
	}

	return 0, errors.New("Can't find action")
}

func encodeBoardEvent(eventType BoardEventType, data any) ArenaBoardEventData {
	return ArenaBoardEventData{
		BoardEvent: BoardEvent{
			EventType: eventType,
			Data:      data,
		},
	}
}

func encodePotentialAction(data ActionData) ArenaBoardEventData {
	return encodeBoardEvent(
		PotentialActionEventType,
		PotentialActionEventData{ActionData: data},
	)
}

func encodePlayerAction(data ActionData, fromPlayer uint8) ArenaBoardEventData {
	return encodeBoardEvent(
		PlayerActionEventType,
		PlayerActionEventData{ActionData: data, FromPlayer: fromPlayer},
	)
}

func makeMessage(visibility Visibility, sendTo uint8, data ...ArenaBoardEventData) MessageSendInfo {
	return MessageSendInfo{
		Events:     data,
		Visibility: visibility,
		SendTo:     sendTo,
	}
}

func makeGlobalMessage(data ...ArenaBoardEventData) MessageSendInfo {
	return MessageSendInfo{
		Events:     data,
		Visibility: GLOBAL,
		SendTo:     0,
	}
}

func globalPlayerAction(data ActionData, fromPlayer uint8) MessageSendInfo {
	return MessageSendInfo{
		Events: []ArenaBoardEventData{
			encodePlayerAction(data, fromPlayer),
		},
		Visibility: GLOBAL,
		SendTo:     0,
	}
}

func privatePlayerAction(data ActionData, fromPlayer uint8) MessageSendInfo {
	return MessageSendInfo{
		Events: []ArenaBoardEventData{
			encodePlayerAction(data, fromPlayer),
		},
		Visibility: PLAYER,
		SendTo:     fromPlayer,
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
				PLAYER,
				game.currentPlayerIdx(),
				encodePotentialAction(
					ActionData{
						ActionType: TOSS,
						Data:       TossData{Invalid},
					},
				)),
		}
		shouldEnd = false

	case CURRENT_TURN_PLAYED: // Get post-toss actions
		// We should wait for all post toss actions to finish before moving to the next turn
		pendingActions, err := game.getPostTossActions()
		if err != nil {
			panic(err)
		}
		game.PendingActions = pendingActions

		if len(pendingActions) == 0 {
			game.GameState = POST_TURN_PLAYED
			// TODO: Get next event again here?
		}

		for _, pendingAction := range pendingActions {
			actions = append(actions, makeMessage(
				PLAYER,
				pendingAction.fromPlayer,
				encodePotentialAction(pendingAction.ActionData),
			))
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
				PARTIAL,
				game.currentPlayerIdx(),
				encodeBoardEvent(
					PlayerActionEventType,
					PlayerActionEventData{
						ActionData: ActionData{
							ActionType: DRAW,
							Data:       DrawData{tile},
						},
						FromPlayer: game.currentPlayerIdx(),
					})),
			makeMessage(
				PLAYER,
				game.currentPlayerIdx(),
				encodePotentialAction(
					ActionData{
						ActionType: TOSS,
						Data:       TossData{Invalid},
					},
				)),
		}

		// For performing a Riichi
		for _, discard := range game.currentPlayer().GetRiichiDiscards() {
			actions = append(actions,
				makeMessage(
					PLAYER,
					game.currentPlayerIdx(),
					encodePotentialAction(
						ActionData{
							ActionType: RIICHI,
							Data:       RiichiData{discard},
						},
					),
				))
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
// func (game *MahjongGame) RespondToAction(action PlayerActionData) ([]MessageSendInfo, bool) {
// }

func (game *MahjongGame) HandleChii(chiiData ChiiData, fromPlayer uint8) ([]MessageSendInfo, error) {
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

	return []MessageSendInfo{
		globalPlayerAction(ActionData{CHII, chiiData}, fromPlayer),
	}, nil
}

func (game *MahjongGame) HandleKan(kanData KanData, fromPlayer uint8) (info []MessageSendInfo, err error) {
	switch game.GameState {
	case CURRENT_TURN: // Ankan

		if fromPlayer != game.currentPlayerIdx() {
			return nil, BadActionError{}
		}

		err := game.Players[fromPlayer].Ankan(kanData.TileToKan)
		if err != nil {
			break
		}

		info = []MessageSendInfo{
			makeGlobalMessage(
				encodePlayerAction(ActionData{KAN, kanData}, fromPlayer),
			),
		}

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
		info = []MessageSendInfo{
			globalPlayerAction(ActionData{KAN, kanData}, fromPlayer),
		}

	case POST_TURN_PLAYED: // Invalid
		err = BadActionError{}
	case GAME_ENDED: // Invalid
		err = BadActionError{}
	}
	return info, err
}

func (game *MahjongGame) HandlePon(ponData PonData, fromPlayer uint8) ([]MessageSendInfo, error) {
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

	return []MessageSendInfo{
		globalPlayerAction(ActionData{PON, ponData}, fromPlayer),
	}, nil

}

func (game *MahjongGame) HandleRon(ronData RonData, fromPlayer uint8) ([]MessageSendInfo, error) {

	if fromPlayer == game.currentPlayerIdx() {
		return nil, BadActionError{}
	}
	_, err := game.findAction(ActionData{RON, ronData}, fromPlayer)
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
	return []MessageSendInfo{
		globalPlayerAction(ActionData{RON, ronData}, fromPlayer),
	}, nil
}

func (game *MahjongGame) HandleRiichi(riichiData RiichiData, fromPlayer uint8) ([]MessageSendInfo, error) {

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

	return []MessageSendInfo{
		globalPlayerAction(ActionData{RIICHI, riichiData}, fromPlayer),
	}, nil
}

func (game *MahjongGame) HandleSkip(skipData SkipData, fromPlayer uint8) ([]MessageSendInfo, error) {

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
	return []MessageSendInfo{
		privatePlayerAction(ActionData{SKIP, skipData}, fromPlayer),
	}, nil
}

func (game *MahjongGame) HandleToss(tossData TossData, fromPlayer uint8) ([]MessageSendInfo, error) {

	onTile := tossData.TileToToss
	if game.GameState != CURRENT_TURN {
		return nil, BadActionError{}
	}
	if fromPlayer != game.currentPlayerIdx() {
		return nil, BadActionError{}
	}
	err := game.Players[fromPlayer].Toss(onTile)
	if err != nil {
		return nil, BadActionError{}
	}

	game.GameState = CURRENT_TURN_PLAYED
	return []MessageSendInfo{
		globalPlayerAction(ActionData{TOSS, tossData}, fromPlayer),
	}, nil
}

func (game *MahjongGame) HandleTsumo(tsumoData TsumoData, fromPlayer uint8) ([]MessageSendInfo, error) {

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
	return []MessageSendInfo{
		globalPlayerAction(ActionData{TSUMO, tsumoData}, fromPlayer),
	}, nil
}

func (game *MahjongGame) HandleDraw(drawData DrawData, fromPlayer uint8) ([]MessageSendInfo, error) {
	panic("NYI")
}

// Checks the post-toss actions that can be made
func (game *MahjongGame) getPostTossActions() ([]PendingAction, error) {
	if game.GameState != CURRENT_TURN_PLAYED {
		return nil, errors.New("Incorrect state")
	}

	if len(game.PendingActions) != 0 {
		return game.PendingActions, nil
	}

	tileTossed, err := game.lastTile()
	if err != nil {
		panic(err)
	}

	nextPlayerIdx := game.nextPlayerIdx()
	nextPlayer := game.Players[nextPlayerIdx]
	moves := make([]PendingAction, 0)

	// Helper that appends a potential move
	appendMove := func(action ActionData, forPlayer uint8) {
		moves = append(moves,
			PendingAction{action, forPlayer})
	}

	// Iterate through all possible combinations of Chii
	{
		tileNum := tileTossed.GetTileNumber()

		// Call when the chii move is valid
		appendChiiMove := func(chiiSequence [2]Tile) {
			appendMove(ActionData{CHII, ChiiData{
				TileToChii:  tileTossed,
				TilesInHand: chiiSequence,
			}}, nextPlayerIdx)
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
			appendMove(ActionData{KAN, KanData{
				TileToKan: tileTossed,
			}}, uint8(idx))
		}

		if player.TestPon(tileTossed) == nil {
			appendMove(ActionData{PON, PonData{
				TileToPon: tileTossed,
			}}, uint8(idx))
		}

		if player.TestRon(tileTossed) == nil {
			appendMove(ActionData{RON, RonData{
				TileToRon: tileTossed,
			}}, uint8(idx))
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

func GetAltMessage(msg ArenaMessage) (altMsg ArenaMessage, err error) {
	if msg.MessageType != ArenaBoardEventType {
		return altMsg, errors.New("Not correct type")
	}
	eventData := msg.Data.(ArenaBoardEventData)
	BoardEventDispatch(AltMessageHandler{}, eventData.BoardEvent)

	return ArenaMessage{}, nil
}
