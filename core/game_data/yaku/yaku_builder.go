package yaku

import (
	"errors"
)

type YakuBuilder struct {
	value    YakuType
	HandOpen bool
	isKazoe  bool // Always checked when adding han
}

func NewYakuBuilder(isOpen bool) YakuBuilder {
	return YakuBuilder{
		value:    0,
		HandOpen: false,
		isKazoe:  false,
	}
}

func (yakuBuilder YakuBuilder) Build() (list YakuList, err error) {
	list.IsOpen = yakuBuilder.HandOpen

	if yakuBuilder.value == NO_YAKU {
		return list, errors.New("No yaku")
	}

	// Check if there is yakuman which overrides other yaku
	if isNonKazoeYakuman(yakuBuilder.value) {
		for yaku := range iterateYaku(yakuBuilder.value) {
			if isNonKazoeYakuman(yaku) {
				list.Yakus = append(list.Yakus, Yaku{
					HanValue: hanValue(yaku, yakuBuilder.HandOpen),
					YakuName: yaku.String(),
				})
			}
		}

		return list, nil
	}

	if yakuBuilder.isKazoe {
		list.IsKazoe = true
	}

	for yaku := range iterateYaku(yakuBuilder.value) {
		list.Yakus = append(list.Yakus, Yaku{
			HanValue: hanValue(yaku, yakuBuilder.HandOpen),
			YakuName: yaku.String(),
		})
	}

	return list, nil
}

func (yaku *YakuBuilder) AddYaku(added YakuType) error {
	if added == NO_YAKU {
		return errors.New("Can't add no yaku")
	}

	if added == KAZOE_YAKUMAN_YAKU {
		return errors.New("Can't add Kazoe Yakuman directly")
	}

	if !added.canBeOpen() && yaku.HandOpen {
		return errors.New("Can't use this yaku in an open hand")
	}

	yaku.value = yaku.value&(^NO_YAKU) | added
	yaku.isKazoe = yaku.han() >= 13

	return nil
}

func (yaku *YakuBuilder) Reset() {
	yaku.isKazoe = false
	yaku.value = NO_YAKU
}

// Without counting kazoe yakuman
func (yakuBuilder YakuBuilder) han() (totalHan int) {
	if yakuBuilder.value == NO_YAKU {
		return 0
	}

	for yaku := range iterateYaku(yakuBuilder.value) {
		totalHan += hanValue(yaku, yakuBuilder.HandOpen)
	}
	return totalHan
}
