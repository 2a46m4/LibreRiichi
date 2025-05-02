package core

type ActionType uint8

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

type PlayerAction struct {
	Action     ActionType
	FromPlayer uint8
	On         any
}

type SetupType uint8

const (
	INITIAL_TILES SetupType = iota
	DORA
	STARTING_POINTS_SELF
	STARTING_POINTS_OTHER
	PLAYER_NUMBER
	PLAYER_ORDER
	ROUND_WIND
	ROUND_NUMBER
)

type Setup struct {
	Type     SetupType
	ToPlayer uint8
	Data     any
}
