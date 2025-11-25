package score

import (
	MeldFinder "codeberg.org/ijnakashiar/LibreRiichi/core/game/meld_finder"
	game "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	tile "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
	"codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type WinMethodValue int

const (
	WinByClosedRon WinMethodValue = 10
	WinByTsumo     WinMethodValue = 2
)

// Represents the scoring system
type Points struct {
	Han uint16
	Fu  uint16
}

// Represents the amount of points that each role has to pay. In the
// case where the winner is the dealer and he wins by tsumo,
// DealerTsumo is invalid.
type PointValue struct {
	NonDealerTsumo uint
	DealerTsumo    uint
	Ron            uint
}

func ComputePoints(
	points Points,
	isDealer bool,
) (value PointValue) {

	roundUp := func(score uint) uint {
		return score + (100 - (score % 100))
	}

	// Tsumo, how much non/dealers pay
	var nonDealerMultiplier uint
	var dealerMultiplier uint
	// Ron, how much everyone pays
	var ronMultiplier uint

	if isDealer {
		nonDealerMultiplier = 2
		dealerMultiplier = 0 // Invalid since the winner is the dealer
		ronMultiplier = 6
	} else {
		nonDealerMultiplier = 1
		dealerMultiplier = 2
		ronMultiplier = 4
	}

	var multiplyBy float32
	var manganScore uint = 2000
	baseScore := uint(points.Fu) * (core.IntPow(uint(2), uint(2+points.Han)))
	if points.Han <= 4 && roundUp(baseScore*ronMultiplier) < manganScore*ronMultiplier {
		value.NonDealerTsumo = roundUp(baseScore * nonDealerMultiplier)
		value.DealerTsumo = roundUp(baseScore * dealerMultiplier)
		value.Ron = roundUp(baseScore * ronMultiplier)
		return value
	} else {
		baseScore = manganScore
	}

	switch points.Han {
	case 5:
		multiplyBy = 1
	case 6, 7:
		multiplyBy = 1.5
	case 8, 9, 10:
		multiplyBy = 2
	case 11, 12:
		multiplyBy = 3
	case 13:
		multiplyBy = 4
	default:
		multiplyBy = float32(points.Han / 13)
	}

	value.NonDealerTsumo = uint(float32(baseScore) * float32(nonDealerMultiplier) * multiplyBy)
	value.DealerTsumo = uint(float32(baseScore) * float32(dealerMultiplier) * multiplyBy)
	value.Ron = uint(float32(baseScore) * float32(ronMultiplier) * multiplyBy)

	return value
}

// TODO
func ComputeFu(
	winningCombination MeldFinder.WinningCombination,
	waitsAtTenpai []tile.Tile,
	winMethod WinMethodValue,
	playerWind game.Wind,
	roundWind game.Wind,
) uint16 {
	return 0
}
