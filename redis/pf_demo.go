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

	/*for i := 0; i < 100000; i++ {
		rdb.PFAdd(ctx, "pf_test_1", fmt.Sprintf("pf1key%d", i))
	}
	ret, err := rdb.PFCount(ctx, "pf_test_1").Result()
	log.Println(ret, err)*/

	//  设置 HyperLogLog 类型的键  pf_test_2
	for i := 0; i < 10; i++ {
		rdb.PFAdd(ctx, "pf_test_2", fmt.Sprintf("pf2key%d", i))
	}
	ret2, _ := rdb.PFCount(ctx, "pf_test_2").Result()
	log.Println(ret2)

	for i := 0; i < 10; i++ {
		rdb.PFAdd(ctx, "pf_test_3", fmt.Sprintf("pf3key%d", i))
	}
	ret3, _ := rdb.PFCount(ctx, "pf_test_3").Result()
	log.Println(ret3)

	//  合并两个 HyperLogLog 类型的键  pf_test_2 + pf_test_3
	rdb.PFMerge(ctx, "pf_test", "pf_test_2", "pf_test_3")
	retAll, _ := rdb.PFCount(ctx, "pf_test").Result()
	log.Println(retAll)
}
