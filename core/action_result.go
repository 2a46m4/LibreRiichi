package core

type Visibility uint8

const (
	PLAYER Visibility = 0
	GLOBAL            = 255
)

// Represents a change in the board state. It can be either a change that can happen or a change that has occured. Usually it is an action that a player has taken.
type ActionResult struct {
	// Whether this is an action that a player can take, not an action that a player took
	VisibleTo Visibility `json:"visibility"`
}
