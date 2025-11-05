package game

import (
	"errors"
	"fmt"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

// TileStateValidator validates that actions performed on the game state are valid
// according to Mahjong rules and game state constraints
type TileStateValidator struct {
	tileState *TileState
	turnState *TurnState
	gameState *GameState
}

func NewTileStateValidator(tileState *TileState, turnState *TurnState, gameState *GameState) *TileStateValidator {
	return &TileStateValidator{
		tileState: tileState,
		turnState: turnState,
		gameState: gameState,
	}
}

func (v *TileStateValidator) ValidateGameState(action Action, playerIdx uint8) error {
	switch a := action.(type) {
	case Draw:
		return v.gameState.Try(DRAW_TRANSITION)
	case Toss:
		return v.gameState.Try(TOSS_TRANSITION)
	case Pon:
		return v.gameState.Try(NAKI_CALLED_TRANSITION)
	case Kan:
		return v.gameState.Try(NAKI_CALLED_TRANSITION)
	case Chii:
		return v.gameState.Try(NAKI_CALLED_TRANSITION)
	case Riichi:
		return v.gameState.Try(NAKI_CALLED_TRANSITION)
	case Tsumo:
		return v.gameState.Try(TSUMO_TRANSITION)
	case Ron:
		return v.gameState.Try(RON_TRANSITION)
	case Skip:
		return v.gameState.Try(SKIP_TRANSITION)
	default:
		return fmt.Errorf("unknown action type: %T", action)
	}
}

// ValidateAction validates that an action can be performed by the given player
// Returns an error if the action is invalid
func (v *TileStateValidator) ValidateAction(action Action, playerIdx uint8) error {
	// First check if it's the player's turn and the action is valid for turn state
	if err := v.validateTurnState(action, playerIdx); err != nil {
		return err
	}

	// Then check specific action validation rules
	switch a := action.(type) {
	case Draw:
		return v.validateDraw(playerIdx)
	case Toss:
		return v.validateToss(a, playerIdx)
	case Pon:
		return v.validatePon(playerIdx)
	case Kan:
		return v.validateKan(a, playerIdx)
	case Chii:
		return v.validateChii(a, playerIdx)
	case Riichi:
		return v.validateRiichi(a, playerIdx)
	case Tsumo:
		return v.validateTsumo(a, playerIdx)
	case Ron:
		return v.validateRon(a, playerIdx)
	case Skip:
		return v.validateSkip(a, playerIdx)
	default:
		return fmt.Errorf("unknown action type: %T", action)
	}
}

func (v *TileStateValidator) validateTurnState(action Action, playerIdx uint8) error {
	return v.turnState.ProcessAction(action, playerIdx)
}

func (v *TileStateValidator) validateDraw(playerIdx uint8) error {

	if v.tileState.LiveWall.End() {
		return errors.New("cannot draw: no tiles remaining in live wall")
	}

	hand := v.tileState.Hands[playerIdx]
	if hand.FullHand() {
		return errors.New("cannot draw: hand already has too many tiles")
	}

	return nil
}

// validateToss validates that a player can discard the given tile
func (v *TileStateValidator) validateToss(action Toss, playerIdx uint8) error {
	tile := action.TileToToss

	hand := v.tileState.Hands[playerIdx]
	if !hand.ClosedHand.HasTile(tile) {
		return fmt.Errorf("cannot discard %s: tile not in hand", tile)
	}

	if v.turnState.GetCurrentPlayer() != playerIdx {
		return errors.New("cannot discard: not player's turn")
	}

	return nil
}

// validatePon validates that a player can call Pon
func (v *TileStateValidator) validatePon(playerIdx uint8) error {
	// Check if player is not calling Pon on their own discard
	if v.turnState.GetCurrentPlayer() == playerIdx {
		return errors.New("cannot call Pon on own discard")
	}

	// Check if player has at least two of the discarded tile in their hand
	// This is a basic check - the actual tile availability would be checked by the caller
	hand := v.tileState.Hands[playerIdx]

	// Get the most recently discarded tile to check against
	discardPile := v.tileState.DiscardPile[v.turnState.GetCurrentPlayer()]
	if discardPile.IsEmpty() {
		return errors.New("cannot call Pon: no discarded tile available")
	}

	// Basic check that player has room for an open meld
	if hand.ClosedHand.index < 2 {
		return errors.New("cannot call Pon: not enough tiles in hand")
	}

	return nil
}

// validateKan validates that a player can call Kan
func (v *TileStateValidator) validateKan(action Kan, playerIdx uint8) error {
	tile := action.TileToKan

	// For different types of Kan, we need different validation
	// This is a basic validation - more specific validation would be needed for each Kan type

	// Check if player has enough tiles in hand for the Kan
	hand := v.tileState.Hands[playerIdx]
	if hand.ClosedHand.index < 3 {
		return errors.New("cannot call Kan: not enough tiles in hand")
	}

	return nil
}

// validateChii validates that a player can call Chii with the given tiles
func (v *TileStateValidator) validateChii(action Chii, playerIdx uint8) error {
	// Chii can only be called by the next player (playerIdx + 1) % 4
	expectedPlayer := (v.turnState.GetCurrentPlayer() + 1) % 4
	if playerIdx != expectedPlayer {
		return errors.New("cannot call Chii: only the next player can call Chii")
	}

	// Check if player has the specified tiles in their hand
	hand := v.tileState.Hands[playerIdx]
	tiles := action.TilesInHand

	for _, tile := range tiles {
		if !hand.ClosedHand.HasTile(tile) {
			return fmt.Errorf("cannot call Chii: tile %s not in hand", tile)
		}
	}

	// Check if tiles form a valid sequence (basic check)
	// This would normally involve more complex sequence validation
	if !v.isValidSequence(tiles) {
		return errors.New("cannot call Chii: tiles do not form a valid sequence")
	}

	return nil
}

// validateRiichi validates that a player can declare Riichi
func (v *TileStateValidator) validateRiichi(action Riichi, playerIdx uint8) error {
	tile := action.TileToRiichi

	// Check if player hasn't already declared Riichi
	hand := v.tileState.Hands[playerIdx]
	if hand.InRiichi {
		return errors.New("cannot declare Riichi: already in Riichi")
	}

	// Check if player has the tile they're discarding
	if !hand.ClosedHand.HasTile(tile) {
		return fmt.Errorf("cannot declare Riichi: tile %s not in hand", tile)
	}

	// Check if hand is in tenpai (ready to win)
	// This would require complex tenpai detection - for now, basic check
	if hand.ClosedHand.index != 14 {
		return errors.New("cannot declare Riichi: hand not in proper state")
	}

	return nil
}

// validateTsumo validates that a player can win by Tsumo
func (v *TileStateValidator) validateTsumo(action Tsumo, playerIdx uint8) error {
	tile := action.TileToTsumo

	// Check if it's the player's turn
	if v.turnState.GetCurrentPlayer() != playerIdx {
		return errors.New("cannot win by Tsumo: not player's turn")
	}

	// Check if player has the drawn tile
	hand := v.tileState.Hands[playerIdx]
	if !hand.ClosedHand.HasTile(tile) {
		return fmt.Errorf("cannot win by Tsumo: tile %s not in hand", tile)
	}

	// Check if hand forms a winning combination
	// This would require complex winning hand detection
	if !v.isWinningHand(hand) {
		return errors.New("cannot win by Tsumo: hand is not a winning combination")
	}

	return nil
}

// validateRon validates that a player can win by Ron
func (v *TileStateValidator) validateRon(action Ron, playerIdx uint8) error {
	tile := action.TileToRon

	// Check if player is not trying to win on their own discard
	if v.turnState.GetCurrentPlayer() == playerIdx {
		return errors.New("cannot win by Ron: cannot win on own discard")
	}

	// Check if the discarded tile exists and can be claimed
	discardPile := v.tileState.DiscardPile[v.turnState.GetCurrentPlayer()]
	if discardPile.IsEmpty() {
		return errors.New("cannot win by Ron: no discarded tile available")
	}

	// Check if adding the tile creates a winning hand
	hand := v.tileState.Hands[playerIdx]
	tempHand := hand // Create a copy for validation
	tempHand.ClosedHand.Add(tile)

	if !v.isWinningHand(tempHand) {
		return errors.New("cannot win by Ron: does not create a winning hand")
	}

	return nil
}

// validateSkip validates that a player can skip an action
func (v *TileStateValidator) validateSkip(action Skip, playerIdx uint8) error {
	// Skip is generally always valid as long as there's an action to skip
	if action.ActionToSkip == nil {
		return errors.New("cannot skip: no action specified")
	}

	return nil
}

// isValidSequence checks if two tiles can form a valid sequence with a third tile
func (v *TileStateValidator) isValidSequence(tiles [2]Tile) bool {
	// Basic sequence validation - this would need to be more comprehensive
	// For now, just check that tiles are not the same and are from the same suit
	if tiles[0].Suit != tiles[1].Suit {
		return false
	}

	// Check if tiles can form part of a sequence
	diff := tiles[1].Value - tiles[0].Value
	if diff < 1 || diff > 2 {
		return false
	}

	return true
}

// isWinningHand checks if a hand forms a winning combination
// This is a placeholder - full implementation would be complex
func (v *TileStateValidator) isWinningHand(hand Hand) bool {
	// For now, just basic checks - full winning hand detection is complex
	// Would need to check for 4 melds + 1 pair, special hands, etc.
	return hand.ClosedHand.index == 14
}
