package handler

import (
	"log"
	"net/http"
	"simpleboard/internal/auth"
	"simpleboard/internal/repository"
	"simpleboard/pkg/db"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Logs a user in with username and password
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user repository.User
	if err := db.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	updateStreak(&user)

	token, err := auth.NewUserToken(user.UserID, 24*time.Hour)
	if err != nil {
		log.Printf("token creation failed for user %d: %v", user.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user": gin.H{
			"user_id":        user.UserID,
			"username":       user.Username,
			"email":          user.Email,
			"current_streak": user.CurrentStreak,
			"longest_streak": user.LongestStreak,
		},
	})
}

// updateStreak checks for new daily login and increments the streak
func updateStreak(user *repository.User) {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	last := user.LastLoginDate.UTC().Truncate(24 * time.Hour)

	if last.Equal(today) {
		return
	}

	yesterday := today.Add(-24 * time.Hour)
	if last.Equal(yesterday) {
		user.CurrentStreak++ // kept the streak alive
	} else {
		user.CurrentStreak = 1 // gap or first ever login
	}

	if user.CurrentStreak > user.LongestStreak {
		user.LongestStreak = user.CurrentStreak
	}
	user.LastLoginDate = today

	db.DB.Save(user)
}
