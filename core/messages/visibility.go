package core

type Visibility uint8

const (
	PLAYER Visibility = iota
	PARTIAL
	EXCLUDE
	GLOBAL
)
