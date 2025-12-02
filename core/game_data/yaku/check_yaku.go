package yaku

import (
	"slices"

	meldfinder "codeberg.org/ijnakashiar/LibreRiichi/core/game/meld_finder"
	core "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data/hand"
	"codeberg.org/ijnakashiar/LibreRiichi/core/game_data/score"
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
	HandInRiichi               bool
	RoundWind                  core.Wind
	PlayerWind                 core.Wind
	IsDealer                   bool
}

// Checks whether the current yaku is valid for the given hand
type YakuChecker func(
	hand *Hand,
	context YakuContext,
	winningTile Tile,
) bool

// Checks just the conventional 4 melds + 1 pair win
//
// No kazoe yakuman, nor chiitoitsu and kokushi musou
var YakuCheckerMap = map[YakuType]YakuChecker{
	MENZEN_TSUMO_YAKU:    CheckMenzenTsumoYaku,
	RIICHI_YAKU:          CheckRiichi,
	IPPATSU_YAKU:         CheckIppatsu,
	PINFU_YAKU:           CheckPinfu,
	IIPEIKOU_YAKU:        CheckIipeikou,
	HAITEI_YAOYUE_YAKU:   CheckHaiteiYaoyue,
	HOUTEI_RAOYUI_YAKU:   CheckHouteiRaoyui,
	RINSHAN_KAIHOU_YAKU:  CheckRinshanKaihou,
	CHANKAN_YAKU:         CheckChankan,
	TANYAO_YAKU:          CheckTanyao,
	YAKUHAI_YAKU:         CheckYakuhai,
	DOUBLE_RIICHI_YAKU:   CheckDoubleRiichi,
	CHANTAIYAO_YAKU:      CheckChantaiyao,
	SANSHOKU_DOUJUN_YAKU: CheckSanshokuDoujun,
	ITTSU_YAKU:           CheckIttsu,
	TOITOI_YAKU:          CheckToitoi,
	SANANKOU_YAKU:        CheckSanankou,
	SANSHOKU_DOUKOU_YAKU: CheckSanshokuDoukou,
	SANKANTSU_YAKU:       CheckSankantsu,
	HONROUTOU_YAKU:       CheckHonroutou,
	SHOUSANGEN_YAKU:      CheckShousangen,
	HONITSU_YAKU:         CheckHonitsu,
	JUNCHAN_YAKU:         CheckJunchan,
	RYANPEIKOU_YAKU:      CheckRyanpeikou,
	CHINITSU_YAKU:        CheckChinitsu,
	SUUANKOU_YAKU:        CheckSuuankou,
	DAISANGEN_YAKU:       CheckDaisangen,
	SHOUSUUSHII_YAKU:     CheckShousuushii,
	DAISUUSHII_YAKU:      CheckDaisuushii,
	TSUUIISOU_YAKU:       CheckTsuuiisou,
	CHINROUTOU_YAKU:      CheckChinroutou,
	RYUUIISOU_YAKU:       CheckRyuuiisou,
	CHUUREN_POUTOU_YAKU:  CheckChuurenPoutou,
	SUUKANTSU_YAKU:       CheckSuukantsu,
	TENHOU_YAKU:          CheckTenhou,
	CHIIHOU_YAKU:         CheckChiihou,
}

type NoValidYakus struct {}
func (NoValidYakus) Error() string {
    return "No valid yakus"
}

// Returns yaku and score of a given hand
func CheckYakuAndScore(hand *Hand, yakuContext YakuContext, winningTile Tile) (YakuList, score.PointValue, error) {
	yakuBuilder := NewYakuBuilder(hand.Closed())

	if CheckKokushiMusouThirteenWaits(hand, yakuContext, winningTile) {
		yakuBuilder.AddYaku(KOKUSHI_MUSOU_THIRTEEN_WAITS_YAKU)
		list, err := yakuBuilder.Build()
		if err != nil {
			return YakuList{}, score.PointValue{}, err
		}
		currentScore := score.ComputePoints(score.Points{
			Han: uint16(list.Han()),
			Fu:  0,
		}, yakuContext.IsDealer)
		return list, currentScore, nil
	}

	if CheckKokushiMusou(hand, yakuContext, winningTile) {
		yakuBuilder.AddYaku(KOKUSHI_MUSOU_YAKU)
		list, err := yakuBuilder.Build()
		if err != nil {
			return YakuList{}, score.PointValue{}, err
		}
		currentScore := score.ComputePoints(score.Points{
			Han: uint16(list.Han()),
			Fu:  0,
		}, yakuContext.IsDealer)
		return list, currentScore, err
	}

	if CheckChiitoitsu(hand, yakuContext, winningTile) {
		yakuBuilder.AddYaku(CHIITOITSU_YAKU)
		list, err := yakuBuilder.Build()
		if err != nil {
			return YakuList{}, score.PointValue{}, err
		}
		currentScore := score.ComputePoints(score.Points{
			Han: uint16(list.Han()),
			Fu:  25, // Chiitoitsu always has 25 fu
		}, yakuContext.IsDealer)
		return list, currentScore, nil
	}

	// Check four melds and a pair wins
	combos := meldfinder.FindMelds(hand.ClosedHand.Hand(), int(hand.OpenMeldCount()))
	var maxScore uint = 0
	var scoreValue score.PointValue
	var yakuList YakuList
	foundValidYaku := false

RANGE_OVER_COMBOS:
	for _, combo := range combos {
		yakuBuilder.Reset()

		// Compute han
		for yakuType, check := range YakuCheckerMap {
			if check(hand, yakuContext, winningTile) {
				err := yakuBuilder.AddYaku(yakuType)
				if err != nil {
					goto RANGE_OVER_COMBOS
				}
			}
		}
		list, err := yakuBuilder.Build()
		if err != nil {
			goto RANGE_OVER_COMBOS
		}

		// TODO: Compute score
		fu := score.ComputeFu(combo, nil, score.WinByTsumo, core.East, core.East)
		currentScore := score.ComputePoints(score.Points{
			Han: uint16(list.Han()),
			Fu:  fu,
		}, yakuContext.IsDealer)

		if maxScore < currentScore.NonDealerTsumo {
			scoreValue = currentScore
			yakuList = list
			foundValidYaku = true
		}
	}

	if foundValidYaku {
		return yakuList, scoreValue, nil
	} else {
		return yakuList, scoreValue, NoValidYakus{}
	}
}

// Menzen Tsumo Yaku requires:
// 1. The hand must be self-drawn
// 2. The hand must be closed/concealed (no open calls)
func CheckMenzenTsumoYaku(hand *Hand,
	context YakuContext,
	winningTile Tile,
) bool {

	if !context.IsSelfDrawn {
		return false
	}

	if hand.Open() {
		return false
	}

	return true
}

func CheckRiichi(hand *Hand, context YakuContext, winningTile Tile) bool {
    if hand.Open() {
	return false
    }
    return false
}

func CheckIppatsu(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckPinfu(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckIipeikou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckHaiteiYaoyue(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckHouteiRaoyui(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckRinshanKaihou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChankan(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckTanyao(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckYakuhai(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckDoubleRiichi(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChantaiyao(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckSanshokuDoujun(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckIttsu(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckToitoi(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckSanankou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckSanshokuDoukou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckSankantsu(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChiitoitsu(hand *Hand, context YakuContext, winningTile Tile) bool {
	if hand.Open() {
		return false
	}

	unique, count := hand.ClosedHand.UniqueTiles()
	for i := range count {
		if count[i] != 2 {
			return false
		}
	}
	return len(unique) == 7
}

func CheckHonroutou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckShousangen(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckHonitsu(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckJunchan(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckRyanpeikou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChinitsu(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckKazoeYakuman(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckKokushiMusou(hand *Hand, context YakuContext, winningTile Tile) bool {
    if hand.Open() {
		return false
	}

	typeList := []Tile{Manzu, Pinzu, Souzu}
	numberList := []uint8{1, 9}
	tileList := []Tile{}
	for _, t := range typeList {
		for _, n := range numberList {
			tileList = append(tileList, MakeNumberTile(t, n))
		}
	}

	hand.ClosedHand.Add(winningTile)
	defer hand.ClosedHand.Pop(1)
	totalValid := 0
	for _, tile := range tileList {
		count := hand.ClosedHand.HasTileN(tile)
		// All the edge tiles must exist in the hand
		// But can't have more than 2
		if count == 1 || count == 2 {
			totalValid += count
		} else {
			return false
		}
	}

	// Make sure that there is a pair
	return totalValid == 14
}

func CheckKokushiMusouThirteenWaits(hand *Hand, context YakuContext, winningTile Tile) bool {
    if hand.Open() {
	return false
    }

    typeList := []Tile{Manzu, Pinzu, Souzu}
    numberList := []uint8{1, 9}
    tileList := []Tile{}
    for _, t := range typeList {
	for _, n := range numberList {
	    tileList = append(tileList, MakeNumberTile(t, n))
	}
    }

    return hand.ClosedHand.HasTile(tileList...) && slices.Contains(tileList, winningTile)
}

func CheckSuuankou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckDaisangen(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckShousuushii(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckDaisuushii(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckTsuuiisou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChinroutou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckRyuuiisou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChuurenPoutou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckSuukantsu(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckTenhou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

func CheckChiihou(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}

// TODO: Need to grab the entire board state
func CheckNagashiMangan(hand *Hand, context YakuContext, winningTile Tile) bool {
	return false
}
