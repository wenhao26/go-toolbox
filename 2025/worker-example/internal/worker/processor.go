package worker

import (
	"context"
	"log"
	"time"
)

// ProcessTask 处理具体的长耗时业务逻辑
func ProcessTask(ctx context.Context, id string, payload string) error {
	// 模拟业务耗时
	log.Printf("[Processor] 正在处理任务: %s", id)

	select {
	case <-time.After(2 * time.Second): // 模拟耗时2秒
		log.Printf("[Processor] 任务已完成: %s %s", id, payload)
	case <-ctx.Done():
		log.Printf("[Processor] 任务被取消: %s", id)
		return ctx.Err()
	}
	return nil
}
