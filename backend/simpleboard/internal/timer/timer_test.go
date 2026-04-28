package timer

import (
	"testing"
	"time"

	"simpleboard/internal/chess"
	"simpleboard/internal/repository"
)

// helper to build a fresh in-progress game with a known clock state
func newGame(side string, whiteMs, blackMs int64, lastMove time.Time) *repository.Game {
	return &repository.Game{
		Status:             chess.InProgress.String(),
		Side:               side,
		TimeControlSeconds: 600,
		WhiteRemainingMs:   whiteMs,
		BlackRemainingMs:   blackMs,
		LastMoveAt:         lastMove,
	}
}

func TestInitGame_UsesProvidedControl(t *testing.T) {
	g := &repository.Game{}
	now := time.Now()

	InitGame(g, 300, now)

	if g.TimeControlSeconds != 300 {
		t.Errorf("TimeControlSeconds = %d, want 300", g.TimeControlSeconds)
	}
	if g.WhiteRemainingMs != 300_000 || g.BlackRemainingMs != 300_000 {
		t.Errorf("clocks = (%d, %d), want (300000, 300000)", g.WhiteRemainingMs, g.BlackRemainingMs)
	}
	if !g.LastMoveAt.Equal(now) {
		t.Errorf("LastMoveAt = %v, want %v", g.LastMoveAt, now)
	}
}

func TestInitGame_FallsBackToDefaultWhenZeroOrNegative(t *testing.T) {
	defaultControlSeconds = 600 // ensure baseline
	g := &repository.Game{}

	InitGame(g, 0, time.Now())
	if g.TimeControlSeconds != 600 {
		t.Errorf("zero input: got %d, want 600 (default)", g.TimeControlSeconds)
	}

	g2 := &repository.Game{}
	InitGame(g2, -5, time.Now())
	if g2.TimeControlSeconds != 600 {
		t.Errorf("negative input: got %d, want 600 (default)", g2.TimeControlSeconds)
	}
}

func TestApplyMove_DeductsActiveSideOnly(t *testing.T) {
	now := time.Now()
	// white to move, last move 2s ago
	g := newGame("w", 60_000, 60_000, now.Add(-2*time.Second))

	timedOut, _ := ApplyMove(g, now)
	if timedOut {
		t.Fatal("unexpected flag fall")
	}

	// white should lose ~2000ms; black untouched
	if g.WhiteRemainingMs > 58_100 || g.WhiteRemainingMs < 57_900 {
		t.Errorf("white remaining = %d, want ~58000", g.WhiteRemainingMs)
	}
	if g.BlackRemainingMs != 60_000 {
		t.Errorf("black remaining = %d, want 60000 (untouched)", g.BlackRemainingMs)
	}
	if !g.LastMoveAt.Equal(now) {
		t.Errorf("LastMoveAt should advance to now")
	}
}

func TestApplyMove_FlagFall(t *testing.T) {
	now := time.Now()
	// black has 500ms left, last moved 2s ago -> flag falls
	g := newGame("b", 60_000, 500, now.Add(-2*time.Second))

	timedOut, loser := ApplyMove(g, now)
	if !timedOut {
		t.Fatal("expected flag fall")
	}
	if loser != "b" {
		t.Errorf("loser = %q, want b", loser)
	}
	if g.BlackRemainingMs != 0 {
		t.Errorf("black remaining = %d, want 0", g.BlackRemainingMs)
	}
}

func TestApplyMove_ClampsBackwardClockSkew(t *testing.T) {
	now := time.Now()
	// LastMoveAt is in the future -> elapsed would be negative
	g := newGame("w", 60_000, 60_000, now.Add(5*time.Second))

	timedOut, _ := ApplyMove(g, now)
	if timedOut {
		t.Fatal("unexpected flag fall")
	}
	if g.WhiteRemainingMs != 60_000 {
		t.Errorf("white remaining = %d, want 60000 (no negative deduction)", g.WhiteRemainingMs)
	}
}

func TestLiveRemaining_DoesNotMutate(t *testing.T) {
	now := time.Now()
	g := newGame("w", 60_000, 60_000, now.Add(-3*time.Second))

	originalWhite := g.WhiteRemainingMs
	originalBlack := g.BlackRemainingMs
	originalLast := g.LastMoveAt

	_, _, _, _ = LiveRemaining(g, now)

	if g.WhiteRemainingMs != originalWhite || g.BlackRemainingMs != originalBlack {
		t.Error("LiveRemaining mutated stored clocks")
	}
	if !g.LastMoveAt.Equal(originalLast) {
		t.Error("LiveRemaining mutated LastMoveAt")
	}
}

func TestLiveRemaining_ComputesActiveSide(t *testing.T) {
	now := time.Now()
	g := newGame("b", 60_000, 60_000, now.Add(-4*time.Second))

	whiteMs, blackMs, timedOut, _ := LiveRemaining(g, now)
	if timedOut {
		t.Fatal("unexpected timeout")
	}
	if whiteMs != 60_000 {
		t.Errorf("white = %d, want 60000 (inactive)", whiteMs)
	}
	if blackMs > 56_100 || blackMs < 55_900 {
		t.Errorf("black = %d, want ~56000", blackMs)
	}
}

func TestLiveRemaining_FrozenWhenGameFinished(t *testing.T) {
	now := time.Now()
	g := newGame("w", 1_000, 60_000, now.Add(-10*time.Second))
	g.Status = chess.WinBlack.String() // already finished

	whiteMs, blackMs, timedOut, _ := LiveRemaining(g, now)
	if timedOut {
		t.Error("finished game should not flag-fall")
	}
	if whiteMs != 1_000 || blackMs != 60_000 {
		t.Errorf("clocks should be frozen, got (%d, %d)", whiteMs, blackMs)
	}
}

func TestLiveRemaining_DetectsFlagFall(t *testing.T) {
	now := time.Now()
	g := newGame("w", 100, 60_000, now.Add(-2*time.Second))

	whiteMs, _, timedOut, loser := LiveRemaining(g, now)
	if !timedOut || loser != "w" {
		t.Fatalf("expected white flag fall, got timedOut=%v loser=%q", timedOut, loser)
	}
	if whiteMs != 0 {
		t.Errorf("white live remaining = %d, want 0", whiteMs)
	}
}

func TestFlagFallStatus(t *testing.T) {
	if FlagFallStatus("w") != chess.WinBlack {
		t.Error("white flagged should give WinBlack")
	}
	if FlagFallStatus("b") != chess.WinWhite {
		t.Error("black flagged should give WinWhite")
	}
}
