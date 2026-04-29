package timer

import (
	"testing"
	"time"

	"simpleboard/internal/chess"
	"simpleboard/internal/repository"
)

func TestMarkIfTimedOut_FlagsWhiteOutOfTime(t *testing.T) {
	now := time.Now()
	g := &repository.Game{
		Status:           chess.InProgress.String(),
		Side:             "w",
		WhiteRemainingMs: 100,
		BlackRemainingMs: 60_000,
		LastMoveAt:       now.Add(-5 * time.Second),
	}

	changed := markIfTimedOut(g, now)
	if !changed {
		t.Fatal("expected timeout to be marked")
	}
	if g.Status != chess.WinBlack.String() {
		t.Errorf("status = %q, want WinBlack", g.Status)
	}
	if g.WhiteRemainingMs != 0 {
		t.Errorf("white remaining = %d, want 0", g.WhiteRemainingMs)
	}
}

func TestMarkIfTimedOut_LeavesActiveGameAlone(t *testing.T) {
	now := time.Now()
	g := &repository.Game{
		Status:           chess.InProgress.String(),
		Side:             "b",
		WhiteRemainingMs: 60_000,
		BlackRemainingMs: 60_000,
		LastMoveAt:       now.Add(-1 * time.Second),
	}

	if markIfTimedOut(g, now) {
		t.Fatal("active game wrongly flagged as timed out")
	}
	if g.Status != chess.InProgress.String() {
		t.Errorf("status changed to %q on a non-timeout", g.Status)
	}
}

func TestMarkIfTimedOut_IgnoresFinishedGames(t *testing.T) {
	now := time.Now()
	// stale clocks would normally flag, but the game is already won
	g := &repository.Game{
		Status:           chess.WinWhite.String(),
		Side:             "b",
		WhiteRemainingMs: 30_000,
		BlackRemainingMs: 100,
		LastMoveAt:       now.Add(-10 * time.Minute),
	}

	if markIfTimedOut(g, now) {
		t.Fatal("finished game must not be re-flagged")
	}
	if g.Status != chess.WinWhite.String() {
		t.Errorf("status = %q, want WinWhite (untouched)", g.Status)
	}
}
