package domain

import (
	"github.com/google/uuid"
)

type Queue struct {
	GameID        uuid.UUID `json:"game_id"`
	WhitePlayerID uint      `json:"white_player_id"`
	BlackPlayerID uint      `json:"black_player_id"`
	WhiteGuestID  string    `json:"white_guest_id"`
	BlackGuestID  string    `json:"black_guest_id"`
	Active        bool      `json:"active"`
	Open          bool      `json:"open"`
}
