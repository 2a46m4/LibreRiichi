package core

import (
	"encoding/json"
	"fmt"
)

type SetupType uint8

const (
	INITIAL_TILES SetupType = iota
	DORA
	STARTING_POINTS
	PLAYER_NUMBER
	PLAYER_ORDER
	ROUND_WIND
	ROUND_NUMBER
)

type Setup struct {
	Type SetupType `json:"setup_type"`
	Data any       `json:"data"`
}

func (msg *Setup) UnmarshalJSON(rawData []byte) error {
	var raw struct {
		SetupType SetupType       `json:"setup_type"`
		Data      json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	msg.Type = raw.SetupType

	switch msg.Type {
	case DORA:
	case INITIAL_TILES:
	case PLAYER_NUMBER:
	case PLAYER_ORDER:
	case ROUND_NUMBER:
	case ROUND_WIND:
	case STARTING_POINTS:
	default:
		return fmt.Errorf("unexpected core.SetupType: %#v", msg.Type)
	}
	return nil
}

func SetupDecode[T any, E any](handler ActionHandler[T, E], data ActionData) error {
	panic("NYI")
}
