package score

import (
	"testing"
)

func TestComputePoints_BasicScenarios(t *testing.T) {
	tests := []struct {
		name      string
		points    Points
		isDealer  bool
		expected  PointValue
	}{
		{
			name:     "1 han 30 fu non-dealer",
			points:   Points{han: 1, fu: 30},
			isDealer: false,
			expected: PointValue{
				NonDealerTsumo: 300,  // 30 * 2^1 * 1 = 60, rounded to 100, then *3 = 300
				DealerTsumo:   600,  // 30 * 2^1 * 2 = 120, rounded to 100, then *6 = 600
				Ron:           1000, // 30 * 2^1 * 4 = 240, rounded to 100, then *10 = 1000
			},
		},
		{
			name:     "1 han 30 fu dealer",
			points:   Points{han: 1, fu: 30},
			isDealer: true,
			expected: PointValue{
				NonDealerTsumo: 600,  // 30 * 2^1 * 2 = 120, rounded to 100, then *6 = 600
				DealerTsumo:   0,    // Invalid since winner is dealer
				Ron:           1500, // 30 * 2^1 * 6 = 360, rounded to 100, then *15 = 1500
			},
		},
		{
			name:     "2 han 30 fu non-dealer",
			points:   Points{han: 2, fu: 30},
			isDealer: false,
			expected: PointValue{
				NonDealerTsumo: 700,  // 30 * 2^2 * 1 = 120, rounded to 100, then *7 = 700
				DealerTsumo:   1300, // 30 * 2^2 * 2 = 240, rounded to 100, then *13 = 1300
				Ron:           2600, // 30 * 2^2 * 4 = 480, rounded to 100, then *26 = 2600
			},
		},
		{
			name:     "3 han 40 fu non-dealer",
			points:   Points{han: 3, fu: 40},
			isDealer: false,
			expected: PointValue{
				NonDealerTsumo: 1300, // 40 * 2^3 * 1 = 320, rounded to 100, then *13 = 1300
				DealerTsumo:   2600, // 40 * 2^3 * 2 = 640, rounded to 100, then *26 = 2600
				Ron:           5200, // 40 * 2^3 * 4 = 1280, rounded to 100, then *52 = 5200
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComputePoints(tt.points, tt.isDealer)
			if result.NonDealerTsumo != tt.expected.NonDealerTsumo {
				t.Errorf("NonDealerTsumo = %d, want %d", result.NonDealerTsumo, tt.expected.NonDealerTsumo)
			}
			if result.DealerTsumo != tt.expected.DealerTsumo {
				t.Errorf("DealerTsumo = %d, want %d", result.DealerTsumo, tt.expected.DealerTsumo)
			}
			if result.Ron != tt.expected.Ron {
				t.Errorf("Ron = %d, want %d", result.Ron, tt.expected.Ron)
			}
		})
	}
}


func TestComputePoints_ManganAndAbove(t *testing.T) {
	tests := []struct {
		name     string
		points   Points
		isDealer bool
		expected PointValue
	}{
		{
			name:     "Mangan (5 han) non-dealer",
			points:   Points{han: 5, fu: 30},
			isDealer: false,
			expected: PointValue{
				NonDealerTsumo: 2000, // 2000 * 1 = 2000
				DealerTsumo:   4000, // 2000 * 2 = 4000
				Ron:           8000, // 2000 * 4 = 8000
			},
		},
		{
			name:     "Mangan (5 han) dealer",
			points:   Points{han: 5, fu: 30},
			isDealer: true,
			expected: PointValue{
				NonDealerTsumo: 4000, // 2000 * 2 = 4000
				DealerTsumo:   0,    // Invalid
				Ron:           12000, // 2000 * 6 = 12000
			},
		},
		{
			name:     "Haneman (6 han) non-dealer",
			points:   Points{han: 6, fu: 30},
			isDealer: false,
			expected: PointValue{
				NonDealerTsumo: 3000, // 2000 * 1.5 = 3000
				DealerTsumo:   6000, // 2000 * 2 * 1.5 = 6000
				Ron:           12000, // 2000 * 4 * 1.5 = 12000
			},
		},
		{
			name:     "Haneman (7 han) non-dealer",
			points:   Points{han: 7, fu: 30},
			isDealer: false,
			expected: PointValue{
				NonDealerTsumo: 3000, // 2000 * 1.5 = 3000
				DealerTsumo:   6000, // 2000 * 2 * 1.5 = 6000
				Ron:           12000, // 2000 * 4 * 1.5 = 12000
			},
		},
		{
			name:     "Baiman (8 han) non-dealer",
			points:   Points{han: 8, fu: 30},
			isDealer: false,
			expected: PointValue{
				NonDealerTsumo: 4000, // 2000 * 2 = 4000
				DealerTsumo:   8000, // 2000 * 2 * 2 = 8000
				Ron:           16000, // 2000 * 4 * 2 = 16000
			},
		},
		{
			name:     "Sanbaiman (11 han) non-dealer",
			points:   Points{han: 11, fu: 30},
			isDealer: false,
			expected: PointValue{
				NonDealerTsumo: 6000, // 2000 * 3 = 6000
				DealerTsumo:   12000, // 2000 * 2 * 3 = 12000
				Ron:           24000, // 2000 * 4 * 3 = 24000
			},
		},
		{
			name:     "Yakuman (13 han) non-dealer",
			points:   Points{han: 13, fu: 30},
			isDealer: false,
			expected: PointValue{
				NonDealerTsumo: 8000, // 2000 * 4 = 8000
				DealerTsumo:   16000, // 2000 * 2 * 4 = 16000
				Ron:           32000, // 2000 * 4 * 4 = 32000
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComputePoints(tt.points, tt.isDealer)
			if result.NonDealerTsumo != tt.expected.NonDealerTsumo {
				t.Errorf("NonDealerTsumo = %d, want %d", result.NonDealerTsumo, tt.expected.NonDealerTsumo)
			}
			if result.DealerTsumo != tt.expected.DealerTsumo {
				t.Errorf("DealerTsumo = %d, want %d", result.DealerTsumo, tt.expected.DealerTsumo)
			}
			if result.Ron != tt.expected.Ron {
				t.Errorf("Ron = %d, want %d", result.Ron, tt.expected.Ron)
			}
		})
	}
}

func TestComputePoints_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		points   Points
		isDealer bool
	}{
		{
			name:     "Minimum values (1 han 20 fu)",
			points:   Points{han: 1, fu: 20},
			isDealer: false,
		},
		{
			name:     "High han but low fu (4 han 20 fu)",
			points:   Points{han: 4, fu: 20},
			isDealer: false,
		},
		{
			name:     "Very high han (20 han)",
			points:   Points{han: 20, fu: 30},
			isDealer: false,
		},
		{
			name:     "Very high fu (110 fu)",
			points:   Points{han: 3, fu: 110},
			isDealer: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComputePoints(tt.points, tt.isDealer)

			// Basic sanity checks
			if tt.isDealer && result.DealerTsumo != 0 {
				t.Errorf("DealerTsumo should be 0 when winner is dealer, got %d", result.DealerTsumo)
			}
			if result.Ron == 0 {
				t.Errorf("Ron should never be 0, got %d", result.Ron)
			}
			if !tt.isDealer && result.DealerTsumo == 0 {
				t.Errorf("DealerTsumo should not be 0 when winner is not dealer")
			}
		})
	}
}

func TestComputePoints_DealerTsumoInvalid(t *testing.T) {
	tests := []struct {
		name     string
		points   Points
	}{
		{"1 han 30 fu", Points{han: 1, fu: 30}},
		{"3 han 40 fu", Points{han: 3, fu: 40}},
		{"5 han mangan", Points{han: 5, fu: 30}},
		{"13 han yakuman", Points{han: 13, fu: 30}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComputePoints(tt.points, true)
			if result.DealerTsumo != 0 {
				t.Errorf("DealerTsumo should be 0 when winner is dealer, got %d", result.DealerTsumo)
			}
		})
	}
}

func TestComputePoints_RoundingBehavior(t *testing.T) {
	// Test that scores are properly rounded up to the next 100
	tests := []struct {
		name           string
		points         Points
		isDealer       bool
		expectedRon    uint
		description    string
	}{
		{
			name:        "1 han 40 fu non-dealer should round up",
			points:      Points{han: 1, fu: 40},
			isDealer:    false,
			expectedRon: 1300, // 40 * 2^2 * 4 = 640, should round up to 1300
			description: "Base score: 160, rounded to 200, then *4 = 800, but should be higher due to proper calculation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComputePoints(tt.points, tt.isDealer)
			if result.Ron != tt.expectedRon {
				t.Errorf("Ron = %d, want %d (%s)", result.Ron, tt.expectedRon, tt.description)
			}
		})
	}
}

func TestComputePoints_Consistency(t *testing.T) {
	// Test that non-dealer tsumo payments * 3 = dealer tsumo payment * 2 = ron payment
	testCases := []Points{
		{han: 1, fu: 30},
		{han: 2, fu: 30},
		{han: 3, fu: 40},
		{han: 4, fu: 30},
	}

	for _, points := range testCases {
		t.Run("Consistency check", func(t *testing.T) {
			result := ComputePoints(points, false)

			// For non-dealer winner: non-dealer * 3 + dealer * 2 = ron
			// Actually, this should be: non-dealer * 3 + dealer * 1 = ron
			// Wait, let me check the logic again...
			// In riichi mahjong: when non-dealer wins by tsumo:
			// - 2 non-dealers pay non-dealer amount
			// - 1 dealer pays dealer amount
			// So: 2 * non-dealer + 1 * dealer = ron
			expectedRon := 2*result.NonDealerTsumo + result.DealerTsumo

			if result.Ron != expectedRon {
				t.Errorf("Ron consistency failed: Ron=%d, 2*NonDealerTsumo+DealerTsumo=%d (NonDealerTsumo=%d, DealerTsumo=%d)",
					result.Ron, expectedRon, result.NonDealerTsumo, result.DealerTsumo)
			}
		})
	}
}

func TestComputePoints_DealerConsistency(t *testing.T) {
	// Test dealer win consistency
	testCases := []Points{
		{han: 1, fu: 30},
		{han: 2, fu: 30},
		{han: 3, fu: 40},
		{han: 5, fu: 30},
	}

	for _, points := range testCases {
		t.Run("Dealer consistency check", func(t *testing.T) {
			result := ComputePoints(points, true)

			// For dealer winner: 3 non-dealers pay non-dealer amount
			// So: 3 * non-dealer = ron
			expectedRon := 3 * result.NonDealerTsumo

			if result.Ron != expectedRon {
				t.Errorf("Dealer Ron consistency failed: Ron=%d, 3*NonDealerTsumo=%d (NonDealerTsumo=%d)",
					result.Ron, expectedRon, result.NonDealerTsumo)
			}
		})
	}
}

// Benchmark tests
func BenchmarkComputePoints_Basic(b *testing.B) {
	points := Points{han: 3, fu: 40}
	isDealer := false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ComputePoints(points, isDealer)
	}
}

func BenchmarkComputePoints_Dealer(b *testing.B) {
	points := Points{han: 5, fu: 30}
	isDealer := true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ComputePoints(points, isDealer)
	}
}

func BenchmarkComputePoints_Parallel(b *testing.B) {
	points := Points{han: 3, fu: 40}
	isDealer := false

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ComputePoints(points, isDealer)
		}
	})
}

func BenchmarkComputePoints_Mangan(b *testing.B) {
	points := Points{han: 13, fu: 30}
	isDealer := false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ComputePoints(points, isDealer)
	}
}
