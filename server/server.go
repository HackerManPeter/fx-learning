package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hackermanpeter/fx-learning/internal/auth"
	"github.com/hackermanpeter/fx-learning/internal/bible"
	"github.com/hackermanpeter/fx-learning/internal/config"
	"github.com/hackermanpeter/fx-learning/internal/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewServer(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger, authService *auth.AuthService, bibleService *bible.BibleService, middleware *middleware.Middleware) *fiber.App {
	defer logger.Sync()
	app := fiber.New()

	auth.NewAuthRoutes(app, authService, middleware)
	bible.NewBibleRoutes(app, bibleService, middleware)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Sugar().Infof("Starting server on port %v", cfg.App.Port)
			go app.Listen(fmt.Sprintf(":%v", cfg.App.Port))
			return nil

		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down server")
			return app.Shutdown()
		},
	})

	return app
}
