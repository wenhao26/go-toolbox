package worker

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"
	"toolbox/2025/worker-example/pkg/task"

	"github.com/redis/go-redis/v9"
)

// Dispatcher 逻辑调度器
type Dispatcher struct {
	rdb        *redis.Client
	queueName  string
	maxWorkers int
	taskCh     chan task.Task
	wg         sync.WaitGroup
}

// NewDispatcher 创建新的调度器
func NewDispatcher(rdb *redis.Client, queue string, workers int) *Dispatcher {
	return &Dispatcher{
		rdb:        rdb,
		queueName:  queue,
		maxWorkers: workers,
		taskCh:     make(chan task.Task, workers), // 缓冲通道，防止频繁阻塞
	}
}

// Start 启动Worker池和分发器
func (d *Dispatcher) Start(ctx context.Context) {
	// 启动固定数量的Worker协程
	for i := 0; i < d.maxWorkers; i++ {
		d.wg.Add(1)
		go d.worker(ctx, i)
	}

	// 循环从Redis拉取任务
	log.Printf("[Dispatcher] 启动成功，并发数: %d", d.maxWorkers)

	for {
		select {
		case <-ctx.Done():
			log.Println("[Dispatcher] 停止抓取新任务")
			close(d.taskCh) // 关闭通道，通知Worker退出
			d.wg.Wait()     // 等到所有已抓取的任务执行完毕
			return
		default:
			// 使用 BLPOP 阻塞方式获取任务，超时设为短时间以便响应ctx.Done
			result, err := d.rdb.BLPop(ctx, time.Second*2, d.queueName).Result()
			if err != nil {
				continue
			}

			var t task.Task
			if err := json.Unmarshal([]byte(result[1]), &t); err != nil {
				log.Printf("解析任务失败: %v", err)
				continue
			}
			d.taskCh <- t // 将任务发送给空闲的Worker
		}
	}
}

// worker 工作者
func (d *Dispatcher) worker(ctx context.Context, id int) {
	defer d.wg.Done()
	log.Printf("[Worker %d] 已就绪", id)

	for t := range d.taskCh {
		// 这里的 ctx 可以根据需要决定是否传递给 Processor，
		// 如果希望 Shutdown 时立即强杀任务则传 ctx，希望处理完当前任务则用 context.Background
		if err := ProcessTask(ctx, t.ID, t.Payload); err != nil {
			log.Printf("[Worker %d] 任务 %s 失败: %v", id, t.ID, err)
		}
	}
}
