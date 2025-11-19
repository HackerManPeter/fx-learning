package auth

import (
	"context"
	"fmt"

	"github.com/hackermanpeter/fx-learning/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUserFromDTO(ctx context.Context, db *gorm.DB, user *models.CreateUserDTO) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return &models.User{}, fmt.Errorf("unable to hash password: %v", err)
	}

	u := models.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  string(hashedPassword),
	}

	if err := gorm.G[models.User](db).Create(ctx, &u); err != nil {
		return &models.User{}, fmt.Errorf("unable to create user: %v", err)
	}

	return &u, nil
}

func GetUserByEmail(ctx context.Context, db *gorm.DB, email string) (*models.User, error) {
	u, err := gorm.G[models.User](db).Where("email = ?", email).Take(ctx)
	if err != nil {
		return &models.User{}, err
	}

	return &u, nil
}
