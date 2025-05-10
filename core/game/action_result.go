package core

type Visibility uint8

const (
	PLAYER Visibility = 0
	GLOBAL            = 255
)

type ActionResult struct {
	ActionPerformed PlayerAction `json:"action"`
	// Whether this is an action that a player can take, not an action that a player took
	IsPotential bool       `json:"is_potential"`
	VisibleTo   Visibility `json:"visibility"`
}
