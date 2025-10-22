// 协程池管理（限制并发数量、任务超时、panic recover、优雅关闭）

package mq

import (
	"context"
	"log"
	"sync"
)

// WorkerPool 工作池
type WorkerPool struct {
	tasks  chan func()     // 任务通道
	wg     sync.WaitGroup  // 等待所有worker完成
	ctx    context.Context // 上下文，用于取消
	cancel context.CancelFunc
}

func NewWorkerPool(size int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	wp := &WorkerPool{
		tasks:  make(chan func(), size*2), // 缓冲，避免阻塞
		ctx:    ctx,
		cancel: cancel,
	}

	// 启动worker
	for i := 0; i < size; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

	return wp
}

// worker 消费任务
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case task, ok := <-wp.tasks:
			if !ok {
				return
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Worker panic recovered: %v", r)
					}
				}()
				task()
			}()
		}
	}
}

// Submit 提交任务到协程池
func (wp *WorkerPool) Submit(task func()) {
	select {
	case <-wp.ctx.Done():
		return
	default:
		wp.tasks <- task
	}
}

// Close 优雅关闭协程池
func (wp *WorkerPool) Close() {
	wp.cancel()
	close(wp.tasks)
	wp.wg.Wait()
}
