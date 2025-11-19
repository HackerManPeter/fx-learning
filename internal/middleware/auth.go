package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hackermanpeter/fx-learning/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (m *Middleware) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		bearerTokenSlice, ok := c.GetReqHeaders()["Authorization"]
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("unauthorized"))
		}

		bearerToken := bearerTokenSlice[0]
		bearerToken = strings.Split(bearerToken, " ")[1]

		token, err := jwt.Parse(bearerToken, func(t *jwt.Token) (any, error) {
			return []byte(m.cfg.App.JWTSecretKey), nil
		}, jwt.WithExpirationRequired(), jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))

		if err != nil {
			m.logger.Sugar().Error("unable to validate token: %v", zap.Error(err), zap.Uint64("request_id", ctx.ID()))
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("unauthorized"))
		}

		uID, err := token.Claims.GetSubject()
		if err != nil {
			m.logger.Sugar().Error("invalid token: %v", zap.Error(err), zap.Uint64("request_id", ctx.ID()))
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("unauthorized"))
		}

		_, err = gorm.G[models.User](m.db.DB).Where("id = ?", uID).Take(ctx)
		if err != nil {
			m.logger.Sugar().Error("unable to get user", zap.Error(err), zap.Uint64("request_id", ctx.ID()))
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("unauthorized"))
		}

		ctx.SetUserValue("userID", uID)

		t, err := m.cache.Get(ctx, fmt.Sprintf("AUTH_TOKEN_USER:%v", uID))
		if err == nil && t == nil {
			m.logger.Sugar().Errorw("🚀 ~ middleware ~ m.AuthMiddleware ~ auth token does not exist", "request_id", ctx.ID())
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("unauthorized"))
		}

		return c.Next()

	}
}
