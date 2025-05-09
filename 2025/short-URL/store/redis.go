package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

// SaveURL 保存URL
func SaveURL(shortURL, longURL string) error {
	return rdb.Set(ctx, shortURL, longURL, 600*time.Second).Err()
}

// GetURL 获取URL
func GetURL(shortURL string) (string, error) {
	return rdb.Get(ctx, shortURL).Result()
}
