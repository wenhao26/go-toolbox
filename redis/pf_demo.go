package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"

	"toolbox/conf"
)

func redisDB() *redis.Client {
	var ctx = context.Background()

	cfg := conf.GetINI()
	sec := cfg.Section("redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     sec.Key("addr").String(),
		Password: sec.Key("password").String(),
		DB:       sec.Key("db").MustInt(),
		//PoolSize: 10,
	})
	//ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	//defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return rdb
}

func main() {
	var ctx = context.Background()
	rdb := redisDB()

	for i := 0; i < 10; i++ {
		rdb.PFAdd(ctx, "pf_test_1", fmt.Sprintf("pf1key%d", i))
	}
	ret, err := rdb.PFCount(ctx, "pf_test_1").Result()
	log.Println(ret, err)
}
