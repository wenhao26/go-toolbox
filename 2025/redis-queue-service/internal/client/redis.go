package client

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"toolbox/2025/redis-queue-service/config"
)

// NewRedisClient 创建一个新Redis客户端实例
func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// 使用带有超时的上下文来ping-redis，避免无限等待
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ping-redis服务器以检查是否连接成功
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Println("Redis服务连接失败")
		return nil, fmt.Errorf("无法连接到Redis:%w", err)
	}

	return client, nil
}

// CloseRedisClient 关闭Redis客户端连接
// 建议在应用程序关闭时调用此函数，释放资源
func CloseRedisClient(client *redis.Client) {
	if client != nil {
		if err := client.Close(); err != nil {
			fmt.Println("Redis客户端关闭失败:", err)
		} else {
			fmt.Println("Redis客户端关闭成功")
		}
	}
}
