package core

type Hand struct {
	ClosedHand   []Tile
	Kans         []Tile
	Pons         []Tile
	Chiis        []Tile // Chiis are the start of the sequence
	HandOpen     bool
	HandInRiichi bool
}
