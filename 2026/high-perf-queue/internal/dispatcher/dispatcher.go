package dispatcher

import (
	"context"
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

var (
	ErrDispatcherClosed = errors.New("dispatcher is closed")
)

// Dispatcher 负责协调高并发的生产与消费流控
type Dispatcher struct {
	capacity    int            // 内部缓冲队列容量
	workerCount int            // 消费者并发协程数
	taskQueue   chan *Task     // 核心流转通道
	wg          sync.WaitGroup // 跟踪所有 Worker 的生命周期
	isClosed    int32          // 原子标记，防止重复关闭导致的 panic
}

// NewDispatcher 初始化调度器
func NewDispatcher(capacity int, workerCount int) *Dispatcher {
	return &Dispatcher{
		capacity:    capacity,
		workerCount: workerCount,
		taskQueue:   make(chan *Task, capacity),
	}
}

// Start 统一初始化并拉起多消费者集群
func (d *Dispatcher) Start() {
	log.Printf("[Dispatcher] 正在启动系统... 缓冲区容量: %d, 消费者数量: %d\n", d.capacity, d.workerCount)
	for i := 0; i < d.workerCount; i++ {
		worker := NewWorker(i, d.taskQueue, &d.wg)
		worker.Start()
	}
	log.Println("[Dispatcher] 所有消费者协程已就绪")
}

// Submit 生产者投递任务（高频调用，必须轻量）
func (d *Dispatcher) Submit(ctx context.Context, payload any, handler TaskHandler) error {
	// 快路径判断：如果调度器已经关闭，拒绝写入
	if atomic.LoadInt32(&d.isClosed) == 1 {
		return ErrDispatcherClosed
	}

	// 从对象池申请内存，降低 GC 开销
	task := AcquireTask(ctx, payload, handler)

	// 阻塞写入通道，同时监听上下文状态
	select {
	case d.taskQueue <- task:
		return nil
	case <-ctx.Done():
		// 如果因通道挤满导致阻塞，且外部设置了超时，则释放资源并安全返回
		ReleaseTask(task)
		return ctx.Err()
	}
}

// Shutdown 优雅停机（核心设计：不丢数据）
func (d *Dispatcher) Shutdown(ctx context.Context) error {
	// CAS 保证并发安全与幂等关闭
	if !atomic.CompareAndSwapInt32(&d.isClosed, 0, 1) {
		return nil
	}

	log.Println("[Dispatcher] 接收到停机指令，正在关闭任务投递通道...")

	// 先关闭通道：此时生产者无法再 Submit 任务
	close(d.taskQueue)

	// 异步等到所有消费者把通道里的积压任务消费干净
	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	// 阻塞等待，同时支持强行终止超时
	select {
	case <-done:
		log.Println("[Dispatcher] 所有积压任务处理完毕，系统优雅退出成功。")
		return nil
	case <-ctx.Done():
		log.Println("[Dispatcher] 优雅停机超时！部分积压任务可能未处理完成，强行退出。")
		return ctx.Err()
	}
}
