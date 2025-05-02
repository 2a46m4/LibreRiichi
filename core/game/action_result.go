package core

type Visibility uint8

const (
	PLAYER Visibility = 0
	GLOBAL            = 255
)

type ActionResult struct {
	ActionPerformed PlayerAction
	VisibleTo       Visibility
}
