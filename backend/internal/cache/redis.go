package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"backend/internal/config"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(cfg *config.Config) *Redis {
	addr := cfg.RedisHost + ":" + cfg.RedisPort

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Println("⚠️ Redis unavailable — running without cache")
		return nil
	}

	log.Println("✅ Redis connected")
	return &Redis{Client: rdb}
}
