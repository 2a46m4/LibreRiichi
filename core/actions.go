package core

const (
	RON ActionType = iota
	TSUMO
	RIICHI
	TOSS
	SKIP
	PON
	KAN
	CHII
)

type ActionType uint8

type PlayerAction struct {
	Action     ActionType
	FromPlayer int
	On         any
}
