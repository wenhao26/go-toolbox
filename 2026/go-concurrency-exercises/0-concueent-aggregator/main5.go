//////////////////////////////
// AI 创作系统的“高并发任务调度器”
//
// AI 创作系统现在面临一个复杂情况:
// 1、请求归并 (SingleFlight)：多个用户可能同时请求生成同一个爆款小说的“大纲”（这个过程很费钱、很慢）。
// 2、并发限流 (Semaphore)：你的 Azure 账号只有 3 个并发额度，不能超载。
// 3、超时控制 (Context)：如果 AI 生成超过 10 秒没反应，必须强制断开，不能卡死。
//

package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// TaskScheduler 综合调度器
type TaskScheduler struct {
	sf    singleflight.Group // 官方归并库
	limit chan struct{}      // 并发限流桶
}

// NewTaskScheduler 初始化任务调度器
func NewTaskScheduler(maxConcurrent int) *TaskScheduler {
	return &TaskScheduler{
		limit: make(chan struct{}, maxConcurrent),
	}
}

// GenerateAction 生成创作内容方法
func (s *TaskScheduler) GenerateAction(topic string) (string, error) {
	// 第一层防御，SingleFlight 归并处理
	v, err, shared := s.sf.Do(topic, func() (interface{}, error) {
		// 第二层防御，限流排队
		log.Printf("⏳ 任务 [%s] 正在等待获取 API 令牌...", topic)

		s.limit <- struct{}{}        // 占位，拿到令牌
		defer func() { <-s.limit }() // 执行完释放

		log.Printf("🚀 任务 [%s] 成功获取令牌，开始调用 AI...", topic)

		// 第三层防御，Context 超时控制
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// TODO 模拟真实的 AI 调用
		return s.mockAiCall(ctx, topic)
	})
	if err != nil {
		return "", err
	}

	result := v.(string)
	if shared {
		log.Printf("🤝 优化成功：话题 [%s] 使用了合并后的共享结果", topic)
	}
	return result, nil
}

// mockAiCall 模拟 AI 调用
func (s *TaskScheduler) mockAiCall(ctx context.Context, topic string) (string, error) {
	select {
	case <-time.After(3 * time.Second): // 模拟生成内容需要 3 秒
		return fmt.Sprintf("【AI 生成完成】话题：%s", topic), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	scheduler := NewTaskScheduler(2)

	var wg sync.WaitGroup

	topics := []string{"修仙类", "修仙类", "赛博朋克", "修仙类", "末世", "末世", "末世", "修仙类", "赛博朋克"}

	for i, topic := range topics {
		wg.Add(1)

		go func(id int, topic string) {
			defer wg.Done()
			log.Printf("用户 %d 发起话题请求: %s", id, topic)
			res, err := scheduler.GenerateAction(topic)
			if err != nil {
				log.Printf("用户 %d 失败: %v", id, err)
			} else {
				log.Printf("用户 %d 拿到结果: %s", id, res)
			}
		}(i, topic)
	}

	wg.Wait()
	log.Println("🎉 所有任务处理完毕")
}
