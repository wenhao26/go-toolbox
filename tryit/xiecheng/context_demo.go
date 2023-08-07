package main

import (
	"context"
	"fmt"
	"time"
)

func handelRequest(ctx context.Context)  {
	go writeRedis(ctx)
	go writeMongo(ctx)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("handelRequest done!")
			return
		default:
			fmt.Println("handelRequest running...")
			time.Sleep(2e9)
		}
	}
}

func writeRedis(ctx context.Context)  {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("write redis done!")
			return
		default:
			fmt.Println("Redis is writing...")
			time.Sleep(2e9)
		}
	}
}

func writeMongo(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("write mongodb done!")
			return
		default:
			fmt.Println("mongodb is writing...")
			time.Sleep(2e9)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go handelRequest(ctx)

	time.Sleep(10e9)
	fmt.Println("停止所有的子协程！")
	cancel()

	time.Sleep(2e9)
}
