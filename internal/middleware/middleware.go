package middleware

import (
	"github.com/hackermanpeter/fx-learning/internal/cache"
	"github.com/hackermanpeter/fx-learning/internal/config"
	"github.com/hackermanpeter/fx-learning/internal/database"
	"go.uber.org/zap"
)

type Middleware struct {
	cfg    *config.Config
	logger *zap.Logger
	db     *database.Database
	cache  *cache.Cache
}

func NewMiddlewareService(cfg *config.Config, logger *zap.Logger, db *database.Database, cache *cache.Cache) *Middleware {
	return &Middleware{
		cfg,
		logger,
		db,
		cache,
	}
}
