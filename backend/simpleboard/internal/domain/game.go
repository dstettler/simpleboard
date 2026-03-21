package domain

import (
	//"gorm.io/gorm"
	"time"
)

// Game is an instance of an active game session
// Stores player IDs, game state, and timestamps
type Game struct {
	ID            uint
	WhitePlayerID uint
	BlackPlayerID uint
	State         string
	Status        string
	Side          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
