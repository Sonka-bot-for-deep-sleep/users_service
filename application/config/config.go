package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type dbCfg struct {
	USER     string `env:"POSTGRES_USER" env-default:"postgres"`
	PASSWORD string `env:"POSTGRES_PASSWORD" env-default:"12345"`
	DB       string `env:"POSTGRES_DB" env-default:"postgres"`
	PORT     string `env:"POSTGRES_PORT" env-default:"5432"`
	HOST     string `env:"POSTGRES_HOST" env-default:"localhost"`
}

type redisCfg struct {
	PASSWORD      string `env:"REDIS_PASSWORD" env-default:"12345"`
	HOST          string `env:"REDIS_HOST" env-default:"localhost"`
	PORT          string `env:"REDIS_PORT" env-default:"6379"`
	USER          string `env:"REDIS_USER" env-default:"redis"`
	USER_PASSWORD string `env:"REDIS_USER_PASSWORD" env-default:"12345"`
}

type appCfg struct {
	PORT string `env:"APP_PORT" env-default:"50051"`
}

type Config struct {
	DSN       string
	PORT      string
	REDIS_URL string
}

func MustLoad() (*Config, error) {
	var db dbCfg
	var app appCfg
	var redis redisCfg

	if err := cleanenv.ReadEnv(db); err != nil {
		return nil, fmt.Errorf("MustLoad: failed read env file and parse to db struct: %w", err)
	}

	if err := cleanenv.ReadEnv(app); err != nil {
		return nil, fmt.Errorf("MustLoad: failed read env file and parse to app struct: %w", err)
	}

	if err := cleanenv.ReadEnv(redis); err != nil {
		return nil, fmt.Errorf("MustLoad: failed read env file and parse to redis struct: %w", err)
	}

	redisURL := fmt.Sprintf("redis://%s:%s@%s:%s/0?protocol=3", redis.USER,
		redis.PASSWORD, redis.HOST, redis.PORT)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.USER, db.PASSWORD, db.HOST, db.PORT, db.DB)
	return &Config{
		DSN:       dsn,
		PORT:      app.PORT,
		REDIS_URL: redisURL,
	}, nil
}
