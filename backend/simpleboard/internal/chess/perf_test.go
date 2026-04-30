package chess

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// perfPositions are representative board positions spanning different move-count
var perfPositions = []struct {
	name string
	fen  string
}{
	{
		name: "start (20 moves)",
		fen:  StartFEN,
	},
	{
		name: "kiwipete (48 moves)",
		fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
	},
	{
		name: "endgame sparse",
		fen:  "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
	},
}

// perft counts leaf nodes in the move tree to the given depth.
//	perft(pos, 1) = number of legal moves from pos
//	perft(pos, n) = sum of perft(child, n-1) over all legal moves
// Known starting-position values: depth 1=20, 2=400, 3=8902, 4=197281.

func perft(g ChessGame, depth int) int {
	if depth <= 0 {
		return 1
	}
	moves := g.LegalMoves()
	if depth == 1 {
		return len(moves)
	}
	total := 0
	for _, m := range moves {
		child := g.Copy()
		child.MakeMove(m)
		total += perft(child, depth-1)
	}
	return total
}

// Go benchmark functions  (go test -bench=. ./internal/chess/)
// Results include ns/op, B/op, allocs/op
// =============================================================================

func BenchmarkLegalMoves(b *testing.B) {
	for _, pos := range perfPositions {
		b.Run(pos.name, func(b *testing.B) {
			game := ReadChessGame(pos.fen, nil, nil)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = game.LegalMoves()
			}
		})
	}
}

func BenchmarkPositionMoves(b *testing.B) {
	for _, pos := range perfPositions {
		b.Run(pos.name, func(b *testing.B) {
			game := ReadChessGame(pos.fen, nil, nil)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = game.PositionMoves()
			}
		})
	}
}

func BenchmarkPerft(b *testing.B) {
	cases := []struct {
		depth int
		nodes int
	}{
		{1, 20},
		{2, 400},
		{3, 8902},
		{4, 197281},
	}
	for _, tc := range cases {
		b.Run(fmt.Sprintf("depth%d", tc.depth), func(b *testing.B) {
			game := ReadChessGame(StartFEN, nil, nil)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if got := perft(game, tc.depth); got != tc.nodes {
					b.Fatalf("perft(%d) = %d, want %d", tc.depth, got, tc.nodes)
				}
			}
		})
	}
}

// Verbose tests  (go test -v -run TestPerf ./internal/chess/)
// print throughput and memory stats via t.Logf
// =============================================================================

func TestPerfLegalMoves(t *testing.T) {
	const iterations = 10_000

	for _, pos := range perfPositions {
		t.Run(pos.name, func(t *testing.T) {
			game := ReadChessGame(pos.fen, nil, nil)
			moveCount := len(game.LegalMoves())

			var before, after runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&before)
			start := time.Now()

			for i := 0; i < iterations; i++ {
				_ = game.LegalMoves()
			}

			elapsed := time.Since(start)
			runtime.GC()
			runtime.ReadMemStats(&after)

			allocPerCall := int64(after.TotalAlloc-before.TotalAlloc) / iterations
			callsPerSec := float64(iterations) / elapsed.Seconds()

			t.Logf("legal moves=%-3d  iterations=%d  elapsed=%-10v  %.0f calls/sec  %d B/call",
				moveCount, iterations, elapsed.Round(time.Millisecond), callsPerSec, allocPerCall)
		})
	}
}

// TestPerfPerft enumerates the move tree from the starting position at depths
func TestPerfPerft(t *testing.T) {
	cases := []struct {
		depth     int
		wantNodes int
	}{
		{1, 20},
		{2, 400},
		{3, 8902},
	}
	game := ReadChessGame(StartFEN, nil, nil)

	for _, tc := range cases {
		t.Run(fmt.Sprintf("depth%d", tc.depth), func(t *testing.T) {
			var before, after runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&before)
			start := time.Now()

			nodes := perft(game, tc.depth)

			elapsed := time.Since(start)
			runtime.GC()
			runtime.ReadMemStats(&after)

			if nodes != tc.wantNodes {
				t.Errorf("perft(%d) = %d, want %d", tc.depth, nodes, tc.wantNodes)
			}

			allocBytes := after.TotalAlloc - before.TotalAlloc
			nodesPerSec := float64(nodes) / elapsed.Seconds()

			t.Logf("depth=%d  nodes=%-6d  elapsed=%-12v  %.0f nodes/sec  %d B allocated",
				tc.depth, nodes, elapsed.Round(time.Microsecond), nodesPerSec, allocBytes)
		})
	}
}
