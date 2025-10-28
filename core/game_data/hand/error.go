package core

type TooManyTilesErr struct{}

func (TooManyTilesErr) Error() string {
	return "Too many tiles in hand"
}

type TooLittleTilesErr struct{}

func (TooLittleTilesErr) Error() string {
	return "Too little tiles in hand"
}
