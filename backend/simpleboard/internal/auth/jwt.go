package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"simpleboard/pkg/config"
	"time"
)

// claims type
type Claims struct {
	UserID  *uint   `json:"user_id,omitempty"`
	GuestID *string `json:"guest_id,omitempty"`
	jwt.RegisteredClaims
}

var jwtKey []byte

// set the jwt secret on startup
func Init(cfg *config.Config) {
	jwtKey = []byte(cfg.JWTSecret)
}

// generates a new user (registered) token upon login
func NewUserToken(id uint, td time.Duration) (string, error) {
	claims := &Claims{
		UserID: &id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(td)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
}

// generates a new user (guest)
func NewGuestToken(id string, td time.Duration) (string, error) {
	claims := &Claims{
		GuestID: &id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(td)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
}
