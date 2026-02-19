package game

import (
	//"gorm.io/gorm"
	"time"
)

// ChessGame is an instance of an active game session
// Stores player IDs, game state, and timestamps
type ChessGame struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	WhiteUserID uint      `json:"white_player_id"`
	BlackUserID uint      `json:"black_player_id"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
