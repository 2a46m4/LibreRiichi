package core

type Wind Tile

const (
	East  Wind = Wind(EastTile)
	South      = Wind(SouthTile)
	West       = Wind(WestTile)
	North      = Wind(NorthTile)
)

func SameWind(wind Wind, tile Tile) bool {
	return wind == Wind(tile)
}
