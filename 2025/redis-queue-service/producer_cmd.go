package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

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

	// 初始化生产者
	producer := queue.NewRedisProducer(redisClient, cfg.Queue.Name)

	// 注册信号处理，用于优雅退出
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM) // 监听中断和终止信号

	// 模拟持续发送消息
	ticker := time.NewTicker(100 * time.Millisecond) // 每100毫秒发送一条消息
	defer ticker.Stop()

	msgCount := 0
	for {
		select {
		case <-ticker.C:
			// 构造示例消息
			msgID := uuid.New().String()
			msg := &model.ExampleMessage{
				ID:        msgID,
				Content:   fmt.Sprintf("Hello from producer! Message %d", msgCount),
				CreatedAt: time.Now(),
			}

			// 序列化消息
			messageBytes, err := msg.ToBytes()
			if err != nil {
				continue
			}

			// 发送消息
			if err := producer.Enqueue(messageBytes); err != nil {
				fmt.Println("消息排队失败 - ", msg.ID)
			} else {
				msgCount++
				fmt.Printf("[%s]消息已排队,总计发送[%d]\n", msg.ID, msgCount)
			}
		case <-stopChan:
			return
		}
	}
}
