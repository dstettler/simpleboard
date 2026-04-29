package repository

import (
	"github.com/google/uuid"
	//"gorm.io/datatypes"
	//"gorm.io/gorm"
	"time"
)

type Queue struct {
	ID            uint      `gorm:"uniqueIndex;primaryKey;autoIncrement`
	GameID        uuid.UUID `json:"game_id"`
	WhitePlayerID uint      `json:"white_player_id"`
	BlackPlayerID uint      `json:"black_player_id"`
	WhiteGuestID  string    `json:"white_guest_id"`
	BlackGuestID  string    `json:"black_guest_id"`
	Active        bool      `json:"active"`
	Open          bool      `json:"open"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
