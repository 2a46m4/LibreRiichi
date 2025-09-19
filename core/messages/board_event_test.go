package core

import (
	"encoding/json"
	"slices"
	"testing"
)

func TestBoardEventMarshalling(t *testing.T) {
	bytes, err := json.Marshal(PlayerActionEvent{
		Action:     nil,
		FromPlayer: 0,
	})

	if err != nil {
		t.Fail()
	}

	if !slices.Equal(bytes, []byte("")) {
		t.Fail()
	}
}

