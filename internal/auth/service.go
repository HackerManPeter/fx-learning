package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hackermanpeter/fx-learning/internal/cache"
	"github.com/hackermanpeter/fx-learning/internal/config"
	"github.com/hackermanpeter/fx-learning/internal/database"
	"github.com/hackermanpeter/fx-learning/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	logger    *zap.Logger
	cfg       *config.Config
	db        *database.Database
	validator *validator.Validate
	cache     *cache.Cache
}

func NewAuthService(logger *zap.Logger, cfg *config.Config, db *database.Database, validator *validator.Validate, cache *cache.Cache) *AuthService {
	return &AuthService{
		logger,
		cfg,
		db,
		validator,
		cache,
	}
}

func (a *AuthService) signUp() fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := models.CreateUserDTO{}
		if err := c.BodyParser(&u); err != nil {
			a.logger.Sugar().Errorf("unable to parse body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("unable to parse body"))
		}

		if err := a.validator.Struct(&u); err != nil {
			a.logger.Sugar().Errorf("unable to validate data: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("unable to validate data"))
		}

		ctx := c.Context()

		_, err := GetUserByEmail(ctx, a.db.DB, u.Email)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			if err != nil {
				a.logger.Sugar().Errorf("unable to validate user: %v", err)
				return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("unable to validate user"))
			}
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(fmt.Sprintf("user with email %v already exists", u.Email)))
		}

		user, err := CreateUserFromDTO(ctx, a.db.DB, &u)
		if err != nil {
			a.logger.Sugar().Errorf("unable to create user: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("unable to create user"))
		}

		return c.Status(fiber.StatusCreated).JSON(models.NewSuccessResponse(user, "successfully created user"))

	}
}

func (a *AuthService) login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		l := models.LoginRequestDTO{}

		if err := c.BodyParser(&l); err != nil {
			a.logger.Sugar().Errorf("unable to process request data: %v", err)
			return c.Status(fiber.StatusUnprocessableEntity).JSON("unable to process request data")
		}

		if err := a.validator.Struct(&l); err != nil {
			a.logger.Sugar().Errorf("unable to process request data: %v", err)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fmt.Sprintf("unable to process request data: %v", err))
		}

		ctx := c.Context()

		u, err := GetUserByEmail(ctx, a.db.DB, l.Email)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("invalid credentials"))
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(l.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("invalid credentials"))
		}

		response, err := signJWT(u.ID.String(), a.cfg.App.JWTSecretKey)
		if err != nil {
			a.logger.Sugar().Errorf("unable to generate jwt: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("unable to generate jwt"))
		}

		expiration := time.Until(response.ExpiresAt)

		err = a.cache.Set(ctx, fmt.Sprintf("AUTH_TOKEN_USER:%v", u.ID.String()), response.Token, expiration)
		if err != nil {
			a.logger.Sugar().Errorw("🚀 ~ auth ~ a.login ~ unable to store jwt in cache", "err", err)
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("unable to login, please retry"))
		}

		return c.Status(fiber.StatusCreated).JSON(models.NewSuccessResponse(response, "login successful"))

	}
}

func (a *AuthService) logout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		err := a.cache.Del(ctx, fmt.Sprintf("AUTH_TOKEN_USER:%v", ctx.UserValue("userID")))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("unable to log out, please try again"))
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(nil, "logout successful"))

	}
}

func signJWT(userID, key string) (*models.LoginResponseDTO, error) {
	exp := time.Now().Add(2 * time.Hour)

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	jwt, err := token.SignedString([]byte(key))
	if err != nil {
		return &models.LoginResponseDTO{}, fmt.Errorf("unable to generate jwt: %v", err)
	}

	return &models.LoginResponseDTO{
		Token:     jwt,
		ExpiresAt: exp,
	}, nil

}
