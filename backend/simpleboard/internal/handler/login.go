package handler

import (
	"net/http"
	//"simpleboard/internal/domain"
	"simpleboard/internal/repository"
	"simpleboard/pkg/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	//"simpleboard/pkg/response"
)

// Logs a user in with username and password
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// bad request; could not bind context into input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get user
	var user repository.User
	if err := db.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// user successfully logged in
	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user": gin.H{
			"user_id":  user.UserID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
