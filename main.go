package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hackermanpeter/fx-learning/internal/auth"
	"github.com/hackermanpeter/fx-learning/internal/bible"
	"github.com/hackermanpeter/fx-learning/internal/cache"
	"github.com/hackermanpeter/fx-learning/internal/config"
	"github.com/hackermanpeter/fx-learning/internal/database"
	"github.com/hackermanpeter/fx-learning/internal/middleware"
	"github.com/hackermanpeter/fx-learning/server"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())

}

func main() {
	fx.New(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Provide(
			// App level
			config.NewConfig,
			server.NewServer,
			zap.NewProduction,
			database.NewDatabase,
			NewValidator,
			cache.NewCache,

			// modules
			auth.NewAuthService,
			bible.NewBibleService,
			middleware.NewMiddlewareService,
		),
		fx.Invoke(
			func(*database.Database) {},
			func(*fiber.App) {},
		),
	).Run()

}
