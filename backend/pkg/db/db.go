package db

import (
	"gorm.io/gorm"
	_ "gorm.io/driver/sqlite"
)

type Store struct {
	DB *gorm.DB
}
