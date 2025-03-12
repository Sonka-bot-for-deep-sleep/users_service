package main

import (
	"fmt"

	"github.com/Sonka-bot-for-deep-sleep/common/pkg/logger"
	"github.com/Sonka-bot-for-deep-sleep/user_service/application/config"
	"github.com/Sonka-bot-for-deep-sleep/user_service/internal/domain/user"
	"github.com/Sonka-bot-for-deep-sleep/user_service/internal/infrastructure/db/postgres"
	"github.com/Sonka-bot-for-deep-sleep/user_service/internal/infrastructure/db/redis"
	"github.com/Sonka-bot-for-deep-sleep/user_service/internal/infrastructure/grpc"
	"github.com/Sonka-bot-for-deep-sleep/user_service/internal/infrastructure/grpc/handlers"
	"github.com/Sonka-bot-for-deep-sleep/user_service/internal/infrastructure/repository"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.New()
	if err != nil {
		fmt.Println("Error create logger instance")
		return
	}

	cfg, err := config.MustLoad()
	if err != nil {
		log.Error("Failed load config data", zap.Error(err))
		return
	}

	pgConn, err := postgres.NewWithConn(cfg.DSN)
	if err != nil {
		log.Error("Failed create conn to postgres database", zap.Error(err))
		return
	}
	redisConn, err := redis.NewWithConn(cfg.REDIS_URL)
	if err != nil {
		log.Error("Failed conn to redis", zap.Error(err))
		return
	}
	repo := repository.NewUser(pgConn.DB)
	entity := user.New(repo, redisConn)
	handler := handlers.New(entity, log)
	srv := grpc.New(handler, log)

	if err := srv.StartServer(cfg.PORT); err != nil {
		log.Error("Failed start server", zap.Error(err))
		return
	}
}
