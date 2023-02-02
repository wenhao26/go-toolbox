package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"toolbox/conf"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func init() {
	cfg := conf.GetINI()
	sec := cfg.Section("redis")

	rdb = redis.NewClient(&redis.Options{
		Addr:     sec.Key("addr").String(),
		Password: sec.Key("password").String(),
		DB:       sec.Key("db").MustInt(),
		PoolSize: 10,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func main() {
	rdb.HSet(ctx, "user", []string{"key3", "value3", "key4", "value4"}).Result()

	fmt.Println("OK")
}
