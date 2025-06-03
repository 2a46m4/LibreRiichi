package core

type YakuType uint64

const (
	NO_YAKU YakuType = 1 << iota
	MENZEN_TSUMO_YAKU
	RIICHI_YAKU
	IPPATSU_YAKU
	PINFU_YAKU
	IIPEIKOU_YAKU

	HAITEI_YAOYUE_YAKU
	HOUTEI_RAOYUI_YAKU
	RINSHAN_KAIHOU_YAKU
	CHANKAN_YAKU
	TANYAO_YAKU
	YAKUHAI_YAKU

	DOUBLE_RIICHI_YAKU
	CHANTAIYAO_YAKU
	SANSHOKU_DOUJUN_YAKU
	ITTSU_YAKU
	TOITOI_YAKU
	SANANKOU_YAKU
	SANSHOKU_DOUKOU_YAKU
	SANKANTSU_YAKU
	CHIITOITSU_YAKU
	HONROUTOU_YAKU
	SHOUSANGEN_YAKU

	HONITSU_YAKU
	JUNCHAN_YAKU
	RYANPEIKOU_YAKU

	CHINITSU_YAKU

	KAZOE_YAKUMAN_YAKU
	KOKUSHI_MUSOU_YAKU
	KOKUSHI_MUSOU_THIRTEEN_WAITS_YAKU
	SUUANKOU_YAKU
	DAISANGEN_YAKU
	SHOUSUUSHII_YAKU
	DAISUUSHII_YAKU
	TSUUIISOU_YAKU
	CHINROUTOU_YAKU
	RYUUIISOU_YAKU
	CHUUREN_POUTOU_YAKU
	SUUKANTSU_YAKU

	TENHOU_YAKU
	CHIIHOU_YAKU

	NAGASHI_MANGAN_YAKU
)

func (yaku YakuType) Han() int {

	return []int{
		NO_YAKU:      0,
		MENZEN_TSUMO_YAKU: 1,
		RIICHI_YAKU:       1,
		IPPATSU_YAKU:      1,
		PINFU_YAKU:        1,
		IIPEIKOU_YAKU:     1,

		HAITEI_YAOYUE_YAKU:  1,
		HOUTEI_RAOYUI_YAKU:  1,
		RINSHAN_KAIHOU_YAKU: 1,
		CHANKAN_YAKU:        1,
		TANYAO_YAKU:         1,
		YAKUHAI_YAKU:        1,

		DOUBLE_RIICHI_YAKU:   2,
		CHANTAIYAO_YAKU:      2,
		SANSHOKU_DOUJUN_YAKU: 2,
		ITTSU_YAKU:           2,
		TOITOI_YAKU:          2,
		SANANKOU_YAKU:        2,
		SANSHOKU_DOUKOU_YAKU: 2,
		SANKANTSU_YAKU:       2,
		CHIITOITSU_YAKU:      2,
		HONROUTOU_YAKU:       2,
		SHOUSANGEN_YAKU:      2,

		HONITSU_YAKU:    3,
		JUNCHAN_YAKU:    3,
		RYANPEIKOU_YAKU: 3,

		CHINITSU_YAKU: 6,

		KAZOE_YAKUMAN_YAKU:                13,
		KOKUSHI_MUSOU_YAKU:                13,
		KOKUSHI_MUSOU_THIRTEEN_WAITS_YAKU: 26,
		SUUANKOU_YAKU:                     13,
		DAISANGEN_YAKU:                    13,
		SHOUSUUSHII_YAKU:                  13,
		DAISUUSHII_YAKU:                   26,
		TSUUIISOU_YAKU:                    13,
		CHINROUTOU_YAKU:                   13,
		RYUUIISOU_YAKU:                    13,
		CHUUREN_POUTOU_YAKU:               13,
		SUUKANTSU_YAKU:                    13,

		TENHOU_YAKU:  13,
		CHIIHOU_YAKU: 13,

		NAGASHI_MANGAN_YAKU: 3,
	}[yaku]
}

func (yaku YakuType) IsYakuman() bool {
	yakuman := []YakuType{
		KAZOE_YAKUMAN_YAKU,
		KOKUSHI_MUSOU_YAKU,
		SUUANKOU_YAKU,
		DAISANGEN_YAKU,
		SHOUSUUSHII_YAKU,
		DAISUUSHII_YAKU,
		TSUUIISOU_YAKU,
		CHINROUTOU_YAKU,
		RYUUIISOU_YAKU,
		CHUUREN_POUTOU_YAKU,
		SUUKANTSU_YAKU,
		TENHOU_YAKU,
		CHIIHOU_YAKU,
	}

	yakuman_bits := NO_YAKU
	for _, ym := range yakuman {
		yakuman_bits |= ym
	}

	if (yaku & yakuman_bits) != 0 {
		return true
	}

	hanCount := 0
	for i := 0; i < 64; i++ {
		hanCount += YakuType(1 << i).Han()
	}

	if hanCount == 13 {
		return true
	}

	return false
}

func (yaku YakuType) OpenHand() bool {
	return map[YakuType]bool{
		NO_YAKU:           true,
		MENZEN_TSUMO_YAKU: false,
		RIICHI_YAKU:       false,
		IPPATSU_YAKU:      false,
		PINFU_YAKU:        false,
		IIPEIKOU_YAKU:     false,

		HAITEI_YAOYUE_YAKU:  true,
		HOUTEI_RAOYUI_YAKU:  true,
		RINSHAN_KAIHOU_YAKU: true,
		CHANKAN_YAKU:        true,
		TANYAO_YAKU:         true,
		YAKUHAI_YAKU:        true,

		DOUBLE_RIICHI_YAKU:   false,
		CHANTAIYAO_YAKU:      true,
		SANSHOKU_DOUJUN_YAKU: true,
		ITTSU_YAKU:           true,
		TOITOI_YAKU:          true,
		SANANKOU_YAKU:        true,
		SANSHOKU_DOUKOU_YAKU: true,
		SANKANTSU_YAKU:       true,
		CHIITOITSU_YAKU:      false,
		HONROUTOU_YAKU:       true,
		SHOUSANGEN_YAKU:      true,

		HONITSU_YAKU:    true,
		JUNCHAN_YAKU:    true,
		RYANPEIKOU_YAKU: false,

		CHINITSU_YAKU: true,

		KAZOE_YAKUMAN_YAKU:  true,
		KOKUSHI_MUSOU_YAKU:  false,
		SUUANKOU_YAKU:       false,
		DAISANGEN_YAKU:      true,
		SHOUSUUSHII_YAKU:    true,
		DAISUUSHII_YAKU:     true,
		TSUUIISOU_YAKU:      true,
		CHINROUTOU_YAKU:     true,
		RYUUIISOU_YAKU:      true,
		CHUUREN_POUTOU_YAKU: false,
		SUUKANTSU_YAKU:      true,

		TENHOU_YAKU:  false,
		CHIIHOU_YAKU: false,

		NAGASHI_MANGAN_YAKU: false,
	}[yaku]
}

func (yaku YakuType) HanLossOnOpen() int {
	return map[YakuType]int{
		NO_YAKU: 0,

		HAITEI_YAOYUE_YAKU:  0,
		HOUTEI_RAOYUI_YAKU:  0,
		RINSHAN_KAIHOU_YAKU: 0,
		CHANKAN_YAKU:        0,
		TANYAO_YAKU:         0,
		YAKUHAI_YAKU:        0,

		CHANTAIYAO_YAKU:      1,
		SANSHOKU_DOUJUN_YAKU: 1,
		ITTSU_YAKU:           1,
		TOITOI_YAKU:          0,
		SANANKOU_YAKU:        0,
		SANSHOKU_DOUKOU_YAKU: 0,
		SANKANTSU_YAKU:       0,
		HONROUTOU_YAKU:       0,
		SHOUSANGEN_YAKU:      0,

		HONITSU_YAKU: 1,
		JUNCHAN_YAKU: 1,

		CHINITSU_YAKU: 1,

		KAZOE_YAKUMAN_YAKU: 0,
		DAISANGEN_YAKU:     0,
		SHOUSUUSHII_YAKU:   0,
		DAISUUSHII_YAKU:    0,
		TSUUIISOU_YAKU:     0,
		CHINROUTOU_YAKU:    0,
		RYUUIISOU_YAKU:     0,
		SUUKANTSU_YAKU:     0,
	}[yaku]
}

func (yaku *YakuType) Set() {

}
