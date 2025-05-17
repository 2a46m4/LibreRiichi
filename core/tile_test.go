package core

import (
	"testing"
)

func TestGetTileList(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want []Tile
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTileList()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetTileList() = %v, want %v", got, tt.want)
			}
		})
	}
}
