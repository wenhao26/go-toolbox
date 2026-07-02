package worker

import (
	"context"
	"log"
	"sync"
	"time"
	"toolbox/2026/novel-creation/internal/handler"
	"toolbox/2026/novel-creation/internal/queue"
	"toolbox/2026/novel-creation/pkg/logger"

	"go.uber.org/zap"
)

type WorkerPool struct {
	queue           *queue.RedisQueue
	workerCount     int
	gracefulTimeout time.Duration
	wg              sync.WaitGroup
	cancelFunc      context.CancelFunc
}

func NewWorkerPool(q *queue.RedisQueue, workerCount int, gracefulTimeoutSec int) *WorkerPool {
	return &WorkerPool{
		queue:           q,
		workerCount:     workerCount,
		gracefulTimeout: time.Duration(gracefulTimeoutSec) * time.Second,
	}
}

// Start 启动所有 worker goroutines
func (p *WorkerPool) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	p.cancelFunc = cancel

	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}
	logger.Log.Info("worker pool started", zap.Int("count", p.workerCount))
}

// worker 是每个工作者的主循环
func (p *WorkerPool) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	logger.Log.Debug("worker started", zap.Int("id", id))

	for {
		select {
		case <-ctx.Done():
			logger.Log.Debug("worker stopping", zap.Int("id", id))
			return
		default:
			log.Printf("[worker %d]Running...", id)

			// 从队列获取任务，阻塞期间也能响应 ctx 取消（因为 BRPopLPush 支持 context 取消）
			taskData, err := p.queue.PopTask(ctx)
			if err != nil {
				logger.Log.Error("pop task error", zap.Int("id", id), zap.Error(err))
				// 短暂休眠避免错误时疯狂重试
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if taskData == "" {
				// 没有任务（超时），继续循环
				continue
			}

			// 处理任务
			err = handler.Handle(ctx, taskData)
			if err == nil {
				// 处理成功，确认任务
				if ackErr := p.queue.AckTask(ctx, taskData); ackErr != nil {
					logger.Log.Error("ack task failed", zap.Int("id", id), zap.String("task", taskData), zap.Error(ackErr))
				} else {
					logger.Log.Debug("task processed successfully", zap.Int("id", id), zap.String("task", taskData))
				}
			} else {
				// 处理失败，根据策略重试或放回队列（这里简单放回源队列）
				logger.Log.Error("task handling failed, requeue", zap.Int("id", id), zap.String("task", taskData), zap.Error(err))
				if requeueErr := p.queue.RequeueTask(ctx, taskData); requeueErr != nil {
					logger.Log.Error("requeue task failed", zap.Int("id", id), zap.Error(requeueErr))
				}
			}
		}
	}
}

// Stop 优雅停止：不再接受新任务，等待正在处理的任务完成，超时则强制退出
func (p *WorkerPool) Stop() {
	logger.Log.Info("stopping worker pool, no new tasks will be accepted")
	p.cancelFunc() // 通知所有 worker 退出 PopTask 循环

	// 等待所有 worker 完成，但不超过 gracefulTimeout
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Log.Info("all workers finished gracefully")
	case <-time.After(p.gracefulTimeout):
		logger.Log.Warn("graceful timeout exceeded, forcing shutdown")
	}
}
