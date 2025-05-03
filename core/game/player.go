package core

type Player struct {
	ClosedHand []Tile
	Kans       []Tile
	Pons       []Tile
	Chiis      []Tile

	InRiichi bool
	HandOpen bool

	Points   uint32
	SeatWind Wind
}

func (player *Player) FreshHand(tiles []Tile) {
	if len(tiles) != 13 {
		panic("Not equal to 13")
	}

	player.ClosedHand = make([]Tile, 14)
	copy(player.ClosedHand, tiles)
	player.Kans = nil
	player.Pons = nil
	player.Chiis = nil

	player.InRiichi = false
	player.HandOpen = false
}
