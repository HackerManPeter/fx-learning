package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hackermanpeter/fx-learning/internal/middleware"
)

func NewAuthRoutes(app fiber.Router, authService *AuthService, middleware *middleware.Middleware) {
	app = app.Group("/auth")

	app.Post("/login", authService.login())
	app.Post("/sign-up", authService.signUp())
	app.Post("/logout", middleware.AuthMiddleware(), authService.logout())
}
