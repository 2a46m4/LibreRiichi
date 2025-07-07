package core

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/errors"
	"encoding/json"
	"fmt"
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

type ActionHandler[T any, E any] interface {
	HandleRon(RonData, E) (T, error)
	HandleTsumo(TsumoData, E) (T, error)
	HandleRiichi(RiichiData, E) (T, error)
	HandleToss(TossData, E) (T, error)
	HandleSkip(SkipData, E) (T, error)
	HandlePon(PonData, E) (T, error)
	HandleKan(KanData, E) (T, error)
	HandleChii(ChiiData, E) (T, error)
	HandleDraw(DrawData, E) (T, error)
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
	ActionToSkip ActionData `json:"action_to_skip"`
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

func ActionDecode[T any, E any](handler ActionHandler[T, E], data ActionData, extraData E) (ret T, err error) {
	switch data.ActionType {
	case CHII:
		message, ok := data.Data.(ChiiData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleChii(message, extraData)
	case DRAW:
		message, ok := data.Data.(DrawData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleDraw(message, extraData)
	case KAN:
		message, ok := data.Data.(KanData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleKan(message, extraData)
	case PON:
		message, ok := data.Data.(PonData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandlePon(message, extraData)
	case RIICHI:
		message, ok := data.Data.(RiichiData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleRiichi(message, extraData)
	case RON:
		message, ok := data.Data.(RonData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleRon(message, extraData)
	case SKIP:
		message, ok := data.Data.(SkipData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleSkip(message, extraData)
	case TOSS:
		message, ok := data.Data.(TossData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleToss(message, extraData)
	case TSUMO:
		message, ok := data.Data.(TsumoData)
		if !ok {
			return ret, BadMessage{}
		}
		return handler.HandleTsumo(message, extraData)
	default:
		return ret, fmt.Errorf("unexpected core.ActionType: %#v", data.ActionType)
	}
}
