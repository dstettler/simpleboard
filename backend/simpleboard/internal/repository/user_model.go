package repository

import "time"

// User is an instance of a registered user;
// Also stores readily available information on user
// game statistics
type User struct {
	UserID        uint      `gorm:"uniqueIndex;primaryKey;autoIncrement"`
	Username      string    `gorm:"uniqueIndex;not null" json:"username"`
	Email         string    `gorm:"uniqueIndex;not null" json:"email"`
	Password      string    `json:"password"`
	LastLoginDate time.Time `json:"last_login_date"`
	CurrentStreak int       `json:"current_streak"`
	LongestStreak int       `json:"longest_streak"`
	TotalGames    int       `json:"total_games"`
	Wins          int       `json:"wins"`
	Losses        int       `json:"losses"`
}
