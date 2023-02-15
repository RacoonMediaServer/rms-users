package db

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/configuration"
	"github.com/RacoonMediaServer/rms-users/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database represents all database methods
type Database interface {
	Users
}

type database struct {
	conn *gorm.DB
}

func Connect(config configuration.Database) (Database, error) {
	db, err := gorm.Open(postgres.Open(config.GetConnectionString()))
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return database{conn: db}, nil
}
