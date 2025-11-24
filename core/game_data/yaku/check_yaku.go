package yaku

import (
	"slices"

	meldfinder "codeberg.org/ijnakashiar/LibreRiichi/core/game/meld_finder"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/hand"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/tile"
)

type YakuContext struct {
	IsSelfDrawn                bool
	HasCalledRiichi            bool
	IsIppatsu                  bool
	IsLastTileDrawnOrDiscarded bool
	IsDeadWallCall             bool
	IsFromOpponentKanCall      bool
	IsDoubleRiichi             bool
	IsTenhou                   bool
	IsChiihou                  bool
	IsHandOpen                 bool
	HandInRiichi               bool
}

// Checks whether the current yaku is valid for the given hand
type YakuChecker func(
	hand []Tile,
	context YakuContext,
	winningTile Tile,
) bool

var YakuCheckerMap = map[YakuType]YakuChecker{
	MENZEN_TSUMO_YAKU:                 CheckMenzenTsumoYaku,
	RIICHI_YAKU:                       CheckRiichi,
	IPPATSU_YAKU:                      CheckIppatsu,
	PINFU_YAKU:                        CheckPinfu,
	IIPEIKOU_YAKU:                     CheckIipeikou,
	HAITEI_YAOYUE_YAKU:                CheckHaiteiYaoyue,
	HOUTEI_RAOYUI_YAKU:                CheckHouteiRaoyui,
	RINSHAN_KAIHOU_YAKU:               CheckRinshanKaihou,
	CHANKAN_YAKU:                      CheckChankan,
	TANYAO_YAKU:                       CheckTanyao,
	YAKUHAI_YAKU:                      CheckYakuhai,
	DOUBLE_RIICHI_YAKU:                CheckDoubleRiichi,
	CHANTAIYAO_YAKU:                   CheckChantaiyao,
	SANSHOKU_DOUJUN_YAKU:              CheckSanshokuDoujun,
	ITTSU_YAKU:                        CheckIttsu,
	TOITOI_YAKU:                       CheckToitoi,
	SANANKOU_YAKU:                     CheckSanankou,
	SANSHOKU_DOUKOU_YAKU:              CheckSanshokuDoukou,
	SANKANTSU_YAKU:                    CheckSankantsu,
	CHIITOITSU_YAKU:                   CheckChiitoitsu,
	HONROUTOU_YAKU:                    CheckHonroutou,
	SHOUSANGEN_YAKU:                   CheckShousangen,
	HONITSU_YAKU:                      CheckHonitsu,
	JUNCHAN_YAKU:                      CheckJunchan,
	RYANPEIKOU_YAKU:                   CheckRyanpeikou,
	CHINITSU_YAKU:                     CheckChinitsu,
	KAZOE_YAKUMAN_YAKU:                CheckKazoeYakuman,
	KOKUSHI_MUSOU_YAKU:                CheckKokushiMusou,
	KOKUSHI_MUSOU_THIRTEEN_WAITS_YAKU: CheckKokushiMusouThirteenWaits,
	SUUANKOU_YAKU:                     CheckSuuankou,
	DAISANGEN_YAKU:                    CheckDaisangen,
	SHOUSUUSHII_YAKU:                  CheckShousuushii,
	DAISUUSHII_YAKU:                   CheckDaisuushii,
	TSUUIISOU_YAKU:                    CheckTsuuiisou,
	CHINROUTOU_YAKU:                   CheckChinroutou,
	RYUUIISOU_YAKU:                    CheckRyuuiisou,
	CHUUREN_POUTOU_YAKU:               CheckChuurenPoutou,
	SUUKANTSU_YAKU:                    CheckSuukantsu,
	TENHOU_YAKU:                       CheckTenhou,
	CHIIHOU_YAKU:                      CheckChiihou,
	NAGASHI_MANGAN_YAKU:               CheckNagashiMangan,
}

func CheckHandCanWin(hand Hand, winningTile Tile, yakuContext YakuContext) bool {
	// Check Kokushi Musou
	if !yakuContext.IsHandOpen {
		typeList := []Tile{Manzu, Pinzu, Souzu}
		numberList := []uint8{1, 9}
		tileList := []Tile{}
		for _, t := range typeList {
			for _, n := range numberList {
				tileList = append(tileList, MakeNumberTile(t, n))
			}
		}

		// Thirteen waits
		if hand.ClosedHand.HasTile(tileList...) && slices.Contains(tileList, winningTile) {
			return true
		}

		hand.ClosedHand.Add(winningTile)
		totalValid := 0
		for _, tile := range tileList {
			count := hand.ClosedHand.HasTileN(tile)
			if count == 1 || count == 2 {
				totalValid += count
			}
		}
		// Regular kokushi
		if totalValid == 14 {
			return true
		}
		hand.ClosedHand.Pop(1)
	}

	// Check Chiitoitsu
	if !yakuContext.IsHandOpen {
		unique, count := hand.ClosedHand.UniqueTiles()
		success := true
		for i := range count {
			if count[i] != 2 {
				success = false
			}
		}
		if len(unique) == 7 && success {
			return true
		}
	}

	// Check four melds and a pair wins
	return meldfinder.HasWinningCombination(hand.ClosedHand.GetHand(), int(hand.OpenMeldCount()))
}

func CheckAllYaku(hand Hand, winningTile Tile, yakuContext YakuContext) map[YakuType]bool {
	return nil
}

// Menzen Tsumo Yaku requires:
// 1. The hand must be self-drawn
// 2. The hand must be closed/concealed (no open calls)
func CheckMenzenTsumoYaku(hand []Tile,
	context YakuContext,
	winningTile Tile,
) bool {

	if !context.IsSelfDrawn {
		return false
	}

	if context.IsHandOpen {
		return false
	}

	return true
}

func CheckRiichi(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckIppatsu(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckPinfu(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckIipeikou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckHaiteiYaoyue(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckHouteiRaoyui(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckRinshanKaihou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChankan(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckTanyao(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckYakuhai(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckDoubleRiichi(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChantaiyao(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckSanshokuDoujun(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckIttsu(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckToitoi(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckSanankou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckSanshokuDoukou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckSankantsu(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChiitoitsu(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckHonroutou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckShousangen(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckHonitsu(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckJunchan(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckRyanpeikou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChinitsu(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckKazoeYakuman(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckKokushiMusou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckKokushiMusouThirteenWaits(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckSuuankou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckDaisangen(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckShousuushii(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckDaisuushii(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckTsuuiisou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChinroutou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckRyuuiisou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChuurenPoutou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckSuukantsu(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckTenhou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChiihou(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckNagashiMangan(hand []Tile, context YakuContext, winningTile Tile) bool {
	return false
}
