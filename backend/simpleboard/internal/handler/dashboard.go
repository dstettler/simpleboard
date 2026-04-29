package handler

import (
	"net/http"
	"simpleboard/internal/auth"
	"simpleboard/internal/repository"
	"simpleboard/pkg/db"

	"github.com/gin-gonic/gin"
)

// Dashboard returns stats for the authenticated user.
// Guests are not supported — a registered account is required.
//
// Response shape (for frontend alignment):
//
//	{
//	  "user_id":        uint,
//	  "username":       string,
//	  "total_games":    int,
//	  "wins":           int,
//	  "losses":         int,
//	  "win_rate":       float64,   // 0.0–1.0; 0 if no games played
//	  "current_streak": int,
//	  "longest_streak": int
//	}
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
