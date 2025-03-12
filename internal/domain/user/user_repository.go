package user

import (
	"context"

	"github.com/Sonka-bot-for-deep-sleep/user_service/application/models"
)

type repository interface {
	GetByTgID(ctx context.Context, tgID string) (*models.User, error)
	CreateUser(ctx context.Context, user models.User) error
}
