package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"toolbox/2025/worker-example/pkg/task"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	queueName := "test_task_queue"
	log.Printf("生产者已启动，正在往队列 [%s] 写入数据...", queueName)

	// 模拟持续写入数据
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond) // 每 0.5 秒产生一个任务
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				// 构造任务数据
				newTask := task.Task{
					ID:      fmt.Sprintf("TASK-%d", t.UnixNano()),
					Payload: fmt.Sprintf("数据内容-%d", rand.Intn(100)),
				}

				// 序列化
				data, _ := json.Marshal(newTask)

				// 推送到 Redis (LPUSH)
				err := rdb.LPush(ctx, queueName, data).Err()
				if err != nil {
					log.Printf("写入失败: %v", err)
				} else {
					log.Printf("已成功投放任务: %s", newTask.ID)
				}
			}
		}
	}()

	<-ctx.Done()
	log.Println("生产者已停止")
}
