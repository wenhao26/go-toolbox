package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"toolbox/2025/redis-queue-service/config"
	"toolbox/2025/redis-queue-service/internal/client"
	"toolbox/2025/redis-queue-service/internal/model"
	"toolbox/2025/redis-queue-service/internal/queue"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	if err := cfg.ValidateConfig(); err != nil {
		panic(err)
	}

	// 初始化redis客户端
	redisClient, err := client.NewRedisClient(&cfg.Redis)
	if err != nil {
		panic(err)
	}
	defer client.CloseRedisClient(redisClient) // 程序退出时关闭redis连接

	// 初始化消费者
	consumer := queue.NewRedisConsumer(redisClient, cfg.Queue.Name, cfg.Queue.Concurrency, cfg.Queue.BlockTimeout)

	// 定义消息处理函数
	// 这是核心业务逻辑，处理从队列中取出的每条消息。
	messageHandler := func(message []byte) error {
		msg, err := model.FromBytes(message)
		if err != nil {
			fmt.Println("对于无法解析的消息")
			return queue.ErrInvalidMessage // 返回自定义错误
		}

		fmt.Printf("[INFO] [%s] Processing message ID: %s, Content: %s, CreatedAt: %s ...\n",
			time.Now().Format("2006-01-02 15:04:05"),
			msg.ID,
			msg.Content,
			msg.CreatedAt,
		)

		// 模拟消息处理逻辑
		// 可以在这里进行数据库操作、API调用等
		// 模拟处理失败的场景
		if msg.ID == "simulate_error_id" { // 可以通过特定ID触发错误
			return queue.ErrProcessingFailed
		}

		// 模拟耗时操作
		time.Sleep(100 * time.Millisecond)

		return nil
	}

	// 启动消费者在一个单独的Goroutine中运行
	// 这样main Goroutine可以继续监听信号
	_, consumerCancel := context.WithCancel(context.Background())
	var consumerWg sync.WaitGroup // 用于等待消费者完全停止
	consumerWg.Add(1)

	go func() {
		defer consumerWg.Done()
		if err := consumer.Consume(messageHandler); err != nil {
			fmt.Printf("消费者遇到严重错误:%v\n", err)
		}
		fmt.Println("消费者协程已结束")
	}()

	// 注册信号处理，以便优雅停机
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM) // 监听中断和终止信号

	fmt.Println("Consumer started. Press Ctrl+C to stop.")

	<-stopChan // 阻塞直到收到停止信号

	// 优雅地停止消费者
	consumer.Stop()   // 通知消费者停止新的消息获取
	consumerCancel()  // 取消消费者上下文
	consumerWg.Wait() // 等待所有消费者goroutine完成当前任务并退出
}
