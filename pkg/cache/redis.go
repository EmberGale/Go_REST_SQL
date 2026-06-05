package cache

import (
	"GoRestSQL/pkg/config"
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedisClient(ctx context.Context, cfg config.RedisConfig, log *zap.SugaredLogger) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		DB:           cfg.DB,
		Username:     cfg.User,
		Password:     cfg.Password,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		log.Info("Redis server ping fail")
		return nil, err
	}

	return db, nil
}
