package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"toolbox/2026/novel-creation/internal/config"
	"toolbox/2026/novel-creation/internal/queue"
	"toolbox/2026/novel-creation/internal/worker"
	"toolbox/2026/novel-creation/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	if err := logger.Init(cfg.Log.Level, cfg.Log.Output); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 3. 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Log.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Log.Info("Redis connected", zap.String("addr", cfg.Redis.Addr))

	// 4. 初始化队列（含恢复未完成的任务）
	redisQueue := queue.NewRedisQueue(rdb, cfg.Queue.Name, cfg.Queue.ProcessingSuffix, cfg.Queue.BlockTimeout)
	if err := redisQueue.RecoverOrphanTasks(ctx); err != nil {
		logger.Log.Error("Failed to recover orphan tasks", zap.Error(err))
	}

	// 5. 启动 worker 池
	workerPool := worker.NewWorkerPool(redisQueue, cfg.Worker.Count, cfg.Worker.GracefulTimeout)
	workerPool.Start(ctx)

	// 6. 等待退出信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	logger.Log.Info("received signal, shutting down", zap.String("signal", sig.String()))

	// 7. 优雅退出
	workerPool.Stop()

	// 8. 关闭 Redis 连接
	if err := redisQueue.Close(); err != nil {
		logger.Log.Error("Failed to close Redis", zap.Error(err))
	}

	logger.Log.Info("service stopped gracefully")
}
