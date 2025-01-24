package rdb

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RedisClient Redis客户端
type RedisStorage struct {
	Rdb *redis.Client
}

// NewRedisClient 创建Redis实例
func NewRedisClient(addr string, db int) *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	// 测试连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisStorage{Rdb: rdb}
}
