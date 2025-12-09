package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game/meld_finder"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/hand"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
)

func CanAnkan(kan Kan, hand *Hand) bool {
    return !hand.InRiichi && hand.FullHand() && (hand.ClosedHand.HasTileN(kan.TileToKan) == 4)
}

func CanRon() {
    
}

// Requires a non full hand
func GetRiichiTargets(hand *Hand, extra Tile) (list []Riichi) {
    if !hand.Closed() {
	return nil
    }

    if hand.InRiichi {
	return nil
    }

    closed := hand.ClosedHand
    closed.Add(extra)
    for _, tile := range GetTileList() {
	handCopy := append([]Tile{}, closed.GetHand()...)
	for _, curTile := range handCopy {
	    // Impossible meld 
	    if curTile != tile && closed.HasTileN(tile) == 4 {
		continue
	    }

	    closed.RemoveTile(curTile)
	    closed.Add(tile)

	    if HasWinningCombination(closed.GetHand(), 0) {
		list = append(list, Riichi{
	    	    TileToRiichi: tile,
		})
	    }

	    closed.Pop(1)
	    closed.Add(curTile)	    
	}
    }

    return list
}


