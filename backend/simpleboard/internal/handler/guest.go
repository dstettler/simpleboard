package handler

import (
	"fmt"
	"net/http"

	"simpleboard/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

// Guest login endpoint
// Serves new guest tokens to GET requests without Auth
func Guest(c *gin.Context) {

	// check auth claims;
	// cannot generate guest tokens if bearing existing authentication
	claims := auth.GetClaims(c)
	if claims == nil || (claims.UserID == nil && claims.GuestID == nil) {
		// serve new guest token
		guestID := uuid.New()
		token, err := auth.NewGuestToken(guestID, 24*time.Hour)
		if err != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "guest creation successful",
			"token":   token,
			"user": gin.H{
				"guest_id": guestID,
			},
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "guest not allowed - existing claims"})
		return
	}
}
