package repository

import (
// "fmt"
)

// User is an instance of a registered user
// Stores user details and login credentials
type User struct {
	UserID   uint   `gorm:"uniqueIndex;primaryKey;autoIncrement"`
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `json:"password"`
}
