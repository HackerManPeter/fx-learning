package bible

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hackermanpeter/fx-learning/internal/middleware"
)

func NewBibleRoutes(app fiber.Router, bibleService *BibleService, ms *middleware.Middleware) {
	app = app.Group("/bible").Use(ms.AuthMiddleware())

	app.Get("/random", bibleService.GetRandomVerse())

}
