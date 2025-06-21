package core

type GameResult struct {
	Result          WinResult
	WonBy           uint8
	PointsTransfers []struct {
		From   uint8
		Amount uint32
	}
}

// TODO: Compute points from result
func ComputeMahjongScore(result WinResult) uint32 {
	return 0
}

// TODO
func GenerateGameResult(result WinResult, wonBy uint8) GameResult {
	return GameResult{}
}
