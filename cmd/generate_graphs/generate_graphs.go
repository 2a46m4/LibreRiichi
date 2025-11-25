package main

import (
	"codeberg.org/ijnakashiar/LibreRiichi/core/game"
	"fmt"
)

func main() {
	fmt.Println(game.InitGameState().GenerateGraphs())
}
