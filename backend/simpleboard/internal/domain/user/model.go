package user

import (
//"fmt"
)

// User is an instance of a registered user
// Stores user details and login credentials
type User struct {
	UserID   uint   `gorm:"primaryKey" json:"user_id"`
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `json:"password"`
}
