package yaku

type YakuList struct {
	Yakus   []Yaku
	IsOpen  bool
	IsKazoe bool
}

type Yaku struct {
	HanValue int
	YakuName string
}

func (yaku YakuList) Han() int {
	if yaku.IsKazoe {
		return 13
	} else {
		res := 0
		for _, yaku := range yaku.Yakus {
			res += yaku.HanValue
		}
		return res
	}
}
