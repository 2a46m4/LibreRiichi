package game

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type yakuContext struct {
	isSelfDrawn                bool
	hasCalledRiichi            bool
	isIppatsu                  bool
	isLastTileDrawnOrDiscarded bool
	isDeadWallCall             bool
	isFromOpponentKanCall      bool
	isDoubleRiichi             bool
	isTenhou                   bool
	isChiihou                  bool
}

type YakuChecker func(
	playerIdx uint8,
	data *MahjongRoundData,
	context yakuContext,
	extraTile ...Tile,
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

func CheckYaku(playerIdx uint8, data *MahjongRoundData, extraTile ...Tile) map[YakuType]bool {

	return nil
}

func CheckMenzenTsumoYaku(playerIdx uint8,
	data *MahjongRoundData,
	context yakuContext,
	extraTile ...Tile) bool {

	// Menzen Tsumo Yaku requires:
	// 1. The hand must be self-drawn
	// 2. The hand must be closed/concealed (no open calls)

	// Check if the hand is self-drawn
	if !context.isSelfDrawn {
		return false
	}

	// Check if the hand is closed (no open calls: chi, pon, kan from opponents)
	// The Open() method returns true when there are no open calls
	if !data.tileState.Hands[playerIdx].Open() {
		return false
	}

	return true
}

func CheckRiichi(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return data.tileState.Hands[playerIdx].InRiichi
}

func CheckIppatsu(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckPinfu(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckIipeikou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckHaiteiYaoyue(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckHouteiRaoyui(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckRinshanKaihou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckChankan(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckTanyao(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckYakuhai(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckDoubleRiichi(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckChantaiyao(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckSanshokuDoujun(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckIttsu(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckToitoi(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckSanankou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckSanshokuDoukou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckSankantsu(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckChiitoitsu(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckHonroutou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckShousangen(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckHonitsu(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckJunchan(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckRyanpeikou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckChinitsu(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckKazoeYakuman(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckKokushiMusou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckKokushiMusouThirteenWaits(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckSuuankou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckDaisangen(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckShousuushii(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckDaisuushii(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckTsuuiisou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckChinroutou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckRyuuiisou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckChuurenPoutou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckSuukantsu(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckTenhou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckChiihou(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}

func CheckNagashiMangan(playerIdx uint8, data *MahjongRoundData, context yakuContext, extraTile ...Tile) bool {
	return false
}
