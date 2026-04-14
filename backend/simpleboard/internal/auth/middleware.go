package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

const claimsKey = "authClaims"

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next() // default - no auth endpoints
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ") // auth

		claims := &Claims{}

		// parse with claims
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// bad token
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// passed auth - set claims in context
		c.Set(claimsKey, claims)
		c.Next()
	}
}

// helper function to get claims from gin context
func GetClaims(c *gin.Context) *Claims {
	v, ok := c.Get(claimsKey)
	if !ok {
		return nil
	}

	claims, ok := v.(*Claims)

	if !ok {
		return nil
	}
	return claims
}
