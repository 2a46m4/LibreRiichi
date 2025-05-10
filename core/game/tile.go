package core

type Tile uint8

const (
	Manzu Tile = 0
	Pinzu Tile = 16
	Souzu Tile = 32

	Kazehai   Tile = 48
	EastTile  Tile = 48
	SouthTile Tile = 49
	WestTile  Tile = 50
	NorthTile Tile = 51

	Sangenpai Tile = 52
	White     Tile = 52
	Red       Tile = 53
	Green     Tile = 54

	DoraTile Tile = 64
	RedTile  Tile = 128
	Hidden   Tile = 254
	Invalid  Tile = 255
)

const (
	TileMask    Tile = 0b11 << 4
	ManzuBit         = 0
	PinzuBit         = 1
	SouzuBit         = 2
	HonourBit        = 3
	NumberMask       = 0b1111
	SpecialMask      = 0b11 << 6
)

func (s Tile) ClearRedOrDora() Tile {
	return s & ^(DoraTile | RedTile)
}

func (s Tile) IsInvalid() bool {
	return s == Invalid
}

func (s Tile) IsHidden() bool {
	return s == Hidden
}

func (s Tile) IsHonour() bool {
	return (s & TileMask) == HonourBit
}

func (s Tile) IsWind() bool {
	return s >= 48 && s <= 51
}

func (s Tile) IsDragon() bool {
	return s >= 52 && s <= 54
}

func (s Tile) IsManzu() bool {
	return s&(TileMask) == ManzuBit
}

func (s Tile) IsPinzu() bool {
	return s&(TileMask) == PinzuBit
}

func (s Tile) GetTileNumber() uint8 {
	return uint8(s & NumberMask)
}

func (s Tile) SetRedTile() Tile {
	return s | RedTile
}

func (s Tile) SetDoraTile() Tile {
	return s | DoraTile
}

func (s Tile) SetTileNumber(num uint8) Tile {
	return (s & (TileMask | SpecialMask)) | Tile(num)
}

// Return the list of tiles
func GetTileList() []Tile {
	tiles := make([]Tile, 136)
	tileItr := 0
	addFour := func(i Tile) {
		for range 4 {
			tiles[tileItr] = i
			tileItr += 1
		}

	}

	for i := Manzu; i < Manzu+10; i++ {
		addFour(i)
	}
	for i := Pinzu; i < Pinzu+10; i++ {
		addFour(i)
	}
	for i := Souzu; i < Souzu+10; i++ {
		addFour(i)
	}
	for i := Kazehai; i <= Green; i++ {
		addFour(i)
	}

	return tiles
}
