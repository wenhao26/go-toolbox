package rdb

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Options struct {
	Addr     string
	Password string
	DB       int
}

type Rdb struct {
	client *redis.Client
}

func NewRdb(opt *Options) *Rdb {
	rdb := redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return &Rdb{client: rdb}
}

func (rdb *Rdb) GetClient() *redis.Client {
	return rdb.client
}

func (rdb *Rdb) Set(key string, value interface{}, ttl time.Duration) {
	err := rdb.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		panic(err)
	}
}

func (rdb *Rdb) Get(key string) string {
	result, err := rdb.client.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return result
}
