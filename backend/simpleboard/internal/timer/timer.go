// Package timer enforces per-side chess clocks server-side.
// The active side's clock ticks down between moves; when it hits
// zero the player flag-falls and loses on time.
package timer

import (
	"time"

	"simpleboard/internal/chess"
	"simpleboard/internal/repository"
	"simpleboard/pkg/config"
)

// Default 10 min/side; overridden by config.Init
var defaultControlSeconds = 600

// Init wires runtime defaults from config
func Init(cfg *config.Config) {
	if cfg.DefaultTimeControlSeconds > 0 {
		defaultControlSeconds = cfg.DefaultTimeControlSeconds
	}
}

// DefaultControlSeconds returns the configured default time control
func DefaultControlSeconds() int {
	return defaultControlSeconds
}

// InitGame stamps initial clock state on a new game
func InitGame(g *repository.Game, controlSeconds int, now time.Time) {
	if controlSeconds <= 0 {
		controlSeconds = defaultControlSeconds
	}
	ms := int64(controlSeconds) * 1000

	g.TimeControlSeconds = controlSeconds
	g.WhiteRemainingMs = ms
	g.BlackRemainingMs = ms
	g.LastMoveAt = now
}

// LiveRemaining returns the up-to-date clock for both sides as of `now`.
// If the active side has flag-fallen, timedOut=true and the loser side ("w"/"b")
// is returned. Does NOT mutate the entry.
func LiveRemaining(g *repository.Game, now time.Time) (whiteMs, blackMs int64, timedOut bool, loser string) {
	whiteMs = g.WhiteRemainingMs
	blackMs = g.BlackRemainingMs

	// clock is frozen for finished games
	if g.Status != chess.InProgress.String() {
		return
	}

	elapsedMs := now.Sub(g.LastMoveAt).Milliseconds()
	if elapsedMs < 0 {
		elapsedMs = 0
	}

	if g.Side == "w" {
		whiteMs -= elapsedMs
		if whiteMs <= 0 {
			whiteMs = 0
			timedOut = true
			loser = "w"
		}
	} else {
		blackMs -= elapsedMs
		if blackMs <= 0 {
			blackMs = 0
			timedOut = true
			loser = "b"
		}
	}
	return
}

// ApplyMove deducts elapsed time from the side that just moved and rolls
// LastMoveAt forward. Call BEFORE switching sides on the chess engine.
// Returns true and the loser side if the moving player flag-fell.
func ApplyMove(g *repository.Game, now time.Time) (timedOut bool, loser string) {
	elapsedMs := now.Sub(g.LastMoveAt).Milliseconds()
	if elapsedMs < 0 {
		elapsedMs = 0
	}

	if g.Side == "w" {
		g.WhiteRemainingMs -= elapsedMs
		if g.WhiteRemainingMs <= 0 {
			g.WhiteRemainingMs = 0
			return true, "w"
		}
	} else {
		g.BlackRemainingMs -= elapsedMs
		if g.BlackRemainingMs <= 0 {
			g.BlackRemainingMs = 0
			return true, "b"
		}
	}

	g.LastMoveAt = now
	return false, ""
}

// FlagFallStatus maps the losing side to a chess.Status (the opponent wins)
func FlagFallStatus(loser string) chess.Status {
	if loser == "w" {
		return chess.WinBlack
	}
	return chess.WinWhite
}
