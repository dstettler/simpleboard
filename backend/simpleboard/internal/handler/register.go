package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"

	"simpleboard/internal/domain/user"
	"simpleboard/pkg/db"
	//"simpleboard/pkg/response"
)

// Registers a new user with username, password, and email
func Register(c *gin.Context) {
	var input user.User // get user type

	// bad request; could not bind context into user type
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create entry
	if err := db.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// user successfully added
	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered",
		"user": gin.H{
			"user_id": input.UserID,
			"username": input.Username,
			"email": input.Email,
		},
	})
}
