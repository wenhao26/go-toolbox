package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"toolbox/2025/worker-example/internal/worker"

	"github.com/redis/go-redis/v9"
)

func main() {
	// 初始化Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	// 初始化调度器
	dispatcher := worker.NewDispatcher(rdb, "test_task_queue", 10)

	// 监听系统信号，优雅退出
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Worker 程序已启动...")
	dispatcher.Start(ctx)
	log.Println("程序已安全退出")
}
