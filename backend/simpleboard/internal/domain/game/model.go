package game

import (
	//"gorm.io/gorm"
	"time"
)

// Game is an instance of an active game session
// Stores player IDs, game state, and timestamps
type Game struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	WhitePlayerID uint      `json:"white_player_id"`
	BlackPlayerID uint      `json:"black_player_id"`
	State         string    `json:"state"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
