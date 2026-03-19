package domain

import (
// "fmt"
)

// User is an instance of a registered user
// Stores user details and login credentials
type User struct {
	UserID   uint
	Username string
	Email    string
	Password string
}
