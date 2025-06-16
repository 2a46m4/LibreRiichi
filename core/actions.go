package core

import (
	"encoding/json"
	"fmt"
)

type Visibility uint8

const (
	PLAYER Visibility = iota
	PARTIAL
	GLOBAL
)

type ActionType uint8

const (
	RON ActionType = iota
	TSUMO
	RIICHI
	TOSS
	SKIP
	PON
	KAN
	CHII
	DRAW
)

type ActionData struct {
	ActionType ActionType `json:"action_type"`
	Data       any        `json:"data"`
}

type ActionHandler[T any] interface {
	HandleRon(RonData, uint8) (T, error)
	HandleTsumo(TsumoData, uint8) (T, error)
	HandleRiichi(RiichiData, uint8) (T, error)
	HandleToss(TossData, uint8) (T, error)
	HandleSkip(SkipData, uint8) (T, error)
	HandlePon(PonData, uint8) (T, error)
	HandleKan(KanData, uint8) (T, error)
	HandleChii(ChiiData, uint8) (T, error)
	HandleDraw(DrawData, uint8) (T, error)
}

type RonData struct {
	TileToRon Tile      `json:"tile_to_ron"`
	WinResult WinResult `json:"win_result"`
}

type TsumoData struct {
	TileToTsumo Tile `json:"tile_to_tsumo"`
}

type RiichiData struct {
	TileToRiichi Tile `json:"tile_to_riichi"`
}

type TossData struct {
	TileToToss Tile `json:"tile_to_toss"`
}

type SkipData struct {
	ActionToSkip PlayerActionData `json:"action_to_skip"`
}

type PonData struct {
	TileToPon Tile `json:"tile_to_pon"`
}

type KanData struct {
	TileToKan Tile `json:"tile_to_kan"`
}

type ChiiData struct {
	TileToChii  Tile    `json:"tile_to_chii"`
	TilesInHand [2]Tile `json:"tiles_in_hand"`
}

type DrawData struct {
	DrawnTile Tile `json:"drawn_tile"`
}

func (msg *ActionData) UnmarshalJSON(rawData []byte) error {
	var raw struct {
		ActionType ActionType      `json:"action_type"`
		Data       json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	msg.ActionType = raw.ActionType

	switch msg.ActionType {
	case CHII:
		message := ChiiData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.Data = message
	case DRAW:
		message := DrawData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.Data = message
	case KAN:
		message := KanData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.Data = message
	case PON:
		message := PonData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.Data = message
	case RIICHI:
		message := RiichiData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.Data = message
	case RON:
		message := RonData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.Data = message
	case SKIP:
		message := SkipData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.Data = message
	case TOSS:
		message := TossData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.Data = message
	case TSUMO:
		message := TsumoData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.Data = message
	default:
		return fmt.Errorf("unexpected core.ActionType: %#v", raw.ActionType)
	}
	return nil
}

func ActionDecode[T any](handler ActionHandler[T], data ActionData, fromPlayer uint8) (ret T, err error) {
	switch data.ActionType {
	case CHII:
		message, ok := data.Data.(ChiiData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleChii(message, fromPlayer)
	case DRAW:
		message, ok := data.Data.(DrawData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleDraw(message, fromPlayer)
	case KAN:
		message, ok := data.Data.(KanData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleKan(message, fromPlayer)
	case PON:
		message, ok := data.Data.(PonData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandlePon(message, fromPlayer)
	case RIICHI:
		message, ok := data.Data.(RiichiData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleRiichi(message, fromPlayer)
	case RON:
		message, ok := data.Data.(RonData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleRon(message, fromPlayer)
	case SKIP:
		message, ok := data.Data.(SkipData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleSkip(message, fromPlayer)
	case TOSS:
		message, ok := data.Data.(TossData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleToss(message, fromPlayer)
	case TSUMO:
		message, ok := data.Data.(TsumoData)
		if !ok {
			return ret, BadTypeError{}
		}
		return handler.HandleTsumo(message, fromPlayer)
	default:
		return ret, fmt.Errorf("unexpected core.ActionType: %#v", data.ActionType)
	}
}

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
		panic(fmt.Sprintf("unexpected core.SetupType: %#v", msg.Type))
	}
	return nil
}

func SetupDecode[T any](handler ActionHandler[T], data ActionData) error {
	panic("NYI")
}
