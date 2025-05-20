package core

import (
	"encoding/json"
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

type PlayerAction struct {
	Action     ActionType `json:"action_type"`
	FromPlayer uint8      `json:"from_player"`
	Data       ActionData `json:"data"`
}

type ActionData interface{ actionDataImpl() }

type RonData struct {
	TileToRon Tile `json:"tile_to_ron"`
}

func (RonData) actionDataImpl() {}

type TsumoData struct {
	TileToTsumo Tile `json:"tile_to_tsumo"`
}

func (TsumoData) actionDataImpl() {}

type RiichiData struct {
	TileToRiichi Tile `json:"tile_to_riichi"`
}

func (RiichiData) actionDataImpl() {}

type TossData struct {
	TileToToss Tile `json:"tile_to_toss"`
}

func (TossData) actionDataImpl() {}

type SkipData struct {
	TileToSkip Tile `json:"tile_to_skip"`
}

func (SkipData) actionDataImpl() {}

type PonData struct {
	TileToPon Tile `json:"tile_to_pon"`
}

func (PonData) actionDataImpl() {}

type KanData struct {
	TileToKan Tile `json:"tile_to_kan"`
}

func (KanData) actionDataImpl() {}

type ChiiData struct {
	TileToChii  Tile    `json:"tile_to_chii"`
	TilesInHand [2]Tile `json:"tiles_in_hand"`
}

func (ChiiData) actionDataImpl() {}

type DrawData struct {
	DrawnTile Tile `json:"drawn_tile"`
}

func (DrawData) actionDataImpl() {}

func (action *PlayerAction) DecodeAction(rawData []byte) error {
	var raw struct {
		Action          ActionType      `json:"action_type"`
		FromPlayer      uint8           `json:"from_player"`
		PotentialAction bool            `json:"potential_action"`
		Data            json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	var ActionToDataMap = []ActionData{
		RonData{},
		TsumoData{},
		RiichiData{},
		TossData{},
		SkipData{},
		PonData{},
		KanData{},
		ChiiData{},
		DrawData{},
	}

	action.Action = raw.Action
	data := ActionToDataMap[raw.Action]
	if err := json.Unmarshal(raw.Data, &data); err != nil {
		return err
	}
	action.Data = data
	return nil
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
	Type     SetupType `json:"setup_type"`
	ToPlayer uint8     `json:"to_player"`
	Data     any       `json:"data"`
}
