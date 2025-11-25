package game

import (
	"testing"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

func TestWindOffset_GetPlayerWind(t *testing.T) {
	tests := []struct {
		name     string
		offset   WindState
		gameIdx  uint8
		expected Wind
	}{
		{
			name:     "Offset 0, Player 0 should be East",
			offset:   0,
			gameIdx:  0,
			expected: East,
		},
		{
			name:     "Offset 0, Player 1 should be South",
			offset:   0,
			gameIdx:  1,
			expected: South,
		},
		{
			name:     "Offset 0, Player 2 should be West",
			offset:   0,
			gameIdx:  2,
			expected: West,
		},
		{
			name:     "Offset 0, Player 3 should be North",
			offset:   0,
			gameIdx:  3,
			expected: North,
		},
		{
			name:     "Offset 1, Player 0 should be South",
			offset:   1,
			gameIdx:  0,
			expected: South,
		},
		{
			name:     "Offset 1, Player 1 should be West",
			offset:   1,
			gameIdx:  1,
			expected: West,
		},
		{
			name:     "Offset 1, Player 2 should be North",
			offset:   1,
			gameIdx:  2,
			expected: North,
		},
		{
			name:     "Offset 1, Player 3 should be East",
			offset:   1,
			gameIdx:  3,
			expected: East,
		},
		{
			name:     "Offset 2, Player 0 should be West",
			offset:   2,
			gameIdx:  0,
			expected: West,
		},
		{
			name:     "Offset 2, Player 1 should be North",
			offset:   2,
			gameIdx:  1,
			expected: North,
		},
		{
			name:     "Offset 2, Player 2 should be East",
			offset:   2,
			gameIdx:  2,
			expected: East,
		},
		{
			name:     "Offset 2, Player 3 should be South",
			offset:   2,
			gameIdx:  3,
			expected: South,
		},
		{
			name:     "Offset 3, Player 0 should be North",
			offset:   3,
			gameIdx:  0,
			expected: North,
		},
		{
			name:     "Offset 3, Player 1 should be East",
			offset:   3,
			gameIdx:  1,
			expected: East,
		},
		{
			name:     "Offset 3, Player 2 should be South",
			offset:   3,
			gameIdx:  2,
			expected: South,
		},
		{
			name:     "Offset 3, Player 3 should be West",
			offset:   3,
			gameIdx:  3,
			expected: West,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.offset.GetPlayerWind(tt.gameIdx)
			if result != tt.expected {
				t.Errorf("WindOffset.GetPlayerWind() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWindOffset_WrapAround(t *testing.T) {
	// Test that the wind calculation wraps around correctly for all 16 combinations (4 offsets * 4 players)
	for offset := uint8(0); offset < 4; offset++ {
		for gameIdx := uint8(0); gameIdx < 4; gameIdx++ {
			wind := WindState(offset).GetPlayerWind(gameIdx)

			// Verify the result is a valid wind
			validWinds := []Wind{East, South, West, North}
			isValid := false
			for _, validWind := range validWinds {
				if wind == validWind {
					isValid = true
					break
				}
			}

			if !isValid {
				t.Errorf("Offset %d, Player %d produced invalid wind: %v", offset, gameIdx, wind)
			}
		}
	}
}

func TestWindOffset_CyclicProperty(t *testing.T) {
	// Test that applying the offset 4 times gets back to the same wind
	offset := WindState(2) // Choose an arbitrary offset

	for gameIdx := uint8(0); gameIdx < 4; gameIdx++ {
		originalWind := offset.GetPlayerWind(gameIdx)

		// Apply full rotation (4 steps)
		wind1 := offset.GetPlayerWind((gameIdx + 1) % 4)
		wind2 := offset.GetPlayerWind((gameIdx + 2) % 4)
		wind3 := offset.GetPlayerWind((gameIdx + 3) % 4)
		wind4 := offset.GetPlayerWind((gameIdx + 4) % 4) // Should be same as original

		// Verify we get all four different winds in order
		winds := []Wind{originalWind, wind1, wind2, wind3}
		uniqueWinds := make(map[Wind]bool)
		for _, w := range winds {
			uniqueWinds[w] = true
		}

		if len(uniqueWinds) != 4 {
			t.Errorf("Expected 4 unique winds for player %d with offset %d, got %d unique winds",
				gameIdx, offset, len(uniqueWinds))
		}

		// Verify full rotation returns to original
		if wind4 != originalWind {
			t.Errorf("Expected wind after full rotation to be %v, got %v", originalWind, wind4)
		}
	}
}

func TestWindOffset_EdgeCases(t *testing.T) {
	// Test edge case: maximum valid values
	t.Run("Maximum values", func(t *testing.T) {
		offset := WindState(3)
		gameIdx := uint8(3)
		expected := West // (3 + 3) % 4 = 6 % 4 = 2, which is West (0=East,1=South,2=West,3=North)
		result := offset.GetPlayerWind(gameIdx)
		if result != expected {
			t.Errorf("WindOffset.GetPlayerWind() with max values = %v, want %v", result, expected)
		}
	})
}

// Benchmark tests
func BenchmarkWindOffset_GetPlayerWind(b *testing.B) {
	offset := WindState(2)
	gameIdx := uint8(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		offset.GetPlayerWind(gameIdx)
	}
}

func BenchmarkWindOffset_GetPlayerWind_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		offset := WindState(2)
		gameIdx := uint8(1)

		for pb.Next() {
			offset.GetPlayerWind(gameIdx)
		}
	})
}
