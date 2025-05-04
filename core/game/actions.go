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
	Action          ActionType     `json:"action_type"`
	FromPlayer      uint8          `json:"from_player"`
	PotentialAction bool           `json:"potential_action"`
	Data            map[string]any `json:"on"`
}

type SetupType uint8

const (
	INITIAL_TILES SetupType = iota
	DORA
	STARTING_POINTS
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
