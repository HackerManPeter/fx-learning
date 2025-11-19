package database

import (
	"context"
	"fmt"

	"github.com/hackermanpeter/fx-learning/internal/config"
	"github.com/hackermanpeter/fx-learning/internal/models"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) *Database {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Sugar().Fatalf("unable to connect to database")
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Migrating Database")
			return db.AutoMigrate(&models.User{})
		},

		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Database connection")
			sqlDB, err := db.DB()
			if err != nil {
				logger.Error("unable to cleanly close database connection")
			}

			return sqlDB.Close()
		},
	})

	return &Database{
		DB: db,
	}
}
