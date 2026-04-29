package repository

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

// Game is an instance of an active game session
// Stores player IDs, game state, and timestamps
type Game struct {
	ID                 uuid.UUID      `gorm:"uniqueIndex;primaryKey`
	WhitePlayerID      uint           `json:"white_player_id"`
	BlackPlayerID      uint           `json:"black_player_id"`
	WhiteGuestID       string         `json:"white_guest_id"`
	BlackGuestID       string         `json:"black_guest_id"`
	State              string         `json:"state"`
	Status             string         `json:"status"`
	Side               string         `gorm:"size:1; not null"`
	NextMoves          datatypes.JSON `gorm:"type:json"`
	PrevMoves          datatypes.JSON `gorm:"type:json"`
	TimeControlSeconds int            `json:"time_control_seconds"`
	WhiteRemainingMs   int64          `json:"white_remaining_ms"`
	BlackRemainingMs   int64          `json:"black_remaining_ms"`
	LastMoveAt         time.Time      `json:"last_move_at"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

func (g *Game) BeforeCreate(tx *gorm.DB) (err error) {
	g.ID = uuid.New()
	return
}
