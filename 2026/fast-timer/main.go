package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

func main() {
	// 优雅退出控制
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("定时任务服务启动...")

	// 启动 0.5 秒高频处理
	go startTicker(ctx, "timer-0.5s", 200*time.Millisecond, HandleHealthCheck)

	<-ctx.Done()
	log.Println("服务正在关闭，等待任务收尾...")
}

// HandleLogStats 模拟每分钟的统计任务
func HandleLogStats(ctx context.Context) {
	log.Printf("[Task] 正在执行分钟级统计... 当前时间: %v", time.Now().Format("15:04:05"))
}

// HandleHealthCheck 模拟秒级的健康检查
func HandleHealthCheck(ctx context.Context) {
	log.Println("[Task] Ping OK")
}

// startTicker 生产级通用调度器
func startTicker(ctx context.Context, name string, interval time.Duration, task func(context.Context)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[-] 任务 %s 已停止", name)
			return
		case <-ticker.C:
			runSafe(ctx, name, task)
		}
	}
}

// runSafe 异常捕获封装
func runSafe(ctx context.Context, name string, task func(context.Context)) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Panic] 任务 %s 异常: %v\n堆栈: %s", name, r, string(debug.Stack()))
		}
	}()
	task(ctx)
}
