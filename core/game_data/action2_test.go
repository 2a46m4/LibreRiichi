package core

import (
	"encoding/json"
	"testing"
)

func TestAction(t *testing.T) {
	var a Action = Ron{
		TileToRon: Green,
		WinResult: WinResult{},
	}

	a.PerformAction()

	bytes, err := json.Marshal(a)
	if err != nil {
		t.Errorf("%v\n", err.Error())
	}

	t.Log(string(bytes))

}
