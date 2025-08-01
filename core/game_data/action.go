package core

import (
	"encoding/json"
	"fmt"
)

type ActionType uint8

//go:generate go run generate_actions.go -- Action
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

type Action interface {
	Data() any
}

type ActionHandler[T any, E any] interface {
	HandleRon(Ron, E) (T, error)
	HandleTsumo(Tsumo, E) (T, error)
	HandleRiichi(Riichi, E) (T, error)
	HandleToss(Toss, E) (T, error)
	HandleSkip(Skip, E) (T, error)
	HandlePon(Pon, E) (T, error)
	HandleKan(Kan, E) (T, error)
	HandleChii(Chii, E) (T, error)
	HandleDraw(Draw, E) (T, error)
}

type Ron struct {
	TileToRon Tile      `json:"tile_to_ron"`
	WinResult WinResult `json:"win_result"` // TODO: Remove
}

type Tsumo struct {
	TileToTsumo Tile `json:"tile_to_tsumo"`
}

type Riichi struct {
	TileToRiichi Tile `json:"tile_to_riichi"`
}

type Toss struct {
	TileToToss Tile `json:"tile_to_toss"`
}

type Skip struct {
	ActionToSkip Action `json:"action_to_skip"`
}

type Pon struct {
	TileToPon Tile `json:"tile_to_pon"`
}

type Kan struct {
	TileToKan Tile `json:"tile_to_kan"`
}

type Chii struct {
	TileToChii  Tile    `json:"tile_to_chii"`
	TilesInHand [2]Tile `json:"tiles_in_hand"`
}

type Draw struct {
	DrawnTile Tile `json:"drawn_tile"`
}

func (msg *Action) MarshalJSON() ([]byte, error) {
	var raw struct {
		ActionType ActionType `json:"action_type"`
		Data       ActionData `json:"data"`
	}

	switch msg.ActionData.(type) {
	case Chii:
		raw.ActionType = CHII
	case Draw:
		raw.ActionType = DRAW
	case Kan:
		raw.ActionType = KAN
	case Pon:
		raw.ActionType = PON
	case Riichi:
		raw.ActionType = RIICHI
	case Ron:
		raw.ActionType = RON
	case Skip:
		raw.ActionType = SKIP
	case Toss:
		raw.ActionType = TOSS
	case Tsumo:
		raw.ActionType = TSUMO
	default:
		panic(fmt.Sprintf("unexpected core.ActionData: %#v", msg.ActionData))
	}
	raw.Data = msg.ActionData

	return json.Marshal(raw)
}

func (msg *Action) UnmarshalJSON(rawData []byte) error {
	var raw struct {
		ActionType ActionType      `json:"action_type"`
		Data       json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	switch raw.ActionType {
	case CHII:
		message := Chii{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.ActionData = message
	case DRAW:
		message := Draw{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.ActionData = message
	case KAN:
		message := Kan{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.ActionData = message
	case PON:
		message := Pon{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.ActionData = message
	case RIICHI:
		message := Riichi{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.ActionData = message
	case RON:
		message := Ron{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.ActionData = message
	case SKIP:
		message := Skip{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.ActionData = message
	case TOSS:
		message := Toss{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.ActionData = message
	case TSUMO:
		message := Tsumo{}
		err := json.Unmarshal(raw.Data, &message)
		if err != nil {
			return err
		}
		msg.ActionData = message
	default:
		return fmt.Errorf("unexpected core.ActionType: %#v", raw.ActionType)
	}
	return nil
}

func ActionDecode[T any, E any](handler ActionHandler[T, E], data ActionData, extraData E) (ret T, err error) {
	switch v := data.(type) {
	case Chii:
		return handler.HandleChii(v, extraData)
	case Draw:
		return handler.HandleDraw(v, extraData)
	case Kan:
		return handler.HandleKan(v, extraData)
	case Pon:
		return handler.HandlePon(v, extraData)
	case Riichi:
		return handler.HandleRiichi(v, extraData)
	case Ron:
		return handler.HandleRon(v, extraData)
	case Skip:
		return handler.HandleSkip(v, extraData)
	case Toss:
		return handler.HandleToss(v, extraData)
	case Tsumo:
		return handler.HandleTsumo(v, extraData)
	default:
		return ret, fmt.Errorf("unexpected core.ActionType: %#v", data)
	}
}
