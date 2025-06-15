package core

import (
	"encoding/json"
	"fmt"
)

type Visibility uint8

const (
	PLAYER Visibility = 0

	GLOBAL = 255
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

type ActionHandler interface {
	HandleRon(RonData) error
	HandleTsumo(TsumoData) error
	HandleRiichi(RiichiData) error
	HandleToss(TossData) error
	HandleSkip(SkipData) error
	HandlePon(PonData) error
	HandleKan(KanData) error
	HandleChii(ChiiData) error
	HandleDraw(DrawData) error
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

func ActionDecodeAndDispatch(handler ActionHandler, rawData []byte) error {
	var raw struct {
		ActionType ActionType      `json:"action_type"`
		Data       json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	switch raw.ActionType {
	case CHII:
		message := ChiiData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleChii(message)
	case DRAW:
		message := DrawData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleDraw(message)
	case KAN:
		message := KanData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleKan(message)
	case PON:
		message := PonData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandlePon(message)
	case RIICHI:
		message := RiichiData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleRiichi(message)
	case RON:
		message := RonData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleRon(message)
	case SKIP:
		message := SkipData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleSkip(message)
	case TOSS:
		message := TossData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleToss(message)
	case TSUMO:
		message := TsumoData{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		return handler.HandleTsumo(message)
	default:
		return fmt.Errorf("unexpected core.ActionType: %#v", raw.ActionType)
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
