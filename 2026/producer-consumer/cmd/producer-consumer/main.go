// Command producer-consumer 是一个高并发生产者-消费者示例程序：
// 多个生产者模拟实时数据推送，可配置数量的 worker 并发消费并模拟不同耗时，
// 并支持接收 Ctrl+C（Windows）/ SIGINT / SIGTERM 信号实现优雅退出。
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"toolbox/2026/producer-consumer/internal/app"
	"toolbox/2026/producer-consumer/internal/config"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(cfg.LogLevel),
	}))
	slog.SetDefault(logger)

	// signal.NotifyContext 会在收到指定信号时自动 cancel 返回的 ctx，
	// 是 Go 1.16+ 之后官方推荐的优雅退出信号监听方式，替代了手写的
	// signal.Notify + select 样板代码。
	//
	// 在 Windows 上，Ctrl+C 会被 runtime 转换为 os.Interrupt；
	// syscall.SIGTERM 在 Windows 上没有真实的信号语义，但保留声明
	// 不影响编译和运行，同时保证了该程序可以无修改地跨平台部署到 Linux 容器。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application := app.New(cfg, logger)
	if err := application.Run(ctx); err != nil {
		logger.Error("应用程序因错误退出", "err", err)
		os.Exit(1)
	}

	logger.Info("应用程序已干净退出，再见")
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
