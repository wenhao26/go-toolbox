package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"

	"toolbox/conf"
)

var (
	prefix = "test_"
	ctx    = context.Background()
	rdb    *redis.Client
)

func init() {
	cfg := conf.GetINI()
	sec := cfg.Section("redis")

	rdb = redis.NewClient(&redis.Options{
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
}

// 模拟数据
func mock() {
	//rdb.HSet(ctx, "chain", "key1", "value1", "key2", "value2").Result()
	//rdb.HSet(ctx, "chain", []string{"key3", "value3", "key4", "value4"})
	//rdb.HSet(ctx, "chain", map[string]interface{}{"key5": "value5", "key6": "value6"})

	for i := 1; i <= 50; i++ {
		hKey := prefix + strconv.Itoa(i)
		fmt.Println("hKey=", hKey)

		for n := 1; n <= 100; n++ {
			key := "key_" + strconv.Itoa(n)
			val := n + i
			rdb.HSet(ctx, hKey, key, val)
		}
	}
	fmt.Println("模拟数据初始化完成")
}

// 获取Keys
func keys() {
	result, err := rdb.Keys(ctx, "test_*").Result()
	if err != nil {
		panic(err)
	}

	for _, hKey := range result {
		fmt.Println("遍历Key=", hKey, "所有数据...")

		allData, err := rdb.HGetAll(ctx, hKey).Result()
		if err != nil {
			panic(err)
		}

		for k, v := range allData {
			fmt.Printf(" --key=%s,value=%s \n", k, v)
		}

		fmt.Println("\n")
	}
}

func main() {
	//mock()
	keys()
	fmt.Println("DONE")
}
