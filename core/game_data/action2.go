package core

import "encoding/json"

type ActionMessage struct {
	ActionType ActionType `json:"action_type"`
	Data       any        `json:"data"`
}

type Game interface {
	CurrentPlayerIdx() uint8
	IsPendingAction(Action, fromPlayer uint8)
	ApplyGameResult(result GameResult)
}

type Action interface {
	PerformAction(game Game, fromPlayer uint8)
}

type ActionWrapper struct {
	Action
}

func (action ActionWrapper) MarshalJSON() ([]byte, error) {
	switch action.Action.(type) {
	case Ron:
		json.Marshal(ActionMessage{
			RON,
			action,
		})
	}
	return nil, nil
}

type Ron struct {
	TileToRon Tile      `json:"tile_to_ron"`
	WinResult WinResult `json:"win_result"` // TODO: Remove
}

func (ron Ron) PerformAction() {

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
