package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"toolbox/2026/high-perf-queue/internal/dispatcher"
)

func main() {
	log.Println("[Main] 初始化企业级高并发生产/消费系统...")

	// 初始化调度器，设置缓冲区为 1000，启动 5 个并发消费者
	engine := dispatcher.NewDispatcher(1000, 10)
	engine.Start()

	// 创建一个全链路控制 Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO 生产者
	go func() {
		for i := 1; i <= 100000; i++ {
			taskID := i

			// 定义每个任务专属的业务逻辑
			handler := func(taskCtx context.Context, payload any) error {
				// 模拟业务处理耗时
				time.Sleep(100 * time.Millisecond)

				// 模拟某一条数据发生不可预知的 panic
				if taskID == 100 {
					panic("模拟大厂生产环境偶发的运行时 Panic！")
				}

				fmt.Printf("[Business] 成功处理任务 ID: %d, 数据: %v\n", taskID, payload)
				return nil
			}

			// 投递到高性能调度器
			err := engine.Submit(ctx, fmt.Sprintf("Data-Payload-%d", taskID), handler)
			if err != nil {
				log.Printf("[Producer] 任务投递失败: %v\n", err)
				break
			}
		}
	}()

	// 监听操作系统信号，实现工业级优雅退出（Win & Linux 通用）
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞在此，直到收到 `Ctrl+C` 或 `kill` 信号
	sig := <-sigChan
	log.Printf("[Main] 捕获到系统信号: %v，准备触发停机流程...\n", sig)

	// 执行优雅退出，并设置硬性超时时间为 5 秒，防止进程悬挂
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := engine.Shutdown(shutdownCtx); err != nil {
		log.Printf("[Main] 停机过程中出现异常: %v\n", err)
	} else {
		log.Println("[Main] 进程安全退出。")
	}
}
