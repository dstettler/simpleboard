package handler

import (
	"net/http"
	"simpleboard/internal/auth"
	"simpleboard/internal/repository"
	"simpleboard/pkg/db"
	"time"

	"github.com/gin-gonic/gin"
)

// Dashboard returns lifetime stats for the authenticated user.
// Guests are not supported — a registered account is required.
func Dashboard(c *gin.Context) {
	claims := auth.GetClaims(c)
	if claims == nil || claims.UserID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "registered account required"})
		return
	}

	var user repository.User
	if err := db.DB.First(&user, *claims.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	winRate := 0.0
	if user.TotalGames > 0 {
		winRate = float64(user.Wins) / float64(user.TotalGames)
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":        user.UserID,
		"username":       user.Username,
		"total_games":    user.TotalGames,
		"wins":           user.Wins,
		"losses":         user.Losses,
		"win_rate":       winRate,
		"current_streak": user.CurrentStreak,
		"longest_streak": user.LongestStreak,
	})
}

// gameHistoryEntry is what we return per game in the history list.
// Kept flat and small — the frontend doesn't need the full board state here.
type gameHistoryEntry struct {
	GameID     string    `json:"game_id"`
	Status     string    `json:"status"`
	PlayedAs   string    `json:"played_as"`   // "w" or "b"
	OpponentID uint      `json:"opponent_id"` // 0 if opponent was a guest
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Games returns the authenticated user's game history, newest first.
// Only finished or in-progress games where the user was a registered player are included.
func Games(c *gin.Context) {
	claims := auth.GetClaims(c)
	if claims == nil || claims.UserID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "registered account required"})
		return
	}

	uid := *claims.UserID

	// find all games where this user was white or black
	var games []repository.Game
	if err := db.DB.Where("white_player_id = ? OR black_player_id = ?", uid, uid).
		Order("created_at DESC").
		Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch games"})
		return
	}

	history := make([]gameHistoryEntry, 0, len(games))
	for _, g := range games {
		entry := gameHistoryEntry{
			GameID:    g.ID.String(),
			Status:    g.Status,
			CreatedAt: g.CreatedAt,
			UpdatedAt: g.UpdatedAt,
		}
		if g.WhitePlayerID == uid {
			entry.PlayedAs = "w"
			entry.OpponentID = g.BlackPlayerID
		} else {
			entry.PlayedAs = "b"
			entry.OpponentID = g.WhitePlayerID
		}
		history = append(history, entry)
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": uid,
		"games":   history,
	})
}
