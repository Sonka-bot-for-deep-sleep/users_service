package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Sonka-bot-for-deep-sleep/user_service/application/dto"
	"github.com/Sonka-bot-for-deep-sleep/user_service/application/models"
)

type redisInterface interface {
	Get(ctx context.Context, key string, out interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

type user struct {
	repo  repository
	redis redisInterface
}

func New(repo repository, redis redisInterface) *user {
	return &user{
		repo:  repo,
		redis: redis,
	}
}

func (u *user) GetByTgID(ctx context.Context, user *dto.GetUserByTgID) (*models.User, error) {
	if user.TgId == "" {
		return nil, fmt.Errorf("GetByTgID: TgId is empty")
	}

	var receivedUserFromRedis models.User
	userKey := u.getUserRedisKey(user.TgId)
	ok, err := u.redis.Get(ctx, userKey, receivedUserFromRedis)
	if err != nil {
		return nil, fmt.Errorf("GetByTgID: failed get data from redis: %w", err)
	}

	if !ok {
		receivedUser, err := u.repo.GetByTgID(ctx, user.TgId)
		if err != nil {
			return nil, fmt.Errorf("GetByTgID: failed get user by tg id: %w", err)
		}

		healthTimeCacheData := time.Minute * 10
		if err := u.redis.Set(ctx, u.getUserRedisKey(user.TgId), receivedUser, healthTimeCacheData); err != nil {
			return nil, fmt.Errorf("GetByTgID: failed set user to cache: %w", err)
		}
		return receivedUser, nil
	}

	return &receivedUserFromRedis, nil
}

func (u *user) CreateUser(ctx context.Context, user *dto.CreateUser) error {
	if user.Login == "" || user.Name == "" || user.TgID == "" {
		return fmt.Errorf("CreateUser: user data is empty")
	}

	wrappedUser := models.User{
		TgId:  user.TgID,
		Name:  user.Name,
		Login: user.Login,
	}
	if err := u.repo.CreateUser(ctx, wrappedUser); err != nil {
		return fmt.Errorf("CreateUser: failed creating user: %w", err)
	}

	return nil
}

func (u *user) getUserRedisKey(tgID string) string {
	return fmt.Sprintf("users:user:%s", tgID)
}
