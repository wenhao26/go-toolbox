package engine

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"toolbox/2026/async-engine/model"
)

// Engine 提供异步任务消费能力，通过有界队列隔离生产与消费的速率差异。
// 队列满时投递不会阻塞，直接丢弃任务以保证生产侧零延迟。
type Engine struct {
	cfg     Config
	queue   chan *model.Task
	workers []*worker
	wg      sync.WaitGroup

	submitTotal    int64
	submitSuccess  int64
	submitDrop     int64
	consumeTotal   int64
	consumeLatency int64
}

// Snapshot 保存某一时刻的引擎指标快照。
type Snapshot struct {
	SubmitTotal   int64
	SubmitSuccess int64
	SubmitDrop    int64
	ConsumeTotal  int64
	QueueDepth    int64
	QueueCapacity int64
	AvgLatencyMs  float64
}

// NewEngine 创建引擎实例，此时消费协程尚未启动。
func NewEngine(cfg Config) *Engine {
	return &Engine{
		cfg:   cfg,
		queue: make(chan *model.Task, cfg.QueueSize),
	}
}

// Start 启动所有消费协程，必须在 Submit 之前调用。
func (e *Engine) Start() {
	for i := 0; i < e.cfg.WorkerNum; i++ {
		w := &worker{
			id:     i,
			queue:  e.queue,
			engine: e,
		}
		e.workers = append(e.workers, w)
		e.wg.Add(1)
		go w.run()
	}
}

// Submit 将任务投递到队列，队列满时立即返回 false，不阻塞调用方。
func (e *Engine) Submit(task *model.Task) bool {
	atomic.AddInt64(&e.submitTotal, 1)

	select {
	case e.queue <- task:
		atomic.AddInt64(&e.submitSuccess, 1)
		return true
	default:
		atomic.AddInt64(&e.submitDrop, 1)
		return false
	}
}

// Stop 优雅关闭引擎：关闭队列，等待所有 Worker 消费完剩余任务，或超时强制返回。
func (e *Engine) Stop() {
	close(e.queue)

	done := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("[engine] all workers exited gracefully")
	case <-time.After(10 * time.Second):
		fmt.Println("[engine] shutdown timeout, force exit")
	}
}

// Snapshot 原子读取当前运行指标，用于统计报表采样。
func (e *Engine) Snapshot() Snapshot {
	return Snapshot{
		SubmitTotal:   atomic.LoadInt64(&e.submitTotal),
		SubmitSuccess: atomic.LoadInt64(&e.submitSuccess),
		SubmitDrop:    atomic.LoadInt64(&e.submitDrop),
		ConsumeTotal:  atomic.LoadInt64(&e.consumeTotal),
		QueueDepth:    int64(len(e.queue)),
		QueueCapacity: int64(cap(e.queue)),
		AvgLatencyMs:  e.avgLatencyMs(),
	}
}

func (e *Engine) avgLatencyMs() float64 {
	total := atomic.LoadInt64(&e.consumeTotal)
	if total == 0 {
		return 0
	}
	return float64(atomic.LoadInt64(&e.consumeLatency)) / float64(total) / 1e6
}

func (e *Engine) recordConsume(latency time.Duration) {
	atomic.AddInt64(&e.consumeTotal, 1)
	atomic.AddInt64(&e.consumeLatency, int64(latency))
}
