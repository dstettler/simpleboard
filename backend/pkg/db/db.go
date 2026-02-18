package db

import (
	"gorm.io/gorm"
)

type Store struct {
	DB *gorm.DB
}
