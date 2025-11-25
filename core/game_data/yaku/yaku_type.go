package yaku

import "codeberg.org/ijnakashiar/LibreRiichi/core/util"

//go:generate stringer -type=YakuType
type YakuType uint64

// Whether or not the yaku can be open.
func (yaku YakuType) canBeOpen() bool {
	return core.All(iterateYaku(yaku), func(y YakuType) bool {
		return canBeOpen[y]
	})
}

func isNonKazoeYakuman(yaku YakuType) bool {
	return yaku&yakumanBitset != 0
}

func iterateYaku(yakus YakuType) func(func(YakuType) bool) {
	return func(yield func(YakuType) bool) {
		for itr := 0; itr != 64; itr += 1 {
			if (yakus>>itr)&1 == 0 {
				continue
			}

			if !yield(yakus & (1 << itr)) {
				return
			}
		}
	}
}
