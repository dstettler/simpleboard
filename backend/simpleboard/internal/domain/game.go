package domain

import (
	"time"
)

// Game is an instance of an active game session
// Stores player IDs, game state, and timestamps
type Game struct {
	WhitePlayerID uint
	BlackPlayerID uint
	WhiteGuestID  string
	BlackGuestID  string
	State         string
	Status        string
	Side          string
	NextMoves     []string
	PrevMoves     []string

	// Per-side chess clock; only the active side's time ticks down
	TimeControlSeconds int
	WhiteRemainingMs   int64
	BlackRemainingMs   int64
	LastMoveAt         time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
