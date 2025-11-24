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
