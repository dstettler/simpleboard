package db

import (
    "log"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    "simpleboard/pkg/config"
	"simpleboard/internal/domain/game"
	"simpleboard/internal/domain/user"

)

var DB *gorm.DB

// Connects the SQLite db and returns the instance
func Connect(cfg *config.Config) {

	// get config db path
	db_path := ""
	if (cfg.DBPath) != "" {
		db_path = cfg.DBPath
	}

	// connect
    db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to database '%s': %v", db_path, err)
    }

	// perform migrations
	err = db.AutoMigrate(
		&game.ChessGame{},
		&user.User{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database '%s': %v", db_path, err)
	}

	// set global instance
	DB = db
}
