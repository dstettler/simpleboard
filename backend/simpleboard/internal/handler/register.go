package handler

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"simpleboard/internal/domain"
	"simpleboard/internal/repository"
	"simpleboard/pkg/db"
)

// Registers a new user with username, password, and email
func Register(c *gin.Context) {
	var input domain.User // get user type

	// bad request; could not bind context into user type
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save password"})
		return
	}

	user := repository.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashed),
	}

	// create entry
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// user successfully added
	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered",
		"user": gin.H{
			"user_id":  user.UserID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
