package main

import (
	"fmt"

	"codeberg.org/ijnakashiar/LibreRiichi/core/game"
	"github.com/looplab/fsm"
)

func main() {
	fmt.Println(fsm.Visualize(game.InitRoundState().RoundFSM))
}
