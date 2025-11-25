package main

import (
	"fmt"
	"codeberg.org/ijnakashiar/LibreRiichi/core/game"
)

func main() {
	fmt.Println(game.InitGameState().GenerateGraphs())
}
